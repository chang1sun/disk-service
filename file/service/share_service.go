package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/changpro/disk-service/common/constants"
	"github.com/changpro/disk-service/common/errcode"
	"github.com/changpro/disk-service/file/repo"
	"github.com/changpro/disk-service/file/rpc/auth"
	"github.com/changpro/disk-service/file/util"
	"google.golang.org/grpc/status"
)

func CreateShare(ctx context.Context, dto *repo.CreateShareDTO) (string, error) {
	doc, err := GetFileDetail(ctx, dto.UserID, dto.DocID)
	if err != nil {
		return "", err
	}
	if doc == nil {
		return "", status.Error(errcode.FindNoFileInServerCode, errcode.FindNoFileInServerMsg)
	}
	// cal size and file num
	size, fileNum, err := GetDirSizeAndSubFilesNum(ctx, doc)
	if err != nil {
		return "", err
	}
	// wrap into po
	po := &repo.ShareDetailPO{
		Uploader:    doc.UserID,
		DocID:       doc.ID,
		DocName:     doc.Name,
		DocSize:     size,
		DocType:     doc.IsDir,
		FileNum:     fileNum,
		ExpireHours: dto.ExpireHour,
	}
	data, err := json.Marshal(po)
	if err != nil {
		return "", status.Errorf(errcode.JsonMarshalErrCode, errcode.JsonMarshalErrMsg, err)
	}
	token := util.GetMd5FromJson(data)
	// write in redis
	err = repo.GetShareDao().CreateShareToken(ctx, token, po)
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	// create a record
	if err = createPostShareRecord(ctx, dto.UserID, po.DocName, doc.IsDir); err != nil {
		return "", err
	}
	return token, nil
}

func RetrieveShareFromToken(ctx context.Context, userID, token, path string) error {
	if err := checkPath(ctx, userID, path); err != nil {
		return err
	}
	po, err := GetShareDetail(ctx, token)
	if err != nil {
		return err
	}
	// update user's analytic data
	err = auth.GetAuthCaller().UpdateUserStorage(ctx, userID, po.FileNum, 0, po.DocSize)
	if err != nil {
		return status.Errorf(errcode.RPCCallErrCode, errcode.RPCCallErrMsg, err)
	}
	originalPO, err := GetFileDetail(ctx, po.Uploader, po.DocID)
	if err != nil {
		return err
	}
	// save to this user's file record
	if err := saveDocsForSaver(ctx, []*repo.UserFilePO{originalPO}, userID, path); err != nil {
		return err
	}
	// create a record
	if err = createSaveShareRecord(ctx, userID, po.DocName, po.Uploader, originalPO.IsDir); err != nil {
		return err
	}
	return nil
}

func saveDocsForSaver(ctx context.Context, pos []*repo.UserFilePO, userID, newPath string) error {
	if len(pos) == 0 {
		return nil
	}
	// if it is a folder, then trigger recursively call
	for _, po := range pos {
		// recursively look up sub folder or files and update them first
		subPOs, err := repo.GetUserFileDao().QueryDirByPath(ctx, po.UserID, po.Path+po.Name+"/", true)
		if err != nil {
			return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
		}
		if len(subPOs) > 0 {
			err = saveDocsForSaver(ctx, subPOs, userID, newPath+po.Name+"/")
			if err != nil {
				return err
			}
		}
		newPO := buildDocPOForSaver(userID, newPath, po)
		_, err = repo.GetUserFileDao().AddFileOrDir(ctx, newPO)
		if err != nil {
			return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
		}
	}
	return nil
}

func GetShareDetail(ctx context.Context, token string) (*repo.ShareDetailPO, error) {
	po, err := repo.GetShareDao().GetShareDetail(ctx, token)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if po == nil {
		return nil, status.Error(errcode.NoSuchShareCode, errcode.NoSuchShareMsg)
	}
	return po, nil
}

func buildDocPOForSaver(userID, path string, originalPO *repo.UserFilePO) *repo.UserFilePO {
	return &repo.UserFilePO{
		UserID:    userID,
		UniFileID: originalPO.UniFileID,
		Name:      originalPO.Name,
		FileMd5:   originalPO.FileMd5,
		FileSize:  originalPO.FileSize,
		FileType:  originalPO.FileType,
		Path:      path,
		IsDir:     originalPO.IsDir,
		Status:    originalPO.Status,
		IsHide:    notHide,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}
}

func createPostShareRecord(ctx context.Context, userID, docName string, dirOrFile int32) error {
	var docNotion string
	if dirOrFile == isDir {
		docNotion = "folder"
	} else {
		docNotion = "file"
	}
	record := &repo.ShareRecordPO{
		UserID:  userID,
		DocName: docName,
		Message: fmt.Sprintf(`%v shared a %v "%v" to the public at %v`, userID, docNotion,
			docName, time.Now().Format(constants.StandardTimeFormat)),
		CreateTime: time.Now(),
	}
	if err := repo.GetShareRecordDao().CreateShareRecord(ctx, record); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func createSaveShareRecord(ctx context.Context, userID, docName, uploader string, dirOrFile int32) error {
	var docNotion string
	if dirOrFile == isDir {
		docNotion = "folder"
	} else {
		docNotion = "file"
	}
	record := &repo.ShareRecordPO{
		UserID:  userID,
		DocName: docName,
		Message: fmt.Sprintf(`%v saved %v "%v" from %v at %v`, userID,
			docNotion, docName, uploader, time.Now().Format(constants.StandardTimeFormat)),
		CreateTime: time.Now(),
	}
	if err := repo.GetShareRecordDao().CreateShareRecord(ctx, record); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func GetShareRecordList(ctx context.Context, userID string, start int32, limit int32) ([]*repo.ShareRecordPO, int64, error) {
	list, count, err := repo.GetShareRecordDao().QueryRecordList(ctx, userID, start, limit)
	if err != nil {
		return nil, 0, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count == 0 {
		return nil, 0, nil
	}
	return list, count, nil
}

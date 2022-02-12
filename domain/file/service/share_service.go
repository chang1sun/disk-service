package service

import (
	"context"
	"encoding/json"
	"time"

	arepo "github.com/changpro/disk-service/domain/auth/repo"
	aservice "github.com/changpro/disk-service/domain/auth/service"
	"github.com/changpro/disk-service/domain/file/repo"
	"github.com/changpro/disk-service/infra/constants"
	"github.com/changpro/disk-service/infra/errcode"
	"github.com/changpro/disk-service/infra/util"
	"google.golang.org/grpc/status"
)

const (
	RecordTypeShare = 1
	RecordTypeSave  = 2
)

func CreateShare(ctx context.Context, dto *repo.CreateShareDTO) (string, string, error) {
	doc, err := GetFileDetail(ctx, dto.UserID, dto.DocID)
	if err != nil {
		return "", "", err
	}
	if doc == nil {
		return "", "", status.Error(errcode.FindNoFileInServerCode, errcode.FindNoFileInServerMsg)
	}
	// cal size and file num
	size, fileNum, err := GetDirSizeAndSubFilesNum(ctx, doc)
	if err != nil {
		return "", "", err
	}
	// wrap into po
	po := &repo.ShareDetailPO{
		Uploader:    doc.UserID,
		Password:    util.NewLenRandomString(4), // 4 bytes rand
		DocID:       doc.ID,
		DocName:     doc.Name,
		DocSize:     size,
		DocType:     doc.FileType,
		IsDir:       doc.IsDir,
		FileNum:     fileNum,
		ExpireHours: dto.ExpireHour,
	}
	data, err := json.Marshal(po)
	if err != nil {
		return "", "", status.Errorf(errcode.JsonMarshalErrCode, errcode.JsonMarshalErrMsg, err)
	}
	token := util.GetMd5FromJson(data)
	// write in redis
	err = repo.GetShareDao().CreateShareToken(ctx, token, po)
	if err != nil {
		return "", "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	// create a record
	if err = createPostShareRecord(ctx, dto.UserID, token, po); err != nil {
		return "", "", err
	}
	return token, po.Password, nil
}

func GetShareDetail(ctx context.Context, token, password string) (*repo.ShareDetailPO, error) {
	po, err := repo.GetShareDao().GetShareDetail(ctx, token)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if po == nil {
		return nil, status.Error(errcode.NoSuchShareCode, errcode.NoSuchShareMsg)
	}
	if po.Password != password {
		return nil, status.Error(errcode.WrongSharePasswordCode, errcode.WrongSharePasswordMsg)
	}
	return po, nil
}

func createPostShareRecord(ctx context.Context, userID, token string, share *repo.ShareDetailPO) error {
	record := &repo.ShareRecordPO{
		UserID:     userID,
		DocID:      share.DocID,
		DocName:    share.DocName,
		CreateTime: time.Now(),
		ExpireTime: time.Now().Add(time.Duration(share.ExpireHours) * time.Hour),
		Token:      token,
		Type:       RecordTypeShare,
	}
	if err := repo.GetShareRecordDao().CreateShareRecord(ctx, record); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

func GetShareRecordList(ctx context.Context, query *repo.RecordQuery) ([]*repo.ShareRecordPO, int64, error) {
	list, count, err := repo.GetShareRecordDao().QueryRecordList(ctx, query)
	if err != nil {
		return nil, 0, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if count == 0 {
		return nil, 0, nil
	}
	return list, count, nil
}

func RetrieveShareFromToken(ctx context.Context, userID, token, path string) error {
	if err := CheckPath(ctx, userID, path); err != nil {
		return err
	}
	po, err := repo.GetShareDao().GetShareDetail(ctx, token)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if po == nil {
		return status.Error(errcode.NoSuchShareCode, errcode.NoSuchShareMsg)
	}
	if err != nil {
		return err
	}
	// update user's analytic data
	err = aservice.UpdateUserStorage(ctx, &arepo.UpdateUserAnalysisDTO{
		UserID:  userID,
		FileNum: 1,
		Size:    po.DocSize,
	})
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
	if err = createSaveShareRecord(ctx, userID, token, po); err != nil {
		return err
	}
	return nil
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
		IsHide:    2,
		CreateAt:  time.Now(),
		UpdateAt:  time.Now(),
	}
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

func createSaveShareRecord(ctx context.Context, userID, token string, share *repo.ShareDetailPO) error {
	loc, _ := time.LoadLocation(constants.TimeZoneLocation)
	createTime, _ := time.ParseInLocation(constants.StandardTimeFormat, share.CreateTime, loc)
	record := &repo.ShareRecordPO{
		UserID:     userID,
		DocID:      share.DocID,
		DocName:    share.DocName,
		CreateTime: time.Now(),
		ExpireTime: createTime.Add(time.Duration(share.ExpireHours) * time.Hour),
		Token:      token,
		Type:       RecordTypeSave,
	}
	if err := repo.GetShareRecordDao().CreateShareRecord(ctx, record); err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return nil
}

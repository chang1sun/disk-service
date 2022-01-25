package application

import (
	"context"
	"fmt"
	"time"

	arepo "github.com/changpro/disk-service/domain/auth/repo"
	aservice "github.com/changpro/disk-service/domain/auth/service"
	"github.com/changpro/disk-service/domain/file/repo"
	"github.com/changpro/disk-service/domain/file/service"
	"github.com/changpro/disk-service/infra/constants"
	"github.com/changpro/disk-service/infra/errcode"
	"google.golang.org/grpc/status"
)

func RetrieveShareFromToken(ctx context.Context, userID, token, path string) error {
	if err := service.CheckPath(ctx, userID, path); err != nil {
		return err
	}
	po, err := service.GetShareDetail(ctx, token)
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
	originalPO, err := service.GetFileDetail(ctx, po.Uploader, po.DocID)
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

func createSaveShareRecord(ctx context.Context, userID, docName, uploader string, dirOrFile int32) error {
	var docNotion string
	if dirOrFile == 1 {
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

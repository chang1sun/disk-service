package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/changpro/disk-service/domain/file/repo"
	"github.com/changpro/disk-service/infra/constants"
	"github.com/changpro/disk-service/infra/errcode"
	"github.com/changpro/disk-service/infra/util"
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

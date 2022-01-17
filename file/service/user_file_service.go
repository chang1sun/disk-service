package service

import (
	"context"

	"github.com/changpro/disk-service/common/errcode"
	"github.com/changpro/disk-service/file/repo"
	"github.com/changpro/disk-service/file/rpc/auth"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"google.golang.org/grpc/status"
)

func QuickUpload(ctx context.Context, userID, fileName, md5 string) (string, error) {
	id, err := tryQuickUpload(ctx, userID, fileName, md5)
	if err != nil {
		return "", err
	}
	if id == "" {
		return "", status.Error(errcode.FindNoFileInServerCode, errcode.FindNoFileInServerMsg)
	}
	return id, nil
}

func tryQuickUpload(ctx context.Context, userID, fileName, md5 string) (string, error) {
	file, err := repo.GetUniFileStoreDao().QueryFileByMd5(ctx, md5)
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if file == nil {
		return "", nil
	}
	fileMeta := &repo.UniFileMetaPO{
		Size: file.Length,
		Md5:  md5,
		Type: file.Metadata.Lookup("type").String(),
	}

	// update user's storage size
	err = auth.GetAuthCaller().UpdateUserStorage(ctx, userID, 1, 0, fileMeta.Size)
	if err != nil {
		return "", status.Errorf(errcode.RPCCallErrCode, errcode.RPCCallErrMsg, err)
	}

	// update user's level content
	id, err := repo.GetUserFileDao().AddFile(ctx, userID,
		buildUploadUserFilePO(file.ID.(primitive.ObjectID).String(), fileName, userID, fileMeta))
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return id, nil
}

func GetUserRoot(ctx context.Context, userID string) ([]*repo.UserFilePO, error) {
	content, err := repo.GetUserFileDao().QueryUserRoot(ctx, userID, false)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return content, nil
}

func GetFileDetail(ctx context.Context, userID, fileID string) (*repo.UserFilePO, error) {
	detail, err := repo.GetUserFileDao().QueryFileDetail(ctx, userID, fileID)
	if err != nil {
		return nil, status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return detail, nil
}

func MakeNewFolder(ctx context.Context, userID, dirName, path string) (string, error) {
	// check path correctness
	ok, err := repo.GetUserFileDao().IsPathExist(ctx, userID, path)
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if !ok {
		return "", status.Errorf()
	}

	id, err := repo.GetUserFileDao().CreateFolder(ctx, userID, dirName, path)
	if err != nil {
		return "", status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	return id
}

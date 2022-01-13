package service

import (
	"context"

	"github.com/changpro/disk-service/common/errcode"
	"github.com/changpro/disk-service/file/repo"
	"google.golang.org/grpc/status"
)

func QuickUpload(ctx context.Context, user_id, file_md5 string) error {
	file, err := repo.GetUniFileStoreDao().QueryFileByMd5(ctx, file_md5)
	if err != nil {
		return status.Errorf(errcode.DatabaseOperationErrCode, errcode.DatabaseOperationErrMsg, err)
	}
	if file != nil {
		return nil
	}
	return nil
}

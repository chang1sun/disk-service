package repo

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/gridfs"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UniFileStoreDao struct {
	Database *mongo.Database
	Bucket   *gridfs.Bucket
}

var uniFileStoreDao *UniFileStoreDao

func GetUniFileStoreDao() *UniFileStoreDao {
	return uniFileStoreDao
}

func SetUniFileStoreDao(dao *UniFileStoreDao) {
	uniFileStoreDao = dao
}

func (dao *UniFileStoreDao) QueryFileByMd5(ctx context.Context, md5 string) (*gridfs.File, error) {
	return &gridfs.File{}, nil
}

func (dao *UniFileStoreDao) UploadFile(ctx context.Context, fileName string, f io.Reader, meta *UniFileMetaPO) (string, error) {
	id, err := dao.Bucket.UploadFromStream(fileName, f, &options.UploadOptions{
		Metadata: meta,
	})
	if err != nil {
		return "", nil
	}
	return id.String(), nil
}

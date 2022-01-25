package repo

import (
	"context"
	"io"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
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
	var files []gridfs.File
	filter := &bson.D{
		{"metadata.md5", md5},
	}
	cursor, err := dao.Bucket.Find(filter)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, files)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, nil
	}
	return &files[0], nil
}

func (dao *UniFileStoreDao) QueryFileByID(ctx context.Context, id string) ([]*gridfs.File, error) {
	var files []*gridfs.File
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	filter := &bson.D{
		{"_id", oid},
	}
	cursor, err := dao.Bucket.Find(filter)
	defer cursor.Close(context.TODO())
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, files)
	if err != nil {
		return nil, err
	}
	if len(files) == 0 {
		return nil, nil
	}
	return files, nil
}

func (dao *UniFileStoreDao) UploadFile(ctx context.Context, fileName string, f io.Reader, meta *UniFileMetaPO) (string, error) {
	id, err := dao.Bucket.UploadFromStream(fileName, f, &options.UploadOptions{
		Metadata: meta,
	})
	if err != nil {
		return "", err
	}
	return id.String(), nil
}

func (dao *UniFileStoreDao) GetDownloadStream(ctx context.Context, id string) (*gridfs.DownloadStream, error) {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, err
	}
	s, err := dao.Bucket.OpenDownloadStream(oid)
	if err != nil {
		return nil, err
	}
	return s, nil
}

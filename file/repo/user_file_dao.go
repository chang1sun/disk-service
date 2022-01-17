package repo

import (
	"context"

	"github.com/changpro/disk-service/file/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type UserFileDao struct {
	Database   *mongo.Database
	Collection string
}

var userFileDao *UserFileDao

func GetUserFileDao() *UserFileDao {
	return userFileDao
}

func SetUserFileDao(dao *UserFileDao) {
	dao.Collection = "user_files"
	userFileDao = dao
}

func (dao *UserFileDao) AddFile(ctx context.Context, userID string, filePO *UserFilePO) (string, error) {
	res, err := dao.Database.Collection(dao.Collection).InsertOne(ctx, filePO)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).String(), nil
}

func (dao *UserFileDao) QueryUserRoot(ctx context.Context, userID string, showHide bool) ([]*UserFilePO, error) {
	var content []*UserFilePO
	filter := bson.D{
		{"user_id", userID},
		{"path", "/"},
	}
	if !showHide {
		filter = append(filter, bson.E{"is_hide", constants.FileDisplayStatusShow})
	}
	cursor, err := dao.Database.Collection(dao.Collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var f UserFilePO
		err := cursor.Decode(&f)
		if err != nil {
			return nil, err
		}
		content = append(content, &f)
	}
	return content, nil
}

func (dao *UserFileDao) QueryFileDetail(ctx context.Context, userID, fileID string) (*UserFilePO, error) {
	var content UserFilePO
	oid, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{"_id", oid},
	}
	res := dao.Database.Collection(dao.Collection).FindOne(ctx, filter)
	if res.Err() != nil {
		return nil, res.Err()
	}
	err = res.Decode(&content)
	if err != nil {
		return nil, err
	}
	return &content, nil
}

func (dao *UserFileDao) QueryDirByPath(ctx context.Context, userID, path string, showHide bool) ([]*UserFilePO, error) {
	var content []*UserFilePO
	filter := bson.D{
		{"user_id", userID},
		{"path", path},
	}
	if !showHide {
		filter = append(filter, bson.E{"is_hide", constants.FileDisplayStatusShow})
	}
	cursor, err := dao.Database.Collection(dao.Collection).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	for cursor.Next(ctx) {
		var f UserFilePO
		err := cursor.Decode(&f)
		if err != nil {
			return nil, err
		}
		content = append(content, &f)
	}
	return content, nil
}

func (dao *UserFileDao) IsPathExist(ctx context.Context, userID string, path string) (bool, error) {
	filter := bson.D{
		{"user_id", userID},
		{"path", path},
	}
	res := dao.Database.Collection(dao.Collection).FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return false, nil
	}
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

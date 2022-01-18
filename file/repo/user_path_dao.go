package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

const collUserPath = "user_files"

type UserPathDao struct {
	Database *mongo.Database
}

var userPathDao *UserPathDao

func GetUserPathDao() *UserPathDao {
	return userPathDao
}

func SetUserPathDao(dao *UserPathDao) {
	userPathDao = dao
}

func (dao *UserPathDao) AddPath(ctx context.Context, userID string, path string) error {
	_, err := dao.Database.Collection(collUserPath).InsertOne(ctx, &UserPathPO{UserID: userID, Path: path})
	if err != nil {
		return err
	}
	return nil
}

func (dao *UserPathDao) IsPathExist(ctx context.Context, userID string, path string) (bool, error) {
	res := dao.Database.Collection(collUserPath).FindOne(ctx, &UserPathPO{UserID: userID, Path: path})
	if res.Err() == mongo.ErrNoDocuments {
		return false, nil
	}
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

func (dao *UserPathDao) DeletePath(ctx context.Context, userID string, path string) error {
	_, err := dao.Database.Collection(collUserPath).DeleteOne(ctx, &UserPathPO{UserID: userID, Path: path})
	if err != nil {
		return err
	}
	return nil
}

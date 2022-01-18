package repo

import (
	"context"

	"github.com/changpro/disk-service/file/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

const collUserFiles = "user_files"

type UserFileDao struct {
	Database *mongo.Database
}

var userFileDao *UserFileDao

func GetUserFileDao() *UserFileDao {
	return userFileDao
}

func SetUserFileDao(dao *UserFileDao) {
	userFileDao = dao
}

func (dao *UserFileDao) AddFileOrDir(ctx context.Context, po *UserFilePO) (string, error) {
	res, err := dao.Database.Collection(collUserFiles).InsertOne(ctx, po)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).String(), nil
}

func isPathExist() {}

func (dao *UserFileDao) QueryUserRoot(ctx context.Context, userID string, showHide bool) ([]*UserFilePO, error) {
	var content []*UserFilePO
	filter := bson.D{
		{"user_id", userID},
		{"path", "/"},
	}
	if !showHide {
		filter = append(filter, bson.E{"is_hide", constants.FileDisplayStatusShow})
	}
	cursor, err := dao.Database.Collection(collUserFiles).Find(ctx, filter)
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

func (dao *UserFileDao) QueryDetail(ctx context.Context, fileID string) (*UserFilePO, error) {
	var content UserFilePO
	oid, err := primitive.ObjectIDFromHex(fileID)
	if err != nil {
		return nil, err
	}
	filter := bson.D{
		{"_id", oid},
	}
	res := dao.Database.Collection(collUserFiles).FindOne(ctx, filter)
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
	cursor, err := dao.Database.Collection(collUserFiles).Find(ctx, filter)
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

func (dao *UserFileDao) IsFileOrDirExist(ctx context.Context, userID, name, path string) (bool, error) {
	filter := bson.D{
		{"user_id", userID},
		{"path", path},
		{"name", name},
	}
	res := dao.Database.Collection(collUserFiles).FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		return false, nil
	}
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

func (dao *UserFileDao) ReplaceFileOrDir(ctx context.Context, po *UserFilePO) (string, error) {
	filter := bson.D{
		{"user_id", po.UserID},
		{"path", po.Path},
		{"name", po.Name},
	}
	res, err := dao.Database.Collection(collUserFiles).ReplaceOne(ctx, filter, po)
	if err != nil {
		return "", err
	}
	return res.UpsertedID.(primitive.ObjectID).String(), err
}

func (dao *UserFileDao) UpdateFileOrDir(ctx context.Context, id string, updatePO *UserFilePO) error {
	oid, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return err
	}
	update := bson.D{
		{"$set", updatePO},
	}
	_, err = dao.Database.Collection(collUserFiles).UpdateByID(ctx, oid, update)
	if err != nil {
		return err
	}
	return nil
}

func (dao *UserFileDao) UpdatesFileOrDir(ctx context.Context, ids []string, updatePO *UserFilePO) (int, error) {
	var oids []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, err
		}
		oids = append(oids, oid)
	}
	filter := bson.D{
		{"_id", bson.E{"$in", oids}},
	}
	update := bson.D{
		{"$set", updatePO},
	}
	res, err := dao.Database.Collection(collUserFiles).UpdateMany(ctx, filter, update)
	if err != nil {
		return 0, err
	}
	return int(res.ModifiedCount), err
}

func (dao *UserFileDao) DeleteFileOrDir(ctx context.Context, id string) error {
	_, err := dao.Database.Collection(collUserFiles).DeleteOne(ctx, bson.D{{"_id", id}})
	if err != nil {
		return err
	}
	return nil
}

func (dao *UserFileDao) QueryDocByIDs(ctx context.Context, ids []string) ([]*UserFilePO, error) {
	var content []*UserFilePO
	var oids []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return nil, err
		}
		oids = append(oids, oid)
	}
	filter := bson.D{
		{"_id", bson.E{"$in", oids}},
	}
	cursor, err := dao.Database.Collection(collUserFiles).Find(ctx, filter)
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
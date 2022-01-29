package repo

import (
	"context"
	"log"
	"time"

	"github.com/changpro/disk-service/infra/constants"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
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

func (dao *UserFileDao) MakeNewFolder(ctx context.Context, po *UserFilePO) (string, error) {
	// add a placeholder(fake file) to sync path
	_, err := dao.addFakeFile(ctx, po.UserID, po.Path, po.Name)
	if err != nil {
		return "", err
	}
	res, err := dao.Database.Collection(collUserFiles).InsertOne(ctx, po)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).String(), nil
}

func (dao *UserFileDao) AddFileOrDir(ctx context.Context, po *UserFilePO) (string, error) {
	res, err := dao.Database.Collection(collUserFiles).InsertOne(ctx, po)
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).String(), nil
}

// func (dao *UserFileDao) removeFromParentSubIDs(ctx context.Context, id string, po *UserFilePO) error {
// 	if po.Path == "/" {
// 		return nil
// 	}
// 	path, name := getParentPathAndName(po.Path)
// 	filter := bson.D{
// 		{"user_id", po.UserID},
// 		{"path", path},
// 		{"name", name},
// 	}
// 	update := bson.D{
// 		{"$pull", bson.E{"sub_ids", id}},
// 	}
// 	res := dao.Database.Collection(collUserFiles).FindOneAndUpdate(ctx, filter, update)
// 	if res.Err() != nil {
// 		return res.Err()
// 	}
// 	return nil
// }

// func (dao *UserFileDao) addToParentSubIDs(ctx context.Context, id string, po *UserFilePO) error {
// 	if po.Path == "/" {
// 		return nil
// 	}
// 	path, name := getParentPathAndName(po.Path)
// 	filter := bson.D{
// 		{"user_id", po.UserID},
// 		{"path", path},
// 		{"name", name},
// 	}
// 	update := bson.D{
// 		{"$push", bson.E{"sub_ids", id}},
// 	}
// 	res := dao.Database.Collection(collUserFiles).FindOneAndUpdate(ctx, filter, update)
// 	if res.Err() != nil {
// 		return res.Err()
// 	}
// 	return nil
// }

func (dao *UserFileDao) addFakeFile(ctx context.Context, userID, path, name string) (string, error) {
	res, err := dao.Database.Collection(collUserFiles).
		InsertOne(ctx, &UserFilePO{
			UserID:   userID,
			Status:   5,
			Path:     path + name + "/",
			IsDir:    2,
			CreateAt: time.Now(),
			UpdateAt: time.Now(),
		})
	if err != nil {
		return "", err
	}
	return res.InsertedID.(primitive.ObjectID).String(), nil
}

func (dao *UserFileDao) IsPathExist(ctx context.Context, userID, path string) (bool, error) {
	filter := bson.M{
		"user_id": userID,
		"path":    path,
		"status":  bson.M{"$in": []int32{1, 5}}, // 1: enable, 5: placeholder used to check path
	}
	res := dao.Database.Collection(collUserFiles).FindOne(ctx, filter)
	if res.Err() == mongo.ErrNoDocuments {
		log.Println(res.Err())
		return false, nil
	}
	if res.Err() != nil {
		return false, res.Err()
	}
	return true, nil
}

func (dao *UserFileDao) QueryUserRoot(ctx context.Context, userID string, showHide bool) ([]*UserFilePO, error) {
	var content []*UserFilePO
	filter := bson.M{
		"user_id": userID,
		"path":    "/",
	}
	if !showHide {
		filter["is_hide"] = constants.FileDisplayStatusShow
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
	filter := bson.M{
		"_id": oid,
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
	filter := bson.M{
		"user_id": userID,
		"path":    path,
		"status":  1,
	}
	if !showHide {
		filter["is_hide"] = constants.FileDisplayStatusShow
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
	filter := bson.M{
		"user_id": userID,
		"path":    path,
		"name":    name,
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
	filter := bson.M{
		"user_id": po.UserID,
		"path":    po.Path,
		"name":    po.Name,
	}
	opts := &options.ReplaceOptions{}
	opts.SetUpsert(true)
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

func (dao *UserFileDao) UpdateFileOrDirByIDs(ctx context.Context, ids []string, updatePO *UserFilePO) (int, error) {
	var oids []primitive.ObjectID
	for _, id := range ids {
		oid, err := primitive.ObjectIDFromHex(id)
		if err != nil {
			return 0, err
		}
		oids = append(oids, oid)
	}
	filter := bson.M{
		"_id": bson.M{"$in": oids},
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
	_, err := dao.Database.Collection(collUserFiles).DeleteOne(ctx, bson.M{"_id": id})
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
	filter := bson.M{
		"_id": bson.M{"$in": oids},
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

// func getParentPathAndName(path string) (string, string) {
// 	if path == "/" {
// 		return "", ""
// 	}
// 	p := strings.Split(path, "/")
// 	if len(p) == 0 {
// 		return "", ""
// 	}
// 	return strings.Join(p[:len(p)-2], "/") + "/", p[len(p)-2]
// }

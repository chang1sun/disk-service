package repo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

const collRecycle = "recycle"

type RecycleFileDao struct {
	Database *mongo.Database
}

var recycleFileDao *RecycleFileDao

func GetRecycleFileDao() *RecycleFileDao {
	return recycleFileDao
}

func SetRecycleFileDao(dao *RecycleFileDao) {
	recycleFileDao = dao
}

func (dao *RecycleFileDao) GetRecycleBinList(ctx context.Context, userID string, offset, limit int32) ([]*RecycleFilePO, error) {
	var res []*RecycleFilePO
	filter := &bson.D{
		{"user_id", userID},
	}
	opts := options.FindOptions{}
	opts.SetSkip(int64(offset))
	opts.SetLimit(int64(limit))
	cursor, err := dao.Database.Collection(collRecycle).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	err = cursor.All(ctx, &res)
	if err != nil {
		return nil, err
	}
	return res, nil
}

func (dao *RecycleFileDao) InsertToRecycleBin(ctx context.Context, userID string, pos []*RecycleFilePO) error {
	var docs []interface{}
	for _, po := range pos {
		docs = append(docs, po)
	}
	_, err := dao.Database.Collection(collRecycle).InsertMany(ctx, docs)
	if err != nil {
		return err
	}
	return nil
}

func (dao *RecycleFileDao) DeleteDocs(ctx context.Context, userID string, ids []string) error {
	var oids []primitive.ObjectID
	for _, id := range ids {
		oid, _ := primitive.ObjectIDFromHex(id)
		oids = append(oids, oid)
	}
	filter := bson.D{
		{"id", bson.E{"$in", oids}},
	}
	_, err := dao.Database.Collection(collRecycle).DeleteMany(ctx, filter)
	if err != nil {
		return err
	}
	return nil
}

package util

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMongodbConn(addr, database string) (*mongo.Database, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	client, err := mongo.Connect(ctx, options.Client().ApplyURI(addr).SetMaxPoolSize(5))
	if err != nil {
		return nil, err
	}
	db := client.Database(database)
	return db, nil
}

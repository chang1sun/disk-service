package util

import (
	"github.com/go-redis/redis/v8"
)

func GetRedisConn(addr, user, pw string, db int) *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: pw, // no password set
		DB:       db, // use default DB
	})
}

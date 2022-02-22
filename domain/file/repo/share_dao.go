package repo

import (
	"context"
	"time"

	"github.com/changpro/disk-service/infra/constants"
	"github.com/go-redis/redis/v8"
)

type ShareDao struct {
	Database *redis.Client
}

var shareDao *ShareDao

func GetShareDao() *ShareDao {
	return shareDao
}

func SetShareDao(dao *ShareDao) {
	shareDao = dao
}

// init a share token
func (dao *ShareDao) CreateShareToken(ctx context.Context, token string, po *ShareDetailPO) error {
	p := dao.Database.Pipeline()
	if err := p.HMSet(ctx, token,
		"uploader", po.Uploader,
		"password", po.Password,
		"uniFileId", po.UniFileID,
		"docId", po.DocID,
		"docName", po.DocName,
		"docSize", po.DocSize,
		"isDir", po.IsDir,
		"docType", po.DocType,
		"createTime", time.Now().Format(constants.StandardTimeFormat),
		"expireHours", po.ExpireHours,
		"viewNum", 0,
		"saveNum", 0).Err(); err != nil {
		return err
	}
	if err := p.Expire(ctx, token, time.Duration(po.ExpireHours)*time.Hour).Err(); err != nil {
		return err
	}
	_, err := p.Exec(ctx)
	if err != nil {
		return err
	}
	return nil
}

func (dao *ShareDao) GetShareDetail(ctx context.Context, token string) (*ShareDetailPO, error) {
	var po ShareDetailPO
	err := dao.Database.HGetAll(ctx, token).Scan(&po)
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	if po.DocID == "" {
		return nil, nil
	}
	return &po, nil
}

func (dao *ShareDao) IncrViewNum(ctx context.Context, token string) error {
	err := dao.Database.HIncrBy(ctx, token, "viewNum", 1).Err()
	if err != nil {
		return err
	}
	return nil
}

func (dao *ShareDao) IncrSaveNum(ctx context.Context, token string) error {
	err := dao.Database.HIncrBy(ctx, token, "saveNum", 1).Err()
	if err != nil {
		return err
	}
	return nil
}

func (dao *ShareDao) DeleteShare(ctx context.Context, token string) error {
	if res := dao.Database.Del(ctx, token); res.Err() != nil {
		return res.Err()
	}
	return nil
}

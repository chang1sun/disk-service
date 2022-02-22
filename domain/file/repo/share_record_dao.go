package repo

import (
	"context"
	"time"

	"github.com/changpro/disk-service/infra/constants"
	"gorm.io/gorm"
)

type ShareRecordDao struct {
	Database *gorm.DB
}

var shareRecordDao *ShareRecordDao

func GetShareRecordDao() *ShareRecordDao {
	return shareRecordDao
}

func SetShareRecordDao(dao *ShareRecordDao) {
	shareRecordDao = dao
}

func (dao *ShareRecordDao) CreateShareRecord(ctx context.Context, po *ShareRecordPO) error {
	if err := dao.Database.WithContext(ctx).Create(po).Error; err != nil {
		return err
	}
	return nil
}

func (dao *ShareRecordDao) QueryRecordList(ctx context.Context, query *RecordQuery) ([]*ShareRecordPO, int64, error) {
	var count int64
	cond := dao.Database.WithContext(ctx).Model(&ShareRecordPO{}).Where("user_id = ?", query.UserID)
	if query.Type != 0 {
		cond.Where("type = ?", query.Type)
	}
	if query.StartTime != 0 && query.EndTime != 0 {
		cond.Where("create_time > ? and create_time < ?",
			time.Unix(query.StartTime, 0).Format(constants.StandardTimeFormat),
			time.Unix(query.EndTime, 0).Format(constants.StandardTimeFormat))
	}
	if err := cond.Count(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	var list []*ShareRecordPO
	if err := cond.Offset(int(query.Offset)).Limit(int(query.Limit)).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

func (dao *ShareRecordDao) DeleteShareRecord(ctx context.Context, token string) error {
	if err := dao.Database.WithContext(ctx).Model(&ShareRecordPO{}).
		Where("token = ?", token).Updates(&ShareRecordPO{
		Status: 2,
	}).Error; err != nil {
		return err
	}
	return nil
}

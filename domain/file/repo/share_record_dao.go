package repo

import (
	"context"

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

func (dao *ShareRecordDao) QueryRecordList(ctx context.Context, userID string, start int32,
	limit int32) ([]*ShareRecordPO, int64, error) {
	var count int64
	if err := dao.Database.WithContext(ctx).Where("user_id = ?", userID).Count(&count).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, 0, nil
		}
		return nil, 0, err
	}
	var list []*ShareRecordPO
	if err := dao.Database.WithContext(ctx).Offset(int(start)).Limit(int(limit)).
		Where("user_id = ?", userID).Order("id DESC").Find(&list).Error; err != nil {
		return nil, 0, err
	}
	return list, count, nil
}

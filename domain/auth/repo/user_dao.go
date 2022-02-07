package repo

import (
	"context"
	"time"

	"github.com/changpro/disk-service/infra/config"
	"gorm.io/gorm"
)

type UserDao struct {
	Database *gorm.DB
}

var userDaoImpl *UserDao

func GetUserDao() *UserDao {
	return userDaoImpl
}

func SetUserDao(dao *UserDao) {
	userDaoImpl = dao
}

func (dao *UserDao) QueryUserByID(ctx context.Context, userID string) (*UserPO, error) {
	var po UserPO
	err := dao.Database.WithContext(ctx).Where("user_id = ? and status != ?", userID,
		3).First(&po).Error
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (dao *UserDao) RegisterNewUser(ctx context.Context, userPO *UserPO) error {
	userPO.Status = 1
	userPO.CreateTime, userPO.UpdateTime = time.Now(), time.Now()
	err := dao.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		userPO.CreateTime = time.Now()
		userPO.UpdateTime = time.Now()
		userPO.LastLogin = time.Now()
		userPO.Status = 1
		err := tx.WithContext(ctx).Create(userPO).Error
		if err != nil {
			return err
		}
		// insert a record into table user_analysis
		err = tx.WithContext(ctx).Create(&UserAnalysisPO{
			UserID:        userPO.UserID,
			TotalSize:     config.GetConfig().InitUserStorageSize,
			UsedSize:      0,
			FileNum:       0,
			UploadFileNum: 0,
			CreateTime:    time.Now(),
			UpdateTime:    time.Now(),
		}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return err
	}
	return nil
}

func (dao *UserDao) SignIn(ctx context.Context, userID string, pwMask string) (*UserPO, error) {
	var po UserPO
	err := dao.Database.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
		// find user
		err := tx.WithContext(ctx).Where("user_id = ? and user_pw = ? and status != ?", userID,
			pwMask, 3).First(&po).Error
		if err != nil {
			return err
		}
		// update last_login
		err = tx.WithContext(ctx).Where("user_id = ?", userID).Updates(&UserPO{LastLogin: time.Now()}).Error
		if err != nil {
			return err
		}
		return nil
	})
	if err == gorm.ErrRecordNotFound {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (dao *UserDao) UpdatePassword(ctx context.Context, userID string, newPw string) error {
	if err := dao.Database.WithContext(ctx).
		Where("user_id = ? and status != ?",
			userID, 3).
		Updates(&UserPO{UserPW: newPw, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

func (dao *UserDao) UpdateUserProfile(ctx context.Context, dto *ModifyUserProfileDTO) error {
	if err := dao.Database.WithContext(ctx).
		Where("user_id = ? and status != ?", dto.UserID, 3).
		Updates(&UserPO{AuthEmail: dto.AuthEmail, UserIcon: dto.Icon, UpdateTime: time.Now()}).Error; err != nil {
		return err
	}
	return nil
}

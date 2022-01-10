package repo

import (
	"context"

	"github.com/changpro/disk-service/auth/common/constants"
	"gorm.io/gorm"
)

type UserAnalysisDao struct {
	Database *gorm.DB
}

var UserAnalysisDaoImpl *UserAnalysisDao

func GetUserAnalysisDao() *UserAnalysisDao {
	return UserAnalysisDaoImpl
}

func SetUserAnalysisDao(dao *UserAnalysisDao) {
	UserAnalysisDaoImpl = dao
}

func (dao *UserAnalysisDao) QueryUserAnalysisByUserID(ctx context.Context, userID string) (*UserAnalysisPO, error) {
	var po UserAnalysisPO
	err := dao.Database.WithContext(ctx).Table("user_analysis as ua").Select("ua.*").
		Joins("left join user_info as ui on ua.user_id = ui.user_id").
		Where("ua.user_id = ? and ui.status != ?", userID, constants.UserStatusCancelled).
		First(&po).Error
	if err != nil {
		return nil, err
	}
	return &po, nil
}

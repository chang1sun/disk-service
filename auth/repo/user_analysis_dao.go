package repo

import (
	"context"
	"time"

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
		Where("ua.user_id = ? and ui.status != ?", userID, 3).
		First(&po).Error
	if err != nil {
		return nil, err
	}
	return &po, nil
}

func (dao *UserAnalysisDao) UpdateUserStorage(ctx context.Context, dto *UpdateUserAnalysisDTO) error {
	if err := dao.Database.WithContext(ctx).Model(&UserAnalysisPO{}).Where("user_id = ?", dto.UserID).Updates(
		map[string]interface{}{
			"file_num":        gorm.Expr("file_num + ?", dto.FileNum),
			"file_upload_num": gorm.Expr("file_upload_num + ?", dto.UploadFileNum),
			"used_size":       gorm.Expr("used_size + ?", dto.Size),
			"update_time":     time.Now(),
		},
	).Error; err != nil {
		return err
	}
	return nil
}

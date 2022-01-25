package repo

import (
	"time"
)

type UserAnalysisPO struct {
	ID            int32     `gorm:"column:id;type:int;primaryKey;autoIncrement"`
	UserID        string    `gorm:"column:user_id;type:varchar(20);not null"`
	FileNum       int32     `gorm:"column:file_num;type:int;default:0;not null"`
	UploadFileNum int32     `gorm:"column:file_upload_num;type:int;default:0;not null"`
	TotalSize     int64     `gorm:"column:total_size;type:bigint(20);not null;default:0"`
	UsedSize      int64     `gorm:"column:used_size;type:bigint(20);not null;default:0"`
	CreateTime    time.Time `gorm:"column:create_time;type:timestamp"`
	UpdateTime    time.Time `gorm:"column:update_time;type:timestamp"`
}

func (UserAnalysisPO) TableName() string {
	return "user_analysis"
}

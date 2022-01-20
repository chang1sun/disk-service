package repo

import "time"

type ShareRecordPO struct {
	ID         int32     `gorm:"column:id;type:int;primaryKey;autoIncrement"`
	UserID     string    `gorm:"column:user_id;type:varchar(20);not null"`
	Message    string    `gorm:"column:message;type:text;not null"`
	DocName    string    `gorm:"column:doc_name;type:varchar(100);not null"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null"`
}

func (ShareRecordPO) TableName() string {
	return "share_record"
}

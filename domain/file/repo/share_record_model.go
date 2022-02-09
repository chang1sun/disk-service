package repo

import "time"

type ShareRecordPO struct {
	ID         int32     `gorm:"column:id;type:int;primaryKey;autoIncrement"`
	UserID     string    `gorm:"column:user_id;type:varchar(20);not null"`
	Token      string    `gorm:"column:token;type:char(32);not null"`
	DocID      string    `gorm:"column:doc_id;type:char(24);not null"`
	DocName    string    `gorm:"column:doc_name;type:varchar(100);not null"`
	ExpireTime time.Time `gorm:"column:expire_time;type:timestamp;not null"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp;not null"`
	Type       int32     `gorm:"column:type;type:tinyint(1);not null"`
}

func (ShareRecordPO) TableName() string {
	return "share_record"
}

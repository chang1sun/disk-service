package repo

import (
	"time"
)

type UserPO struct {
	ID         int32     `gorm:"column:id;type:int;primaryKey;autoIncrement"`
	UserID     string    `gorm:"column:user_id;type:varchar(20);not null"`
	UserPW     string    `gorm:"column:user_pw;type:varchar(100);not null"`
	UserIcon   string    `gorm:"column:user_icon;type:varchar(255)"`
	Status     int32     `gorm:"column:status;type:tinyint(1);not null;default:1"`
	CreateTime time.Time `gorm:"column:create_time;type:timestamp"`
	UpdateTime time.Time `gorm:"column:update_time;type:timestamp"`
	AuthEmail  string    `gorm:"column:auth_email;type:varchar(30);not null"`
	LastLogin  time.Time `gorm:"column:last_login;type:timestamp"`
}

func (UserPO) TableName() string {
	return "user_info"
}

package util

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var gormConn *gorm.DB

func GetGormConn() *gorm.DB {
	return gormConn
}

func InitGormConn() error {
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		GetConfig().Mysql.User,
		GetConfig().Mysql.Password,
		GetConfig().Mysql.Addr,
		GetConfig().Mysql.Database)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return err
	}
	gormConn = conn
	return nil
}

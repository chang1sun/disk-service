package util

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func GetGormConn(user, pw, addr, database string) (*gorm.DB, error) {
	dsn := fmt.Sprintf("%v:%v@tcp(%v)/%v?charset=utf8mb4&parseTime=True&loc=Local",
		user,
		pw,
		addr,
		database)
	conn, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return conn, nil
}

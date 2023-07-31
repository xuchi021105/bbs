package utils

import (
	"log"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	host     string
	port     string
	user     string
	password string
	dbName   string
)

func GetDB() *gorm.DB {
	return DB
}

func InitAndLoadDB() {

	var err error

	host = `127.0.0.1`
	port = `3306`
	user = `your user`
	password = `your password` // 这里的密码和用户名请用自己的密码和用户名
	dbName = `your database name`

	dsn := user + ":" + password + "@tcp(" + host + ":" + port + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		log.Fatalf("数据库连接错误: %+v", err)
	}

}

package store

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log-record/utils"
)

var DB *gorm.DB

func init() {
	dbUrl := utils.EnvGetOrDefaultStringValue("DB_URL", "root:123456@tcp(10.43.33.6:3306)/test?charset=utf8mb4&parseTime=True&loc=Local")
	db, err := gorm.Open(mysql.Open(dbUrl), &gorm.Config{
		//Logger: logger.New(log.New(os.Stdout, "\r\n", log.LstdFlags), logger.Config{LogLevel: logger.Info}),
	})
	if err != nil {
		panic("failed to connect database")
	}
	fmt.Printf("DB_URL %v \r\n", dbUrl)
	DB = db
}

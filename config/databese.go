package config

import (
	"log"
	"main/model"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	PageSize = 10
)

func InitDB() *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{PrepareStmt: true})
	if err != nil {
		log.Fatalln("failed to connect to database:", err)
	}
	err = DB.AutoMigrate(&model.User{}, &model.File{}) // 自动迁移 User 模型
	if err != nil {
		log.Fatalln("Table User failed Migrate")
	}
	return DB
}

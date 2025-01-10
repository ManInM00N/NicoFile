package config

import (
	"main/model"
	. "main/util"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDB() {
	var err error
	DB, err = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{})
	if err != nil {
		ErrorLog.Println("failed to connect to database:", err)
		panic(err)
	}
	err = DB.AutoMigrate(&model.User{}) // 自动迁移 User 模型
	if err != nil {
		ErrorLog.Println("Table User failed Migrate")
	}
}

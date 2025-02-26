package config

import (
	"gorm.io/gorm/logger"
	"main/model"
	"main/pkg/util"
	"time"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	PageSize = 10
)

func InitDB() *gorm.DB {
	var err error
	DB, err = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Discard,
	})
	// 获取底层 SQL DB 实例
	sqlDB, _ := DB.DB()
	// 设置连接池
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间
	if err != nil {
		util.Log.Fatalln("failed to connect to database:", err)
	}
	err = DB.AutoMigrate(
		&model.User{},
		&model.File{},
		&model.Chunk{},
	) // 自动迁移 User 模型
	if err != nil {
		util.Log.Fatalln("Table User failed Migrate")
	}
	return DB
}

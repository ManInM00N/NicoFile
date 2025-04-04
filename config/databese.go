package config

import (
	"fmt"
	"gorm.io/gorm/logger"
	"main/model"
	"main/pkg/encrypt"
	"main/pkg/util"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var (
	DB       *gorm.DB
	PageSize = 10
	username = "root"      //账号
	password = "root"      //密码
	host     = "127.0.0.1" //数据库地址，可以是Ip或者域名
	port     = 23306       //数据库端口
	Dbname   = "nicofile"  //数据库名
	timeout  = "10s"       //连接超时，10秒
)

func GetDB() *gorm.DB {
	return DB
}
func InitDB(host string) *gorm.DB {
	var err error
	dsn := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local&timeout=%s", username, password, host, port, Dbname, timeout)
	DB, err = gorm.Open(mysql.Open(dsn), &gorm.Config{
		PrepareStmt: true,
		Logger:      logger.Discard,
	})
	if err != nil {
		util.Log.Fatalln("failed to connect to database:", err)
	}
	//DB, err = gorm.Open(sqlite.Open("./data.db"), &gorm.Config{
	//	PrepareStmt: true,
	//	Logger:      logger.Discard,
	//})
	// 获取底层 SQL DB 实例
	sqlDB, _ := DB.DB()
	// 设置连接池
	sqlDB.SetMaxOpenConns(100)          // 最大连接数
	sqlDB.SetMaxIdleConns(10)           // 最大空闲连接数
	sqlDB.SetConnMaxLifetime(time.Hour) // 连接最大存活时间

	err = DB.AutoMigrate(
		&model.User{},
		&model.File{},
		&model.Chunk{},
		&model.Article{},
		&model.Image{},
		&model.Comment{},
	) // 自动迁移 User 模型
	if err != nil {
		util.Log.Fatalln("Table User failed Migrate")
	}
	pwd := encrypt.EncPassword("123456")
	DB.Model(&model.User{}).Where("username = ?", "admin").FirstOrCreate(&model.User{Username: "admin", Password: pwd, Priority: 2})
	return DB
}

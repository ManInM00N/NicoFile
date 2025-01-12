package svc

import (
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
	config2 "main/config"
	"main/nicofile/internal/config"
	"main/nicofile/internal/middleware"
)

type ServiceContext struct {
	Config              config.Config
	DB                  *gorm.DB
	UserExistMiddleware rest.Middleware
}

func NewServiceContext(c config.Config) *ServiceContext {
	return &ServiceContext{
		Config:              c,
		DB:                  config2.InitDB(),
		UserExistMiddleware: middleware.NewUserExistMiddleware().Handle,
	}
}

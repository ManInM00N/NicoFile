package svc

import (
	"fmt"
	"github.com/IBM/sarama"
	"github.com/redis/go-redis/v9"
	"github.com/zeromicro/go-zero/rest"
	"gorm.io/gorm"
	config2 "main/config"
	"main/kafka"
	"main/nicofile/grpc"
	"main/nicofile/internal/config"
	"main/nicofile/internal/middleware"
	CacheRedis "main/redis"
)

type ServiceContext struct {
	Config              config.Config
	DB                  *gorm.DB
	Rdb                 *redis.Client
	Producer            *sarama.AsyncProducer
	UserExistMiddleware rest.Middleware
	HotArticlePool      *grpc.Pool
}

func NewServiceContext(c config.Config) *ServiceContext {
	pool, err := grpc.NewPool(fmt.Sprintf("%s%d", c.Services.ArticleRank.Host, c.Services.ArticleRank.Port), c.GrpcPool.Size)
	if err != nil {
		panic(err)
	}
	return &ServiceContext{
		Config:              c,
		DB:                  config2.InitDB(),
		Rdb:                 CacheRedis.InitRedis(c.Redis.Host, c.Redis.Port, c.Redis.Password, c.Redis.DB, c.Redis.Disabled),
		Producer:            kafka.Subscribe(c.Kafka.Disabled, c.Kafka.Host, c.Kafka.Port),
		UserExistMiddleware: middleware.NewUserExistMiddleware().Handle,
		HotArticlePool:      pool,
	}
}
func (s *ServiceContext) Close() {
	if s.HotArticlePool != nil {
		s.HotArticlePool.Close()
	}
}

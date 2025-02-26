package config

import "github.com/zeromicro/go-zero/rest"

type Config struct {
	rest.RestConf
	ChunkStorePath string
	StoragePath    string
	Auth           struct {
		AccessSecret string
		AccessExpire int64
	}
	Redis struct {
		Host            string
		Port            int
		Password        string
		DB              int
		Disabled        bool
		Expire          int64
		RefreshInterval int64
	}
	Kafka struct {
		Host     string
		Port     int
		Topic    string
		Broker   string
		Disabled bool
	}
}

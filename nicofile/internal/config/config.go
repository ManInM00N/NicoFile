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
}

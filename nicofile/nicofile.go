package main

import (
	"flag"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
	"main/nicofile/internal/middleware"
	"main/pkg/util"
	CacheRedis "main/redis"
	"net/http"
	"os"
	"path"
	"strings"
	"time"

	"main/nicofile/internal/config"
	"main/nicofile/internal/handler"
	"main/nicofile/internal/svc"

	"github.com/zeromicro/go-zero/core/conf"
	"github.com/zeromicro/go-zero/rest"
)

var configFile = flag.String("f", "nicofile/etc/nicofile-api.yaml", "the config file")
var inTest = flag.Lookup("test.v") != nil

const basename = "/static"

func main() {
	util.NewLog("nicofile/nicofile-log")
	flag.Parse()
	var c config.Config
	conf.MustLoad(*configFile, &c)

	os.MkdirAll(c.ChunkStorePath, os.ModePerm)
	os.MkdirAll(c.StoragePath, os.ModePerm)
	domains := []string{"*"}
	fs := http.Dir("storage")
	fileServer := http.FileServer(fs)
	server := rest.MustNewServer(c.RestConf,
		rest.WithNotFoundHandler(&NotFoundHandler{ // 自定义 NotFoundHandler，对虚拟路由做处理
			fs:         fs,
			fileServer: fileServer,
		}),
		rest.WithFileServer(basename, fs),
		rest.WithCors(domains...),
		rest.WithCustomCors(func(header http.Header) {
			header.Add("Access-Control-Allow-Headers", "Content-Type,AccessToken,Range,X-CSRF-Token,Authorization,Token,X-Token,X-User-Id,OS,Platform, Version")
			header.Set("Access-Control-Allow-Methods", "GET,POST,PUT,DELETE,OPTIONS,PATCH")
			header.Set("Access-Control-Expose-Headers", "Content-Length, Content-Type")
			header.Set("Access-Control-Allow-Origin", "*")
		}, nil, "*"),
	)
	defer server.Stop()

	ctx := svc.NewServiceContext(c)
	if !c.Redis.Disabled {
		CacheRedis.PullData(ctx.Rdb, ctx.DB)
		go func(rdb *redis.Client, DB *gorm.DB) {
			timer := time.NewTimer(time.Duration(c.Redis.RefreshInterval) * time.Second)
			for range timer.C {
				CacheRedis.Transport(ctx.Rdb, ctx.DB)
			}
		}(ctx.Rdb, ctx.DB)
		defer CacheRedis.Transport(ctx.Rdb, ctx.DB)
	}
	if !c.Kafka.Disabled {
		defer func() {
			if err := (*ctx.Producer).Close(); err != nil {
				util.Log.Errorf("Error closing Kafka producer: %v", err)
			}
		}()
	}
	//if c.DevServer.EnableMetrics {
	server.Use(middleware.PrometheusMiddleware)
	//}
	handler.RegisterHandlers(server, ctx)

	util.Log.Printf("Starting server at %s:%d...\n", c.Host, c.Port)
	server.Start()
}

type NotFoundHandler struct {
	fs         http.FileSystem
	fileServer http.Handler
}

func (n NotFoundHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate, proxy-revalidate")
	w.Header().Set("Pragma", "no-cache")
	w.Header().Set("Expires", "0")
	w.Header().Set("Surrogate-Control", "no-store")
	filePath := strings.TrimPrefix(path.Clean(r.URL.Path), basename)
	if len(filePath) == 0 {
		filePath = basename
	}

	file, err := n.fs.Open(filePath)
	switch {
	case err == nil:
		n.fileServer.ServeHTTP(w, r)
		_ = file.Close()
		return
	case os.IsNotExist(err):
		r.URL.Path = "/" // all virtual routes in react app means visit index.html
		n.fileServer.ServeHTTP(w, r)
		return
	default:
		http.Error(w, "not found", http.StatusNotFound)
		return
	}
}

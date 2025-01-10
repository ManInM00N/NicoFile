package main

import (
	"main/config"
	"main/server/handler"
	"main/server/proto"
	. "main/util"
	"net"

	"github.com/orandin/lumberjackrus"
	"github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

var log = logrus.New()

func newRotateHook() logrus.Hook {
	hook, _ := lumberjackrus.NewHook(
		&lumberjackrus.LogFile{ // 通用日志配置
			Filename:   "general.log",
			MaxSize:    100,
			MaxBackups: 1,
			MaxAge:     1,
			Compress:   false,
			LocalTime:  false,
		},
		logrus.InfoLevel,
		&logrus.TextFormatter{DisableColors: true},
		&lumberjackrus.LogFileOpts{ // 针对不同日志级别的配置
			logrus.TraceLevel: &lumberjackrus.LogFile{
				Filename: "trace.log",
			},
			logrus.ErrorLevel: &lumberjackrus.LogFile{
				Filename:   "error.log",
				MaxSize:    10,    // 日志文件在轮转之前的最大大小，默认 100 MB
				MaxBackups: 10,    // 保留旧日志文件的最大数量
				MaxAge:     10,    // 保留旧日志文件的最大天数
				Compress:   true,  // 是否使用 gzip 对日志文件进行压缩归档
				LocalTime:  false, // 是否使用本地时间，默认 UTC 时间
			},
		},
	)
	return hook
}

func main() {
	// 初始化数据库
	config.InitDB()
	InfoLog.Printf("DB init successfully %d %d", 2, 4)
	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		ErrorLog.Printf("failed to listen: %v", err)
		panic(err)
	}

	grpcServer := grpc.NewServer()

	proto.RegisterAuthServiceServer(grpcServer, &handler.AuthServiceServer{})

	InfoLog.Println("gRPC server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		ErrorLog.Printf("failed to serve: %v", err)
	}
	defer func() {
		grpcServer.Stop()
		lis.Close()
	}()
}

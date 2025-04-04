package main

import (
	"fmt"
	"google.golang.org/grpc"
	"main/config"
	"main/pkg/util"
	CacheRedis "main/redis"
	"main/server/handler"
	"main/server/proto/articleRank"
	"main/server/proto/auth"
	"net"
)

func main() {
	prefix := "server"
	util.NewLog(prefix)
	// 初始化数据库
	config.InitDB("127.0.0.1")
	CacheRedis.InitRedis("127.0.0.1:", 6380, "", 0, false)
	util.Log.Tracef("DB init successfully %d %d", 2, 4)
	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", "127.0.0.1:50051")
	if err != nil {
		util.Log.Errorf("failed to listen: %v\n", err)
		panic(err)
	}

	grpcServer := grpc.NewServer()
	articleRankServer := handler.NewArticleRankService()
	auth.RegisterAuthServiceServer(grpcServer, &handler.AuthServiceServer{})
	articleRank.RegisterArticleRankServiceServer(grpcServer, articleRankServer)

	util.Log.Traceln("gRPC server is running on port 127.0.0.1:50051")
	fmt.Println("gRPC server is running on port 127.0.0.1:50051")
	if err := grpcServer.Serve(lis); err != nil {
		util.Log.Errorf("failed to serve: %v\n", err)
	}
	defer func() {
		grpcServer.Stop()
		lis.Close()
	}()
}

package main

import (
	"google.golang.org/grpc"
	"main/config"
	"main/pkg/util"
	"main/server/handler"
	"main/server/proto/auth"
	"net"
)

func main() {
	prefix := "server"
	util.NewLog(prefix)
	// 初始化数据库
	config.InitDB()
	util.Log.Tracef("DB init successfully %d %d", 2, 4)
	// 启动 gRPC 服务
	lis, err := net.Listen("tcp", ":50051")
	if err != nil {
		util.Log.Errorf("failed to listen: %v\n", err)
		panic(err)
	}

	grpcServer := grpc.NewServer()

	auth.RegisterAuthServiceServer(grpcServer, &handler.AuthServiceServer{})

	util.Log.Traceln("gRPC server is running on port :50051")
	if err := grpcServer.Serve(lis); err != nil {
		util.Log.Errorf("failed to serve: %v\n", err)
	}
	defer func() {
		grpcServer.Stop()
		lis.Close()
	}()
}

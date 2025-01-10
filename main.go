package main

import (
	"context"
	"fmt"
	"main/server/proto"
	. "main/util"

	"google.golang.org/grpc"
)

func main() {
	InfoLog.Println("log init")
	conn, err := grpc.NewClient(":50051", grpc.WithInsecure())
	if err != nil {
		fmt.Println("err")
		ErrorLog.Printf("did not connect: %v", err)
		//panic(err)
		return

	}
	defer conn.Close()

	client := proto.NewAuthServiceClient(conn)

	// 尝试登录
	res, err := client.Login(context.Background(), &proto.LoginRequest{
		Username: "user1",
		Password: "password123",
	})
	if err != nil {
		ErrorLog.Printf("could not login: %v", err)
		panic(err)
	}

	fmt.Printf("Login Response: %v\n", res.Message)
}

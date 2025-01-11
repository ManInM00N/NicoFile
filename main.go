package main

import (
	"context"
	"fmt"
	"google.golang.org/grpc"
	"log"
	"main/server/proto"
)

func _() {

	log.Println("log init")
	conn, err := grpc.NewClient(":50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		panic(err)
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
		log.Printf("could not login: %v", err)
		panic(err)
	}

	fmt.Printf("Login Response: %v\n", res.Message)
}

func main() {

}

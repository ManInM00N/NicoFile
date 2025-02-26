package main

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/proto"
	"log"
	"main/pkg/util"
	"main/server/proto/auth"
	"main/server/proto/kafka"
	"time"
)

var (
	Users = int64(0)
	Files = int64(0)
)

type ConsumerGroupHandler struct{}

func (h *ConsumerGroupHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (h *ConsumerGroupHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	for msg := range claim.Messages() {
		fmt.Println("Received messages", string(msg.Value))
		switch msg.Topic {
		case "data-monitor-test":
			var data kafka.UserMonitor
			err := proto.Unmarshal(msg.Value, &data)
			if err != nil {
				util.Log.Errorf("Failed to unmarshal message: %v (Partition: %d, Offset: %d)\n", err, msg.Partition, msg.Offset)
			}
			session.MarkMessage(msg, "") // 标记消息为已处理

		case "user-monitor":
			var data kafka.UserMonitor
			err := proto.Unmarshal(msg.Value, &data)
			if err != nil {
				util.Log.Errorf("Failed to unmarshal uesr-message: %v (Partition: %d, Offset: %d)\n", err, msg.Partition, msg.Offset)
			} else {
				if data.Message == "A new user has been registered" {
					Users++
				} else {
					Users--
				}
			}
			session.MarkMessage(msg, "") // 标记消息为已处理
		case "file-monitor":
			var data kafka.FileMonitor
			err := proto.Unmarshal(msg.Value, &data)
			if err != nil {
				util.Log.Errorf("Failed to unmarshal file-message: %v (Partition: %d, Offset: %d)\n", err, msg.Partition, msg.Offset)
			} else {
				if data.Message == "A file has been uploaded" {
					Files++
				} else {
					Files--
				}
			}
			session.MarkMessage(msg, "") // 标记消息为已处理
		default:
			session.MarkMessage(msg, "") // 标记消息为已处理
		}

	}
	return nil
}

func _() {

	log.Println("log init")
	conn, err := grpc.NewClient(":50051", grpc.WithInsecure())
	if err != nil {
		log.Printf("did not connect: %v\n", err)
		panic(err)
		return

	}
	defer conn.Close()

	client := auth.NewAuthServiceClient(conn)

	// 尝试登录
	res, err := client.Login(context.Background(), &auth.LoginRequest{
		Username: "user1",
		Password: "password123",
	})
	if err != nil {
		util.Log.Printf("could not login: %v", err)
		panic(err)
	}

	fmt.Printf("Login Response: %v\n", res.Message)
}

func main() {
	util.NewLog("monitor-log")
	util.Log.Println("monitor started")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "data-monitor-test", config)
	if err != nil {
		panic(fmt.Errorf("failed to start Sarama consumer: %v\n", err))
	}
	defer func() {
		if err = consumer.Close(); err != nil {
			panic(fmt.Errorf("Failed to close Sarama consumer: %v\n", err))
		}
	}()
	handler := &ConsumerGroupHandler{}
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	go func() {
		for {
			err := consumer.Consume(ctx, []string{"data-monitor-test", "user-monitor", "file-monitor"}, handler)
			if err != nil {
				panic(fmt.Errorf("Failed to consume data-monitor topic: %v\n", err))
			}
			if ctx.Err() != nil {
				panic(fmt.Errorf("Failed	%v\n", ctx.Err()))
			}
			time.Sleep(time.Millisecond * 50)
		}
	}()
	ticker := time.NewTicker(time.Second * 15)
	for {
		select {
		case <-ticker.C:

			fmt.Printf("Users: %d, Files: %d\n", Users, Files)
			util.Log.Tracef("Users: %d, Files: %d\n", Users, Files)
		default:

		}
	}
}

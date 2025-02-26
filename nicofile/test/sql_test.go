package test

import (
	"context"
	"fmt"
	"github.com/IBM/sarama"
	"github.com/ManInM00N/go-tool/statics"
	"google.golang.org/protobuf/proto"
	"gorm.io/gorm"
	config2 "main/config"
	"main/kafka"
	"main/model"
	"main/pkg/util"
	CacheRedis "main/redis"
	kafka2 "main/server/proto/kafka"
	"testing"
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
		switch msg.Topic {
		case "data-monitor-test":
			fmt.Println("Received messages", string(msg.Value))
			session.MarkMessage(msg, "") // 标记消息为已处理

		case "user-monitor":
			var data kafka2.UserMonitor
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
			var data kafka2.FileMonitor
			err := proto.Unmarshal(msg.Value, &data)
			if err != nil {
				util.Log.Errorf("Failed to unmarshal file-message: %v (Partition: %d, Offset: %d) , data len:%d\n", err, msg.Partition, msg.Offset, len(msg.Value))
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
func TestQuery(t *testing.T) {
	util.NewLog("test")
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	consumer, err := sarama.NewConsumerGroup([]string{"127.0.0.1:9092"}, "test", config)
	if err != nil {
		t.Errorf("failed to start Sarama consumer: %v\n", err)
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
	go func() {
		for {
			select {
			case <-ticker.C:

				t.Logf("Users: %d, Files: %d\n", Users, Files)
			default:

			}
		}
	}()
	pro := kafka.Subscribe(false, "127.0.0.1:", 9092)
	if pro == nil {
		t.Errorf("Subscribe failed\n")
		return
	}
	//pro, err := sarama.NewAsyncProducer([]string{"127.0.0.1:9092"}, config)

	for {
		event := &kafka2.FileMonitor{
			Message: "A file has been uploaded",
			Warning: true,
			UserId:  uint32(3),
		}
		data, err2 := proto.Marshal(event)
		if err2 != nil {
			t.Errorf("unexpected error: %v", err2)
			return
		}
		mess := &sarama.ProducerMessage{
			Topic: "file-monitor",
			Value: sarama.ByteEncoder(data),
		}
		(*pro).Input() <- mess
		time.Sleep(time.Second * 5)
	}
	(*pro).Close()
}
func TestQueryLimit(t *testing.T) {
	DB := config2.InitDB()
	var ed model.Chunk

	for i := 0; i < 10000; i++ {
		var chunk []model.Chunk
		DB.Model(model.Chunk{}).Select("file_name,author_id,chunk_index").Where("author_id = ?  and file_name = 'test4'", 4).Find(&chunk)

		//DB.Model(model.Chunk{}).Where("author_id = ? and file_name = 'test4'", 4).Find(&chunk)
		if len(chunk) > 0 && i == 0 {
			ed = chunk[0]
		}
	}
	if ed.AuthorID == 0 {
		t.Errorf("Query failed %v \n", ed)
	}

}

func TestSqlQuery(t *testing.T) {
	DB := config2.InitDB()
	file := model.File{
		FileName: "118268382_p0_0",
	}
	DB.Model(&model.File{}).Where("file_name = ?", file.FileName).FirstOrCreate(&file)
	t.Logf("file: %v\n", file)
	num := int64(0)
	DB.Model(&model.File{}).Where("file_name = ?", "118268382_p0_0").Count(&num)
	if num == 0 {
		t.Errorf("Query failed num: %d\n", num)
	}
}
func TestSqlUpdate(t *testing.T) {
	DB := config2.InitDB()
	var file model.File
	DB.First(&file, gorm.Expr("id = ?", 111))
	file.DownloadTimes = 3
	file.Description = "test"
	if err := DB.Save(&file).Error; err != nil {
		t.Errorf("Update failed err: %v\n", err)
	}
}

func TestRedis(t *testing.T) {
	rdb := CacheRedis.InitRedis("127.0.0.1", 6380, "", 0, false)
	ctx := context.Background()
	data, err := rdb.HGetAll(ctx, "file:1111312232").Result()
	fmt.Println(data, err, data == nil)
}
func TestTransport(t *testing.T) {

	DB := config2.InitDB()
	rdb := CacheRedis.InitRedis("127.0.0.1", 6380, "", 0, false)
	ctx := context.Background()
	var file = model.File{}
	if err := DB.Model(&model.File{}).Where("id = ?", 111).First(&file).Error; err != nil {
		t.Errorf("Query failed err: %v\n", err)
	}
	rdb.HSet(ctx, "file:"+statics.IntToString(111), map[string]interface{}{
		"description":    "test1ttttt",
		"download_times": 5441,
	})
	rdb.Expire(ctx, "file:"+statics.IntToString(111), time.Second*300)
	timer := time.NewTicker(time.Second * 10)
	tt := 0
	for range timer.C {
		CacheRedis.Transport(rdb, DB)
		tt++
		cnt, _ := rdb.DBSize(ctx).Result()
		if cnt != 1 {
			t.Errorf("Transport failed cnt: %d\n", cnt)
		}
		if tt == 1 {
			timer.Stop()
			break
		}
	}
}

func BenchmarkTransport(b *testing.B) {
	b.StopTimer()
	util.NewLog("test")
	DB := config2.InitDB()
	rdb := CacheRedis.InitRedis("127.0.0.1:", 6380, "", 0, false)
	ctx := context.Background()
	for i := 0; i < 1000000; i++ {
		rdb.HSet(ctx, "file:"+statics.IntToString(i+1000), map[string]interface{}{
			"description":    "test?",
			"download_times": 222,
			"AuId":           i,
			"deleted":        0,
		})
		//rdb.Expire(ctx, "file:"+statics.IntToString(i+1000))
	}
	fmt.Println("create file success")
	start := time.Now()
	fmt.Println("start time: ", start)
	b.StartTimer()
	CacheRedis.Transport(rdb, DB)
	fmt.Println("time: ", time.Since(start))
}

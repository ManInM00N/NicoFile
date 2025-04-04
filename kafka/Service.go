package kafka

import (
	"context"
	"github.com/IBM/sarama"
	"go.etcd.io/etcd/client/v3"
	"main/pkg/util"
	"time"
)

/*
	发送消息
	event := &kafka.TopicMonitor{
		Message: "data-monitor-test",
		Warning: false,
		UserId:  uint32(i),
	}
	data, _ := proto.Marshal(event)
	msg := &sarama.ProducerMessage{
		Topic: "data-monitor-test",
		Value: sarama.ByteEncoder(data),
	}

producer.Input() <- msg
*/
func Subscribe(disabled bool, host string, port string) *sarama.AsyncProducer {
	if disabled {
		return nil
	}
	etcdClient, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"localhost:2379"}, // Etcd 地址
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		util.Log.Fatalf("Failed to connect to Etcd: %v", err)
	}
	defer etcdClient.Close()

	brokerAddr := getKafkaBrokerFromEtcd(etcdClient)
	util.Log.Printf("Using Kafka broker: %s\n", brokerAddr)

	producer, err := createKafkaProducer(brokerAddr)
	if err != nil {
		util.Log.Fatalf("Failed to create Kafka producer: %v", err)
	}

	(*producer).Input() <- &sarama.ProducerMessage{
		Topic: "data-monitor-test",
		Value: sarama.StringEncoder("test"),
	}
	return producer
}

// 从 Etcd 获取 Kafka Broker 地址
func getKafkaBrokerFromEtcd(etcdClient *clientv3.Client) string {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	resp, err := etcdClient.Get(ctx, "/kafka/config/brokers")
	cancel()
	if err != nil {
		util.Log.Fatalf("Failed to get Kafka broker address from Etcd: %v", err)
	}
	if len(resp.Kvs) == 0 {
		util.Log.Fatalf("No Kafka broker address found in Etcd")
	}
	return string(resp.Kvs[0].Value)
}

func createKafkaProducer(brokerAddr string) (*sarama.AsyncProducer, error) {
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Flush.Messages = 1000
	config.Producer.Flush.Bytes = 1024 * 1024          // 1MB
	config.Producer.Flush.Frequency = time.Second * 10 // 30s update
	config.Producer.MaxMessageBytes = 1024 * 1024      // 1MB
	config.Producer.Retry.Max = 5                      // 最大重试次数
	config.Producer.Return.Successes = true            // 接收成功发送的消息
	config.Producer.Return.Errors = true               // 接收发送失败的消息

	producer, err := sarama.NewAsyncProducer([]string{brokerAddr}, config)
	if err != nil {
		panic(err)
	}
	go func() {
		for msg := range producer.Successes() {
			util.Log.Printf("Message sent successfully! Topic: %s, Partition: %d, Offset: %d\n",
				msg.Topic, msg.Partition, msg.Offset)
		}
	}()
	go func() {
		for err := range producer.Errors() {
			util.Log.Errorf("Failed to send message: %v\n", err)
		}
	}()
	return &producer, nil
}

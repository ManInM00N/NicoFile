package kafka

import (
	"github.com/IBM/sarama"
	"main/pkg/util"
	"strconv"
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
func Subscribe(disabled bool, host string, port int) *sarama.AsyncProducer {
	if disabled {
		return nil
	}
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
	producer, err := sarama.NewAsyncProducer([]string{host + strconv.Itoa(port)}, config)
	if err != nil {
		panic(err)
	}
	// 处理成功发送的消息
	go func() {
		for msg := range producer.Successes() {
			util.Log.Printf("Message sent successfully! Topic: %s, Partition: %d, Offset: %d\n",
				msg.Topic, msg.Partition, msg.Offset)
		}
	}()
	// 处理发送失败的消息
	go func() {
		for err := range producer.Errors() {
			util.Log.Errorf("Failed to send message: %v\n", err)
		}
	}()
	producer.Input() <- &sarama.ProducerMessage{
		Topic: "data-monitor-test",
		Value: sarama.StringEncoder("test"),
	}
	return &producer
}

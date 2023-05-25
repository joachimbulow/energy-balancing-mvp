package client

import (
	"context"
	"fmt"
	"math/rand"
	"strings"
	"time"

	"github.com/joachimbulow/pem-energy-balance/src/util"
	"github.com/segmentio/kafka-go"
)

var (
	logger util.Logger
)

type KafkaClient struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewKafkaClient() (*KafkaClient, error) {
	client := &KafkaClient{}
	logger = util.NewLogger("KafkaBroker")

	var conn *kafka.Conn
	var err error
	for i := 0; i < 10; i++ {
		conn, err = kafka.Dial("tcp", util.GetBrokerURL())
		if err == nil {
			break
		}
		if i < 10 {
			logger.ErrorWithMsg(fmt.Sprintf("Failed to connect to Kafka, Retrying in 5-15 (ish) seconds... Try %d/10", i+1), err)
			time.Sleep(time.Duration(rand.Intn(11)+5) * time.Second)
		} else {
			logger.ErrorWithMsg("Could not connect after 10 attempts, aborting mission", err)
			panic(err)
		}

	}
	defer conn.Close()
	return client, err
}

func setupReader(kafkaClient *KafkaClient, topic string, consumerGroupID string) *kafka.Reader {
	kafkaClient.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:        strings.Split(util.GetBrokerURL(), ","),
		GroupID:        consumerGroupID,
		Topic:          topic,
		MinBytes:       1,
		MaxBytes:       10e6, // 10MB
		StartOffset:    kafka.LastOffset,
		CommitInterval: time.Duration(util.GetKafkaOffSetCommitIntervalMillis()) * time.Millisecond,
	})
	return kafkaClient.reader
}

func setupWriter(kafkaClient *KafkaClient) *kafka.Writer {
	kafkaClient.writer = &kafka.Writer{
		Addr:                   kafka.TCP(util.GetBrokerURL()),
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: false,
	}
	return kafkaClient.writer
}

func (k *KafkaClient) Publish(topic string, key string, message string) error {
	if k.writer == nil {
		k.writer = setupWriter(k)
	}

	err := k.writer.WriteMessages(context.Background(),
		kafka.Message{
			Topic: topic,
			Key:   []byte(key),
			Value: []byte(message),
		},
	)
	if err != nil {
		logger.ErrorWithMsg("Failed to write messages", err)
	}
	return err
}

func (k *KafkaClient) Listen(topic string, consumerGroupID string, handler func(params ...[]byte)) error {
	if k.reader == nil {
		setupReader(k, topic, consumerGroupID)
	}
	for {
		m, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			logger.ErrorWithMsg("Failed to read message from Kafka", err)
			continue
		}

		handler(m.Key, m.Value)
	}
}

func (k *KafkaClient) Connect() error {
	// Kafka connection is established when the reader or writer is created
	return nil
}

func (k *KafkaClient) Disconnect() error {
	// Kafka connections are closed after the reader or writer is used
	return nil
}

func (k *KafkaClient) Subscribe(topic string) error {
	// Kafka reader subscribes to the topic when it is listening
	return nil
}

func (k *KafkaClient) Unsubscribe(topic string) error {
	if err := k.reader.Close(); err != nil {
		logger.ErrorWithMsg("Failed to close reader", err)
	}
	return nil
}

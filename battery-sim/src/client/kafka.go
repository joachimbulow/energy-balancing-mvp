package client

import (
	"context"
	"strings"

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
	conn, err := kafka.Dial("tcp", brokerURL)
	if err != nil {
		logger.ErrorWithMsg("Failed to connect to Kafka:", err)
		return client, err
	}
	defer conn.Close()
	return client, nil
}

func setupReader(kafkaClient *KafkaClient, topic string) *kafka.Reader {
	kafkaClient.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(brokerURL, ","),
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return kafkaClient.reader
}

func setupWriter(kafkaClient *KafkaClient) *kafka.Writer {
	kafkaClient.writer = &kafka.Writer{
		Addr:                   kafka.TCP(brokerURL),
		Balancer:               &kafka.Hash{},
		AllowAutoTopicCreation: true,
	}
	return kafkaClient.writer
}

func (k *KafkaClient) Connect() error {
	// Kafka connection is established when the reader or writer is created
	return nil
}

func (k *KafkaClient) Disconnect() error {
	// Kafka connections are closed after the reader or writer is used
	return nil
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

func (k *KafkaClient) Listen(topic string, handler func(params ...[]byte)) error {
	if k.reader == nil {
		setupReader(k, topic)
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

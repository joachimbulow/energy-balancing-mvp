package broker

import (
	"context"
	"log"
	"strings"

	"github.com/segmentio/kafka-go"
)

var (
	brokerURL = "20.105.75.161:9092"
)

type KafkaBroker struct {
	reader *kafka.Reader
	writer *kafka.Writer
}

func NewKafkaBroker() (*KafkaBroker, error) {
	broker := &KafkaBroker{}
	return broker, nil
}

func setupReader(kafkaBroker *KafkaBroker, topic string) *kafka.Reader {
	kafkaBroker.reader = kafka.NewReader(kafka.ReaderConfig{
		Brokers:  strings.Split(brokerURL, ","),
		Topic:    topic,
		MinBytes: 10e3, // 10KB
		MaxBytes: 10e6, // 10MB
	})
	return kafkaBroker.reader
}

func setupWriter(kafkaBroker *KafkaBroker) *kafka.Writer {
	kafkaBroker.writer = &kafka.Writer{
		Addr:     kafka.TCP(brokerURL),
		Balancer: &kafka.Hash{},
	}
	return kafkaBroker.writer
}

func (k *KafkaBroker) Connect() error {
	// Kafka connection is established when the reader or writer is created
	return nil
}

func (k *KafkaBroker) Disconnect() error {
	// Kafka connections are closed after the reader or writer is used
	return nil
}

func (k *KafkaBroker) Publish(topic string, key string, message string) error {
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
		log.Fatal("failed to write messages:", err)
	}
	if err := k.writer.Close(); err != nil {
		log.Fatal("failed to close writer:", err)
	}
	return err
}

func (k *KafkaBroker) Subscribe(topic string) error {
	// Kafka reader subscribes to the topic when it is listening
	return nil
}

func (k *KafkaBroker) Unsubscribe(topic string) error {
	if err := k.reader.Close(); err != nil {
		log.Fatal("failed to close reader:", err)
	}
	return nil
}

func (k *KafkaBroker) Listen(topic string, handler func(params ...[]byte)) error {
	setupReader(k, topic)
	for {
		// read messages from the Kafka topic
		m, err := k.reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("error reading message from Kafka: %v", err)
			continue
		}

		// call the handler function for each message
		handler(m.Key, m.Value)
	}
}

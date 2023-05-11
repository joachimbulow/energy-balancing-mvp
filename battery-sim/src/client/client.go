package client

import (
	"fmt"
	"os"
)

var (
	brokerType = "KAFKA"
	brokerURL  = "127.0.0.1:29092" // Kafka PLAINTEXT_HOST://localhost:29092
)

type Client interface {
	Connect() error
	Disconnect() error
	Publish(topic string, key string, message string) error
	Subscribe(topic string) error
	Unsubscribe(topic string) error
	Listen(topic string, consumerGroupID string, handler func(params ...[]byte)) error
}

func NewClient() (Client, error) {
	if envBrokerURL := os.Getenv("BROKER_URL"); envBrokerURL != "" {
		brokerURL = envBrokerURL
	} else {
		logger.Info("BROKER_URL not set, using default: %s", brokerURL)
	}
	if envBrokerType := os.Getenv("BROKER"); envBrokerType != "" {
		brokerType = envBrokerType
	} else {
		logger.Info("BROKER_TYPE not set, using default: %s", brokerType)
	}

	switch brokerType {
	case "REDIS":
		return NewRedisClient()
	case "KAFKA":
		return NewKafkaClient()
	default:
		return nil, fmt.Errorf("invalid broker type: %s", brokerType)
	}
}

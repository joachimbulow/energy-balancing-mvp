package broker

import (
	"fmt"
	"os"
)

var (
	brokerType = "KAFKA"
	brokerURL  = "127.0.0.1:9092"
)

// define an interface for the broker so that we can use different brokers
type Broker interface {
	Connect() error
	Disconnect() error
	Publish(topic string, key string, message string) error
	Subscribe(topic string) error
	Unsubscribe(topic string) error
	Listen(topic string, handler func(params ...[]byte)) error
}

func NewBroker() (Broker, error) {
	if envBrokerURL := os.Getenv("BROKER_URL"); envBrokerURL != "" {
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
		return NewRedisBroker()
	case "KAFKA":
		return NewKafkaBroker()
	default:
		return nil, fmt.Errorf("invalid broker type: %s", brokerType)
	}
}

package broker

import (
	"fmt"
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

func NewBroker(brokerType string) (Broker, error) {
	switch brokerType {
	case "REDIS":
		return NewRedisBroker()
	case "KAFKA":
		return NewKafkaBroker()
	default:
		return nil, fmt.Errorf("invalid broker type: %s", brokerType)
	}
}

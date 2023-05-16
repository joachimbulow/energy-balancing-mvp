package client

import (
	"fmt"

	"github.com/joachimbulow/pem-energy-balance/src/util"
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
	brokerType := util.GetBroker()

	switch brokerType {
	case "REDIS":
		return NewRedisClient()
	case "KAFKA":
		return NewKafkaClient()
	default:
		return nil, fmt.Errorf("invalid broker type: %s", brokerType)
	}
}

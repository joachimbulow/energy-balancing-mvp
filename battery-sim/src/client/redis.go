package client

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

type RedisClient struct {
	client *redis.Client
	pubsub *redis.PubSub
	ctx    context.Context
}

const (
	address = "localhost:6379"
)

func NewRedisClient() (*RedisClient, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr: address,
	})

	if err := rdb.Ping(context.Background()).Err(); err != nil {
		return nil, fmt.Errorf("failed to ping Redis server: %w", err)
	}

	pubsub := rdb.Subscribe(context.Background(), "default")

	return &RedisClient{
		client: rdb,
		pubsub: pubsub,
		ctx:    context.Background(),
	}, nil
}

func (rb *RedisClient) Connect() error {
	return nil // Redis connection is established in NewRedisBroker()
}

func (rb *RedisClient) Disconnect() error {
	if err := rb.pubsub.Unsubscribe(rb.ctx, "default"); err != nil {
		return fmt.Errorf("failed to unsubscribe from topic: %w", err)
	}
	if err := rb.client.Close(); err != nil {
		return fmt.Errorf("failed to close Redis client connection: %w", err)
	}
	return nil
}

func (rb *RedisClient) Publish(topic string, key string, message string) error {
	if err := rb.client.Publish(rb.ctx, topic, message).Err(); err != nil {
		return fmt.Errorf("failed to publish message: %w", err)
	}
	return nil
}

func (rb *RedisClient) Subscribe(topic string) error {
	if err := rb.pubsub.Subscribe(rb.ctx, topic); err != nil {
		return fmt.Errorf("failed to subscribe to topic: %w", err)
	}
	return nil
}

func (rb *RedisClient) Unsubscribe(topic string) error {
	if err := rb.pubsub.Unsubscribe(rb.ctx, topic); err != nil {
		return fmt.Errorf("failed to unsubscribe from topic: %w", err)
	}
	return nil
}

func (rb *RedisClient) Listen(topic string, consumerGroupID string, handler func(params ...[]byte)) error {
	ch := rb.pubsub.Channel()

	go func() {
		for msg := range ch {
			if msg.Channel != topic {
				continue
			}

			var message map[string]interface{}
			if err := json.Unmarshal([]byte(msg.Payload), &message); err != nil {
				log.Printf("Failed to unmarshal message: %v", err)
				continue
			}

			if _, ok := message["timestamp"]; !ok {
				message["timestamp"] = time.Now().Unix()
			}

			bytes, err := json.Marshal(message)
			if err != nil {
				log.Printf("Failed to marshal message: %v", err)
				continue
			}

			handler(bytes)
		}
	}()

	return nil
}

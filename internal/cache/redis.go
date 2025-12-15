package cache

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Set(key string, value string) error {
	ctx := context.Background()

	err := s.client.Set(ctx, key, value, time.Duration(time.Duration.Minutes(20))).Err()
	if err != nil {
		return err
	}
	fmt.Println(228)
	return nil
}

func (s *Service) Get(key string) (string, error) {
	ctx := context.Background()

	val, err := s.client.Get(ctx, key).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

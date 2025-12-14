package cache

import (
	"github.com/redis/go-redis/v9"
)

type Service struct {
	client *redis.Client
}

func NewRedisService(client *redis.Client) *Service {
	return &Service{client: client}
}

func (s *Service) Set(key string, value string) {
	err := s.client.Set()
}

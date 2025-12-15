package cache

import (
	"log"
	"os"

	"github.com/redis/go-redis/v9"
)

func NewClient() *redis.Client {
	addr, exist := os.LookupEnv("REDIS_URL")
	if !exist {
		log.Println("redis addres not present")
	}

	return redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})
}

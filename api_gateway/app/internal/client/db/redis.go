package redis

import (
	"github.com/WTC-SYSTEM/wtc_system/api_gateway/internal/config"
	"github.com/go-redis/redis/v9"
)

func NewClient(config *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		DB:       0,
		Password: config.Redis.Password,
		Addr:     config.Redis.Addr,
	})
	return client
}

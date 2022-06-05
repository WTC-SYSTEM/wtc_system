package redis

import (
	"github.com/go-redis/redis/v9"
	"github.com/hawkkiller/wtc_system/api_gateway/internal/config"
)

func NewClient(config *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		DB:       0,
		Password: config.Redis.Password,
		Addr:     config.Redis.Addr,
	})
	return client
}

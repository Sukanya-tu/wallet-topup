package config

import (
	"github.com/redis/go-redis/v9"
)

func SetupRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr: GetRedisAddr(),
	})
}

package handlers

import (
	"os"

	"github.com/redis/go-redis/v9"
)

var redisClient *redis.Client

func GetRedisInstance() *redis.Client {

	addr := os.Getenv("REDIS_ADDR")
	psk := os.Getenv("REDIS_PSWD")

	if redisClient == nil {
		return redis.NewClient(&redis.Options{
			Addr:     addr,
			Password: psk,
		})
	}

	return redisClient
}

func DisconnectRedis() error {
	if redisClient == nil {
		return nil
	}

	err := redisClient.Close()
	redisClient = nil

	return err
}

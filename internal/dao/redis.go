package dao

import redis "github.com/go-redis/redis/v8"

func NewRedisClient() *redis.Client {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		},
		)
	return redisClient
}


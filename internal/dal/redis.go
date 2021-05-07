package dal

import (
	redis "github.com/go-redis/redis/v8"
	"sync"
)

var redisClient *redis.Client
var mutex sync.Mutex

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

// init once
func Redis() *redis.Client {

}
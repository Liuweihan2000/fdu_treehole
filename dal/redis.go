package dal

import (
	redis "github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

var redisClient *redis.Client

// 一个 client 对应着多个 connection
// 在同一个程序(生产环境中部署的实例)中，全局使用一个 client 即可
// 这个 client 是 go 程序对连接池的一个抽象，可能包含一个或多个 redis connection
// 这里说的 connection 是 go 程序中的 connection，其对应的就是 redis 中的一个或 client
// 使用 info clients 来查看 redis 当前的连接数
func NewRedisClient() *redis.Client {
	redisClient := redis.NewClient(
		&redis.Options{
			Addr:     viper.GetString("redis_source"),
			Password: viper.GetString("redis_password"),
			DB:       0,
		},
	)
	return redisClient
}

// init once for all connections in this program
func Redis() *redis.Client {
	return redisClient
}

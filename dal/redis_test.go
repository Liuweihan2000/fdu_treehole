package dal

import (
	"GoProject/fudan_bbs/utils"
	"context"
	"fmt"
	"github.com/spf13/viper"
	"testing"
)

type Index struct {
	ThreadID         int32  `json:"thread_id"`
	ThreadCreatedAt  string `json:"thread_created_at"`
	PostCount        int32  `json:"post_count"`
	FirstPostContent string `json:"first_post_content"`
	TimePassed       string `json:"time_passed"`
	ThreadUpdatedAt  string `json:"thread_updated_at"`
}

func TestRedisClient(t *testing.T) {
	viper.SetConfigType("yml")
	//viper.AddConfigPath("..\\")
	viper.SetConfigFile("..\\config.yml")
	err := viper.ReadInConfig() // Find and read the config file
	utils.FatalErrorHandle(err, "error while reading config file")
	fmt.Println(viper.GetString("redis_source"))
	fmt.Println(viper.GetString("redis_password"))
	rdb := NewRedisClient()
	cmd := rdb.Ping(context.Background())
	fmt.Println(cmd.Err())
	val, err := rdb.LRange(context.Background(), "non_exist_list", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val)
}

func TestRedisThread(t *testing.T) {
	rdb := NewRedisClient()

	res, err := rdb.LRange(context.Background(), "123123", 0, -1).Result()
	if err != nil {
		fmt.Printf("ERR: %v", err)
		fmt.Println()
	}
	fmt.Println("RES:")
	fmt.Println(res)
}

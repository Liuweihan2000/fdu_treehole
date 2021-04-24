package dao

import (
	"context"
	"encoding/json"
	"fmt"
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
	rdb := NewRedisClient()
	val, err := rdb.LRange(context.Background(), "non_exist_list", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val)
}

func TestRedisThread(t *testing.T) {
	rdb := NewRedisClient()

	index := &Index{
		ThreadID: 11,
		ThreadCreatedAt: "2021-3-12 08:56:38",
		PostCount: 22,
		FirstPostContent: "请问有人出英语视听教材吗",
		ThreadUpdatedAt: "2021-4-14 16:14:28",
	}
	bytes, _ := json.Marshal(index)
	str := string(bytes)

	cmd := rdb.LPush(context.Background(), "thread_list", str)
	if cmd != nil {
		fmt.Println(cmd)
	}

	redisSlice, e := rdb.LRange(context.Background(), "thread_list", 0, -1).Result()
	if e != nil {
		fmt.Println(e)
	}
	redisStr := redisSlice[0]
	fmt.Println(redisStr)
	indexStruct := &Index{}
	json.Unmarshal([]byte(redisStr), indexStruct)
	fmt.Println(indexStruct)
}

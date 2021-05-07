package dal

import (
	"context"
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

	res, err := rdb.LRange(context.Background(), "123123", 0, -1).Result()
	if err != nil {
		fmt.Printf("ERR: %v", err)
		fmt.Println()
	}
	fmt.Println("RES:")
	fmt.Println(res)
}

package dao

import (
	"context"
	"fmt"
	"testing"
)

func TestRedisClient(t *testing.T) {
	rdb := NewRedisClient()
	val, err := rdb.LRange(context.Background(), "alist", 0, -1).Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val)
}

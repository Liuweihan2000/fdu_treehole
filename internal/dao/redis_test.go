package dao

import (
	"context"
	"fmt"
	"testing"
)

func TestRedisClient(t *testing.T) {
	rdb := NewRedisClient()
	val, err := rdb.Get(context.Background(), "liuweihan").Result()
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(val)
}

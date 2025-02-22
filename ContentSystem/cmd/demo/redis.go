package main

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"time"
)

func main() {
	rdb := connRdb()
	ctx := context.Background()
	if err := rdb.Set(ctx, "session_id:admin", "session", 5*time.Second).Err(); err != nil {
		panic(err)
	}
	sessionID, err := rdb.Get(ctx, "session_id:admin").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(sessionID)
}

func connRdb() *redis.Client {
	rdb := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379", // use default Addr
		Password: "",               // no password set
		DB:       0,                // use default DB
	})
	_, err := rdb.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	return rdb
}

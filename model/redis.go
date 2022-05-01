package model

import (
	"context"
	"github.com/go-redis/redis/v8"
	"log"
)

func init() {
	log.Println("redis init")
	cli := redis.NewClient(&redis.Options{})

	r, err := cli.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
	log.Println("redis ping", r)
	println(cli.Set(context.Background(), "hello", "world", 0).Result())
}

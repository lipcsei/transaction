package main

import (
	"fmt"

	"github.com/foosio/api/lib/services/env"
	redis "gopkg.in/redis.v4"
)

type Redis struct{}

func (r *Redis) connect() *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     env.Get("REDIS_HOST", "localhost:6379"),
		Password: "",
		DB:       0,
	})
	return client
}

func main() {
	var redis Redis

	client := redis.connect()
	err := client.Set("key", "value", 0).Err()
	if err != nil {
		panic(err)
	}

	val, err := client.Get("key").Result()
	if err != nil {
		panic(err)
	}
	fmt.Println("key", val)

}

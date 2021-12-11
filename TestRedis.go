package main

import (
	"fmt"
	"github.com/go-redis/redis"
)

func main() {
	fmt.Println("golang连接redis")

	client := redis.NewClient(&redis.Options{
		Addr:     "124.70.95.61:6379",
		Password: "redis1921@",
		DB:       0,
	})

	pong, err := client.Ping().Result()
	fmt.Println(pong, err)

}

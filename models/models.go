package models

import (
	"github.com/go-redis/redis"
)

func ExampleNewClient() string {
	client := redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password:	"",
		DB:			0,
	})

	pong, err := client.Ping().Result()
	if err != nil {
		return err.Error()
	}
	
	return pong
}
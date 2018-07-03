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

func ExampleClient() string {
	client := redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password:	"",
		DB:			0,
	})

	// err := client.Set("name", "zhenhao", 0).Err()
	// if err != nil {
	// 	panic(err)
	// }

	val, err := client.Get("name").Result()
	if err != nil {
		panic(err)
	}
	return val

}

func TestGetHashTable() map[string]string {
	client := redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password:	"",
		DB:			0,
	})

	val, err := client.HGetAll("test").Result()
	if err != nil {
		panic(err)
	}

	return val
}
package models

import (
	"github.com/go-redis/redis"
	"encoding/json"
)

type Article struct {
	Title string `json:"title"`
	Content string `json:"content"`
	Time string `json:"time"`
	ReadCount string `json:"readCount"`
	Classify string `json:"classify"`
}

type Test struct {
	Title string `json:"title"`
	Content string `json:"content"`
	Id string `json:"id"`
	// ReadCount string `json:"readCount"`
	// Classify string `json:"classify"`
}

func HelloRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password:	"",
		DB:			0,
	})
}

func GetOneArticle(articleId string) string {
	client := HelloRedis()

	val, err := client.HGetAll(articleId).Result()
	if err != nil {
		panic(err)
	}

	return getArticleJSON(val)
}

func getArticleJSON(val map[string]string) string {
	var article Article

	article.Title		= val["title"]
	article.Content 	= val["content"]
	article.Time		= val["time"]
	article.ReadCount 	= val["readcount"]
	article.Classify	= val["classify"]

	var test Test

	test.Title		= val["title"]
	test.Content 	= val["content"]
	test.Id			= val["id"]
	// test.ReadCount 	= val["readcount"]
	// test.Classify	= val["classify"]

	result, err := json.Marshal(test)
	if err != nil {
		panic(err)
	}

	return string(result)
}

func GetAllArticles() map[string]string {
	client := HelloRedis()

	val, err := client.HGetAll("test").Result()
	if err != nil {
		panic(err)
	}

	return val
}

package models

import (
	"time"
	"fmt"
	// "fmt"
	"github.com/go-redis/redis"
	"encoding/json"
)

type Article struct {
	Id string `json:"id"`
	Info Info `json:"info"`  
}

type Info struct {
	Title string `json:"title"`
	Content string `json:"content"`
	Time string `json:"time"`
	ReadCount string `json:"readcount"`
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

func GetOneArticle(key string) string {
	client := HelloRedis()

	val, err := client.HGetAll(key).Result()
	if err != nil {
		panic(err)
	}

	return getArticleJSON(key, val)
}

func GetAllArticles() string {
	client := HelloRedis()

	val, err := client.LRange("testlist", 0, -1).Result()
	if err != nil {
		panic(err)
	}

	var articles string
	articles = "["
	for i := 0; i < len(val); i++ {
		articles += GetOneArticle(val[i])
		if i != len(val)-1 {
			articles += ","
		}
	}
	articles += "]"

	return articles
}

func AddArticle(inputs []byte) string {
	var articleInfo Info
	json.Unmarshal(inputs, &articleInfo)

	// client := HelloRedis()

	fmt.Println( time.Now().Format(time.ANSIC))

	// val, err := client.HMSet(article.Id, map[string]interface{} {
	// 	"title": articleInfo.Title,
	// 	"content": articleInfo.Content,
	// 	"classify": articleInfo.Classify,
	// 	"readcount": "0",
	// 	"time": time.Now().String(),
	// 	}).Result()
	// // fmt.Println(string(val))
	// if err != nil {
	// 	return "[{result: 'bu ok'}]"
	// 	panic(err)
	// }

	// _, err = client.LPush("testlist", article.Id).Result()
	// if err != nil {
	// 	return "[{result: 'bu ok'}]"
	// 	panic(err)
	// }

	return "[{result: 'ok'}]"
}

func getArticleJSON(key string, val map[string]string) string {
	var info Info
	info.Title		= val["title"]
	info.Content 	= val["content"]
	info.Time		= val["time"]
	info.ReadCount 	= val["readcount"]
	info.Classify	= val["classify"]

	var article Article
	article.Id 	= key
	article.Info = info

	result, err := json.Marshal(article)
	if err != nil {
		panic(err)
	}

	return string(result)
}


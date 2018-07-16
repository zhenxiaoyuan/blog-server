package models

import (
	"encoding/hex"
	// "strconv"
	"crypto/md5"
	"io"
	// "crypto/sha1"
	"fmt"
	"time"
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

var client *redis.Client

func init()  {
	client = HelloRedis()
	fmt.Println("Redis clent already initiated !")
}

func HelloRedis() *redis.Client {
	return redis.NewClient(&redis.Options{
		Addr:		"localhost:6379",
		Password:	"",
		DB:			0,
	})
}

func GetOneArticle(key string) string {
	// client := HelloRedis()

	val, err := client.HGetAll(key).Result()
	if err != nil {
		panic(err)
	}

	return getArticleJSON(key, val)
}

func GetAllArticles() string {
	// client := HelloRedis()

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

	articleId := getArticleId(articleInfo.Title)
	_, err := client.HMSet(articleId, map[string]interface{} {
		"title": articleInfo.Title,
		"content": articleInfo.Content,
		"classify": articleInfo.Classify,
		"readcount": "0",
		"time": time.Now().Format(time.ANSIC),
		}).Result()
	// fmt.Println(string(val))
	if err != nil {
		return "[{result: 'bu ok'}]"
		// panic(err)
	}

	_, err = client.LPush("testlist", articleId).Result()
	if err != nil {
		return "[{result: 'bu ok'}]"
		// panic(err)
	}

	return "[{result: 'ok'}]"
}

func UpdateArticle(inputs []byte) string {
	var articleInfo Info
	json.Unmarshal(inputs, &articleInfo)

	// 记得修正数据结构，在传输过程中通过id控制数据
	// 在前端应该加一层校验，如果没有做修改，就不需要回传了
	articleId := getArticleId(articleInfo.Title)
	_, err := client.HMSet(articleId, map[string]interface{} {
		"title": articleInfo.Title,
		"content": articleInfo.Content,
		"classify": articleInfo.Classify,
		}).Result()
	// fmt.Println(string(val))
	if err != nil {
		return "[{result: 'bu ok'}]"
		// panic(err)
	}

	// _, err = client.LPush("testlist", articleId).Result()
	// if err != nil {
	// 	return "[{result: 'bu ok'}]"
	// 	// panic(err)
	// }

	return "[{result: 'ok'}]"
}

func DeleteArticle(inputs []byte) string {
	var article Article
	json.Unmarshal(inputs, &article)

	// 第二个参数区分从头还是从尾进行查找，后期优化。
	_, err := client.LRem("testlist", 1, article.Id).Result()
	// fmt.Println(string(val))
	if err != nil {
		return "[{result: 'bu ok'}]"
		// panic(err)
	}

	_, err = client.Del(article.Id).Result()
	if err != nil {
		return "[{result: 'bu ok'}]"
		// panic(err)
	}

	return "{articleId: " + article.Id + "}"
}

// func isArticleExist(articleTitle string) bool {
// 	val, err := client.LRange("testlist", 0, -1).Result()
// 	if err != nil {
// 		return false
// 	}

// 	for _, article
// }

func getArticleId(articleTitle string) string {
	hasher := md5.New()
	b := []byte{}
	io.WriteString(hasher, articleTitle)
	
	return hex.EncodeToString(hasher.Sum(b))
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


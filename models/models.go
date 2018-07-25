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
	// Info Info `json:"info"`  
	Title string `json:"title"`
	Content string `json:"content"`
	Time string `json:"time"`
	ReadCount string `json:"readcount"`
	Classify string `json:"classify"`
}

// type Info struct {
// 	Title string `json:"title"`
// 	Content string `json:"content"`
// 	Time string `json:"time"`
// 	ReadCount string `json:"readcount"`
// 	Classify string `json:"classify"`
// }

// type Test struct {
// 	Title string `json:"title"`
// 	Content string `json:"content"`
// 	Id string `json:"id"`
// 	// ReadCount string `json:"readCount"`
// 	// Classify string `json:"classify"`
// }

var client *redis.Client
var listName string
var listLength int64

func init()  {
	client = HelloRedis()
	listName = "testlist"
	listLength, _ = client.LLen(listName).Result()
	fmt.Println(listLength)
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

	val, err := client.LRange(listName, 0, -1).Result()
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
	var article Article
	json.Unmarshal(inputs, &article)

	articleId := getArticleId(article.Title)
	currentTime := time.Now().Format(time.ANSIC)
	fmt.Println(articleId)
	// 增加查重操作
	// 解析的数据有误，明天查
	_, err := client.HMSet(articleId, map[string]interface{} {
		"title": article.Title,
		"content": article.Content,
		"classify": article.Classify,
		"readcount": "0",
		"time": currentTime, // 时间需要修改
		}).Result()
	// fmt.Println(string(val))
	if err != nil {
		return "[{id: 'bu ok'}]"
		// panic(err)
	}

	_, err = client.LPush(listName, articleId).Result()
	if err != nil {
		return "[{result: 'bu ok'}]"
		// panic(err)
	}

	article.Id 	= articleId
	// article.Info = info
	// article.Title		= val["title"]
	// article.Content 	= val["content"]
	article.Time		= currentTime
	article.ReadCount 	= "0"
	// article.Classify	= val["classify"]

	result, err := json.Marshal(article)
	if err != nil {
		panic(err)
	}


	// return "[" + string(result) + "]"
	return string(result)
}

func UpdateArticle(inputs []byte) string {
	var article Article
	json.Unmarshal(inputs, &article)

	// 做操作前先保存一个temp
	// 记得修正数据结构，在传输过程中通过id控制数据
	// 在前端应该加一层校验，如果没有做修改，就不需要回传了
	// 因为是根据文章标题进行判断，所以先要检测是否修改标题
	articleId := getArticleId(article.Title)
	if articleId != article.Id {
		// 标题做了修改
		// 修改链表value
		// 删除原来的hash
	} 

	// 获取旧数据

	// 更新新数据
	_, err := client.HMSet(articleId, map[string]interface{} {
		"title": article.Title,
		"content": article.Content,
		"classify": article.Classify,
		}).Result()
	// fmt.Println(string(val))
	if err != nil {
		return "[{result: 'bu ok'}]"
		// panic(err)
	}

	// 返回新id, readcount

	// _, err = client.LPush(listName, article.Id).Result()
	// if err != nil {
	// 	return "[{result: 'bu ok'}]"
	// 	// panic(err)
	// }

	return "[{result: 'update ok'}]"
}

func DeleteArticle(inputs []byte) string {
	var article Article
	json.Unmarshal(inputs, &article)

	// 第二个参数区分从头还是从尾进行查找，后期优化。
	_, err := client.LRem(listName, 1, article.Id).Result()
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
// 	val, err := client.LRange(listName, 0, -1).Result()
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
	// var info Info
	// info.Title		= val["title"]
	// info.Content 	= val["content"]
	// info.Time		= val["time"]
	// info.ReadCount 	= val["readcount"]
	// info.Classify	= val["classify"]

	var article Article
	article.Id 	= key
	// article.Info = info
	article.Title		= val["title"]
	article.Content 	= val["content"]
	article.Time		= val["time"]
	article.ReadCount 	= val["readcount"]
	article.Classify	= val["classify"]

	result, err := json.Marshal(article)
	if err != nil {
		panic(err)
	}

	return string(result)
}


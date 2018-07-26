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

type Result struct {
	Error bool `json:"error"`
	Info string `json:"info"`
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
	// 增加查重操作
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
	var result Result
	// 如果删除失败则使用inputs回滚

	json.Unmarshal(inputs, &article)

	// 查询准备删除的文章是否存在
	index, err := getArticleIndexById(article.Id)
	if err != nil {
		resultJSON := getResultJSON(&result, true, 
			"ERROR: Failed to get article index in list. ")
		
		return string(resultJSON)
	} 

	if index == -1 {
		// 文章不存在
		resultJSON := getResultJSON(&result, true, 
			"ERROR: Non-exist this article in list. ")
		
		return string(resultJSON)
	} 

	// 文章存在，根据索引位置判断Rem检索方向
	var searchDirection int64
	if index < listLength / 2.0 {
		searchDirection = 1
	} else {
		searchDirection = -1
	}
	_, err = client.LRem(listName, searchDirection, article.Id).Result()

	if err != nil {
		resultJSON := getResultJSON(&result, true, 
			"ERROR: Failed to delete article id in list. ")
		
		return string(resultJSON)
		// panic(err) 服务器是否需要panic，到底干什么用
	}

	_, err = client.Del(article.Id).Result()
	if err != nil {
		resultJSON := getResultJSON(&result, true, 
			"ERROR: Failed to delete article in hash. ")

		return string(resultJSON)
	}

	listLength, _ = client.LLen(listName).Result()

	resultJSON := getResultJSON(&result, false, article.Id)
	
	return string(resultJSON)
}

func getArticleIndexById(articleId string) (int64, error) {
	articles, err := client.LRange(listName, 0, -1).Result()
	if err != nil {
		// 获取数据库数据出错
		return int64(-1), err
	}

	for i := 0; i < len(articles); i++ {
		if articles[i] == articleId {
			// 找到索引并返回
			return int64(i), nil
		}
	}

	// 未找到索引并返回
	return int64(-1), nil
}

func getArticleId(articleTitle string) string {
	hasher := md5.New()
	b := []byte{}
	io.WriteString(hasher, articleTitle)
	
	return hex.EncodeToString(hasher.Sum(b))
}

func getArticleJSON(key string, val map[string]string) string {
	var article Article
	article.Id 	= key
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

func getResultJSON(result *Result, err bool, info string) []byte {
	result.Error = err
	result.Info = info
	resultJSON, _ := json.Marshal(result)

	return resultJSON
}


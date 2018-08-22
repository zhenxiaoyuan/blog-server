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
	// 此处是否可用枚举
	Type string `json:"type"`
	// Error bool `json:"error"`
	Info string `json:"info"`
}

// type Return struct {
// 	Type string `json:"type"`
// 	RArticle Article `json:"article"`
// 	RError
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

func GetAllArticles() string {
	var articles string
	var result Result

	articlesId, err := client.LRange(listName, 0, -1).Result()
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to get all articles id in list. ")

		return string(resultJSON)
	}

	articles = "["
	for i := 0; i < len(articlesId); i++ {
		article, err := client.HGetAll(articlesId[i]).Result()
		if err != nil {
			resultJSON := getResultJSON(&result, "error", 
				"ERROR: Failed to get an article in hash. ")

			return string(resultJSON)
		}
		articles += string(getArticleJSON(articlesId[i], article))

		if i != len(articlesId)-1 {
			articles += ","
		}
	}
	articles += "]"

	return articles
}

func AddArticle(inputs []byte) string {
	var article Article
	var result Result

	// 增加报错回滚机制

	json.Unmarshal(inputs, &article)

	articleId := getArticleId(article.Title)
	// 时间格式记得处理
	currentTime := time.Now().Format(time.ANSIC)

	// 查询准备添加的文章是否存在
	index, err := getArticleIndexById(articleId)
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to get article index in list. ")
		
		return string(resultJSON)
	} 

	if index != -1 {
		// 文章已存在
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Already-exist this article in list. ")
		
		return string(resultJSON)
	} 
	
	_, err = client.HMSet(articleId, map[string]interface{} {
		"title": article.Title,
		"content": article.Content,
		"classify": article.Classify,
		"readcount": "0",
		"time": currentTime, 
		}).Result()
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to add article in hash. ")
		
		return string(resultJSON)
	}

	_, err = client.LPush(listName, articleId).Result()
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to add article in list. ")
		
		return string(resultJSON)
	}

	article.Id 			= articleId
	article.Time		= currentTime
	article.ReadCount 	= "0"

	articleJSON, err := json.Marshal(article)
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to get article JSON. ")
		
		return string(resultJSON)
	}
	
	resultJSON := getResultJSON(&result, "article", string(articleJSON))

	return string(resultJSON)
}

func UpdateArticle(inputs []byte) string {
	var article Article
	var result Result
	json.Unmarshal(inputs, &article)

	// 做操作前先保存一个temp
	var tmpArticle map[string]string
	tmpArticle, err := client.HGetAll(article.Id).Result()
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to get old article data in hash. ")

		return string(resultJSON)
	}

	// 记得修正数据结构，在传输过程中通过id控制数据
	// 在前端应该加一层校验，如果没有做修改，就不需要回传了
	// 因为是根据文章标题进行判断，所以先要检测是否修改标题
	articleId := getArticleId(article.Title)
	// 标题做了修改
	if articleId != article.Id {
		// 修改链表value
		// 获取原文章所在index
		index, err := getArticleIndexById(article.Id)
		if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to get article index in list. ")
		
		return string(resultJSON)
		}

		_, err = client.LSet(listName, index, articleId).Result()
		if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to set article in list. ")

		return string(resultJSON)
		}

		// 获取旧数据
		// oldArticle, err := client.HGetAll(article.Id).Result()
		// if err != nil {
		// resultJSON := getResultJSON(&result, "error", 
		// 	"ERROR: Failed to get old article data in hash. ")

		// return string(resultJSON)
		// }

		// 备份旧数据
		article.Time = tmpArticle["time"]
		article.ReadCount = oldArticle.ReadCount
		_, err = client.HMSet(articleId, map[string]interface{} {
			"id": articleId,
			"time": oldArticle["time"],
			"readcount": oldArticle["readcount"],
		}).Result()
		if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to backup article in hash. ")

		return string(resultJSON)
		}

		// 删除原来的hash
		_, err = client.Del(article.Id).Result()
		if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to delete article in hash. ")

		return string(resultJSON)
		}

	} 

	// 无论标题改变与否，统一处理hash的变化
	// 更新新数据，处理hash
	_, err := client.HMSet(articleId, map[string]interface{} {
	"title": article.Title,
	"content": article.Content,
	"classify": article.Classify,
	}).Result()
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to backup article in hash. ")

		return string(resultJSON)
	}

	article.Id 			= articleId

	articleJSON, err := json.Marshal(article)
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to get article JSON. ")
		
		return string(resultJSON)
	}
	
	resultJSON := getResultJSON(&result, "article", string(articleJSON))

	return string(resultJSON)
	getArticleJSON(articleId)

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
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to get article index in list. ")
		
		return string(resultJSON)
	} 

	if index == -1 {
		// 文章不存在
		resultJSON := getResultJSON(&result, "error", 
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
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to delete article id in list. ")
		
		return string(resultJSON)
		// panic(err) 服务器是否需要panic，到底干什么用
	}

	_, err = client.Del(article.Id).Result()
	if err != nil {
		resultJSON := getResultJSON(&result, "error", 
			"ERROR: Failed to delete article in hash. ")

		return string(resultJSON)
	}

	listLength, _ = client.LLen(listName).Result()

	resultJSON := getResultJSON(&result, "article", article.Id)
	
	return string(resultJSON)
}

func getArticleIndexById(articleId string) (int64, error) {
	articlesId, err := client.LRange(listName, 0, -1).Result()
	if err != nil {
		// 获取数据库数据出错
		return int64(-1), err
	}

	for i := 0; i < len(articlesId); i++ {
		if articlesId[i] == articleId {
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

func getArticleJSON(key string, val map[string]string) []byte {
	var article Article

	article.Id 			= key
	article.Title		= val["title"]
	article.Content 	= val["content"]
	article.Time		= val["time"]
	article.ReadCount 	= val["readcount"]
	article.Classify	= val["classify"]

	articleJSON, err := json.Marshal(article)
	if err != nil {
		panic(err)
	}

	return articleJSON
}

func getResultJSON(result *Result, resultType string, resultInfo string) []byte {
	result.Type = resultType
	result.Info = resultInfo

	resultJSON, err := json.Marshal(result)
	if err != nil {
		// LOG
	}

	return resultJSON
}


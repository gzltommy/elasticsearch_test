package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"time"
)

// 索引 mapping 定义，这里仿微博消息结构定义
const mapping = `
{
  "mappings": {
    "properties": {
      "user": {
        "type": "keyword"
      },
      "message": {
        "type": "text"
      },
      "image": {
        "type": "keyword"
      },
      "created": {
        "type": "date"
      },
      "tags": {
        "type": "keyword"
      },
      "location": {
        "type": "geo_point"
      },
      "suggest_field": {
        "type": "completion"
      }
    }
  }
}`

type Weibo struct {
	User     string                `json:"user"`               // 用户
	Message  string                `json:"message"`            // 微博内容
	Retweets int                   `json:"retweets"`           // 转发数
	Image    string                `json:"image,omitempty"`    // 图片
	Created  time.Time             `json:"created,omitempty"`  // 创建时间
	Tags     []string              `json:"tags,omitempty"`     // 标签
	Location string                `json:"location,omitempty"` //位置
	Suggest  *elastic.SuggestField `json:"suggest_field,omitempty"`
}

func main() {
	// 创建client
	client, err := elastic.NewClient(
		//elastic.SetURL("http://192.168.24.132:9200", "http://127.0.0.1:9201"),
		elastic.SetURL("http://192.168.24.132:9200"),
		//elastic.SetBasicAuth("user", "secret")
		elastic.SetSniff(false),
	)
	if err != nil {
		// Handle error
		fmt.Printf("连接失败: %v\n", err)
		return
	} else {
		fmt.Println("连接成功")
	}

	// 执行 ES 请求需要提供一个上下文对象
	ctx := context.Background()

	// 1/首先检测下weibo索引是否存在
	exists, err := client.IndexExists("weibo").Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}

	// weibo 索引不存在，则创建一个
	if !exists {
		_, err := client.CreateIndex("weibo").BodyString(mapping).Do(ctx)
		if err != nil {
			// Handle error
			panic(err)
		}
	}

	// 创建创建一条微博
	//msg1 := Weibo{User: "olivere", Message: "打酱油的一天", Retweets: 0}
	//
	//// 使用client创建一个新的文档
	//put1, err := client.Index().
	//	Index("weibo"). // 设置索引名称
	//	Id("1"). // 设置文档id
	//	BodyJson(msg1). // 指定前面声明的微博内容
	//	Do(ctx) // 执行请求，需要传入一个上下文对象
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//
	//fmt.Printf("文档Id %s, 索引名 %s\n", put1.Id, put1.Index)

	//// 根据 id 查询文档
	//get1, err := client.Get().
	//	Index("weibo"). // 指定索引名
	//	Id("1"). // 设置文档id
	//	Do(ctx) // 执行请求
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//if get1.Found {
	//	fmt.Printf("文档id=%s 版本号=%d 索引名=%s\n", get1.Id, get1.Version, get1.Index)
	//}
	//
	//// 手动将文档内容转换成 go struct 对象
	//msg2 := Weibo{}
	//data, _ := get1.Source.MarshalJSON() // 提取文档内容，原始类型是 json 数据
	//json.Unmarshal(data, &msg2)          // 将 json 转成 struct 结果
	//
	//fmt.Println(msg2.Message)

	//_, err = client.Update().
	//	Index("weibo"). // 设置索引名
	//	Id("1"). // 文档 id
	//	Doc(map[string]interface{}{"retweets": 10}). // 更新 retweets=0，支持传入键值结构
	//	Do(ctx) // 执行 ES
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}

	// 根据 id 删除一条数据
	_, err = client.Delete().
		Index("weibo").
		Id("1").
		Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
}

package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"time"
)

// 定义一个文章索引结构，用来存储文章内容
type Article struct {
	Title   string    // 文章标题
	Content string    // 文章内容
	Author  string    // 作者
	Created time.Time // 发布时间
}

func main() {
	// 创建client
	client, err := elastic.NewClient(
		elastic.SetSniff(false), // docker 里面运行的 es 需要带上这个 option
		// elasticsearch 服务地址，多个服务地址使用逗号分隔
		elastic.SetURL("http://192.168.24.132:9200", "http://127.0.0.1:9201"),
		// 基于 http base auth 验证机制的账号和密码
		//elastic.SetBasicAuth("user", "secret"),
		// 启用 gzip 压缩
		elastic.SetGzip(true),
		// 设置监控检查时间间隔
		elastic.SetHealthcheckInterval(10*time.Second),
		// 设置请求失败最大重试次数
		elastic.SetMaxRetries(5),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置 info 日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))
	if err != nil {
		// Handle error
		fmt.Printf("连接失败: %v\n", err)
		return
	} else {
		fmt.Println("连接成功")
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	//// 定义一篇博客
	//blog := Article{Title: "golang es教程", Content: "go如何操作ES", Author: "tizi", Created: time.Now()}
	//
	//// 使用client创建一个新的文档
	//put1, err := client.Index().
	//	Index("blogs"). // 设置索引名称
	//	Id("1"). // 设置文档id
	//	BodyJson(blog). // 指定前面声明 struct 对象
	//	Do(ctx) // 执行请求，需要传入一个上下文对象
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//
	//fmt.Printf("文档Id %s, 索引名 %s\n", put1.Id, put1.Index)

	// 根据 id 查询文档
	//get1, err := client.Get().
	//	Index("blogs"). // 指定索引名
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
	//msg2 := Article{}
	//data, _ := get1.Source.MarshalJSON() // 提取文档内容，原始类型是 json 数据
	//json.Unmarshal(data, &msg2)          // 将 json 转成 struct 结果
	//fmt.Println(msg2.Title)              // 打印结果

	//// 查询 id 等于 1,2,3 的博客内容
	//result, err := client.MultiGet().
	//	Add(elastic.NewMultiGetItem(). // 通过 NewMultiGetItem 配置查询条件
	//		Index("blogs"). // 设置索引名
	//		Id("1")). // 设置文档 id
	//	Add(elastic.NewMultiGetItem().
	//		Index("blogs").
	//		Id("2")).
	//	Add(elastic.NewMultiGetItem().
	//		Index("blogs").
	//		Id("3")).
	//	Do(ctx) // 执行请求
	//
	//if err != nil {
	//	panic(err)
	//}
	//
	//// 遍历文档
	//for _, doc := range result.Docs {
	//	// 转换成struct对象
	//	var content Article
	//	tmp, _ := doc.Source.MarshalJSON()
	//	err := json.Unmarshal(tmp, &content)
	//	if err != nil {
	//		panic(err)
	//	}
	//
	//	fmt.Println(content.Title)
	//}

	// 根据 id 更新文档
	//_, err = client.Update().
	//	Index("blogs"). // 设置索引名
	//	Id("1"). // 文档id
	//	Doc(map[string]interface{}{"Title": "新的文章标题"}). // 更新Title="新的文章标题"，支持传入键值结构
	//	Do(ctx) // 执行ES查询
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}

	//// 根据条件更新文档
	//_, err = client.UpdateByQuery("blogs").
	//	Query(elastic.NewTermQuery("Author", "tizi")).// 设置查询条件，这里设置 Author=tizi
	//	Script(elastic.NewScript("ctx._source['Title']='1111111'")).// 通过脚本更新内容，将 Title 字段改为 1111111
	//	ProceedOnVersionConflict().// 如果文档版本冲突继续执行
	//	Do(ctx)
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}

	//// 根据 id 删除一条数据
	//_, err = client.Delete().
	//	Index("blogs").
	//	Id("1"). // 文档id
	//	Do(ctx)
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}

	// 根据条件 删除文档
	_, err = client.DeleteByQuery("blogs"). // 设置索引名
						Query(elastic.NewTermQuery("Author", "tizi")). // 设置查询条件为: Author = tizi
						ProceedOnVersionConflict().                    // 文档版本冲突也继续删除
						Do(ctx)
	if err != nil {
		// Handle error
		panic(err)
	}
}

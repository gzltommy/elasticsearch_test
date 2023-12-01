package main

import (
	"context"
	"fmt"
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"reflect"
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
		elastic.SetURL("http://192.168.151.97:9200"),
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
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)),
		// 设置 请求追踪（调试时启用）
		elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)), // 这一 必须的
	)

	if err != nil {
		// Handle error
		fmt.Printf("连接失败: %v\n", err)
		return
	} else {
		fmt.Println("连接成功")
	}

	// 执行ES请求需要提供一个上下文对象
	ctx := context.Background()

	//// 创建 term 查询条件，用于精确查询
	//termQuery := elastic.NewTermQuery("Author", "tizi")
	//
	//searchResult, err := client.Search().
	//	Index("blogs"). // 设置索引名
	//	Query(termQuery). // 设置查询条件
	//	Sort("Created", true). // 设置排序字段，根据 Created 字段升序排序，第二个参数 false 表示逆序
	//	From(0). // 设置分页参数 - 起始偏移量，从第 0 行记录开始
	//	Size(10). // 设置分页参数 - 每页大小
	//	Pretty(true). // 查询结果返回可读性较好的 JSON 格式
	//	Do(ctx) // 执行请求
	//
	//if err != nil {
	//	// Handle error
	//	panic(err)
	//}
	//
	//fmt.Printf("查询消耗时间 %d ms, 结果总数: %d\n", searchResult.TookInMillis, searchResult.TotalHits())
	//
	//if searchResult.TotalHits() > 0 {
	//	// 查询结果不为空，则遍历结果
	//	var b1 Article
	//	// 通过 Each 方法，将 es 结果的 json 结构转换成 struct 对象
	//	for _, item := range searchResult.Each(reflect.TypeOf(b1)) {
	//		// 转换成 Article 对象
	//		if t, ok := item.(Article); ok {
	//			fmt.Println(t.Title)
	//		}
	//	}
	//}

	//// 创建 terms 查询条件
	//termsQuery := elastic.NewTermsQuery("Author", "tizi", "tizi365")
	//
	//searchResult, err := client.Search().
	//	Index("blogs"). // 设置索引名
	//	Query(termsQuery). // 设置查询条件
	//	Sort("Created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
	//	From(0). // 设置分页参数 - 起始偏移量，从第0行记录开始
	//	Size(10). // 设置分页参数 - 每页大小
	//	Do(ctx) // 执行请求
	//
	//_ = searchResult

	//// 创建 match 查询条件
	//matchQuery := elastic.NewMatchQuery("Title", "golang es教程")
	//
	//searchResult, err := client.Search().
	//	Index("blogs"). // 设置索引名
	//	Query(matchQuery). // 设置查询条件
	//	Sort("Created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
	//	From(0). // 设置分页参数 - 起始偏移量，从第0行记录开始
	//	Size(10). // 设置分页参数 - 每页大小
	//	Do(ctx)
	//
	//_ = searchResult

	//// 例1 等价表达式： Created > "2020-07-20" and Created < "2020-07-29"
	//rangeQuery := elastic.NewRangeQuery("Created").
	//	Gt("2020-07-20").
	//	Lt("2020-07-29")
	//
	//// 例2 等价表达式： id >= 1 and id < 10
	//rangeQuery = elastic.NewRangeQuery("id").
	//	Gte(1).
	//	Lte(10)
	//
	//_ = rangeQuery

	// 创建 bool 查询

	// 创建 term 查询
	termQuery := elastic.NewTermQuery("Author", "tizi")
	matchQuery := elastic.NewMatchQuery("Title", "es教程")

	// 设置 bool 查询的 must 条件, 组合了两个子查询
	// 表示搜索匹配 Author=tizi 且 Title 匹配 "golang es 教程" 的文档
	boolQuery := elastic.NewBoolQuery().Must(termQuery, matchQuery)

	searchResult, err := client.Search().
		Index("blogs").        // 设置索引名
		Query(boolQuery).      // 设置查询条件
		Sort("Created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
		From(0).               // 设置分页参数 - 起始偏移量，从第0行记录开始
		Size(10).              // 设置分页参数 - 每页大小
		Do(ctx)                // 执行请求

	if err != nil {
		fmt.Println(err)
	}
	for _, item := range searchResult.Each(reflect.TypeOf(Article{})) {
		raw := item.(Article)
		fmt.Printf("%#v \n", raw)
	}

	//// 创建 bool 查询
	//boolQuery := elastic.NewBoolQuery().Must()
	//
	//// 创建 term 查询
	//termQuery := elastic.NewTermQuery("Author", "tizi")
	//
	//// 设置 bool 查询的 must not 条件
	//boolQuery.MustNot(termQuery)

	// 创建 bool 查询
	//boolQuery := elastic.NewBoolQuery().Must()
	//
	//// 创建 term 查询
	//termQuery := elastic.NewTermQuery("Author", "tizi")
	//matchQuery := elastic.NewMatchQuery("Title", "golang es教程")
	//
	//// 设置bool查询的should条件, 组合了两个子查询
	//// 表示搜索Author=tizi或者Title匹配"golang es教程"的文档
	//boolQuery.Should(termQuery, matchQuery)
	//searchResult, err := client.Search().
	//	Index("blogs").        // 设置索引名
	//	Query(boolQuery).      // 设置查询条件
	//	Sort("Created", true). // 设置排序字段，根据Created字段升序排序，第二个参数false表示逆序
	//	From(0).               // 设置分页参数 - 起始偏移量，从第0行记录开始
	//	Size(10).              // 设置分页参数 - 每页大小
	//	Do(ctx)                // 执行请求
	//
	//_ = searchResult
}

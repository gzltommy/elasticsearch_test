package main

import (
	"github.com/olivere/elastic/v7"
	"log"
	"os"
	"time"
)

func main() {
	client, err := elastic.NewClient(
		elastic.SetSniff(false), // docker 里面运行的 es 需要带上这个 option
		// elasticsearch 服务地址，多个服务地址使用逗号分隔
		elastic.SetURL("http://192.168.24.132:9200", "http://127.0.0.1:9201"),
		// 基于 http base auth 验证机制的账号和密码
		elastic.SetBasicAuth("user", "secret"),
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

		// 打印请求数据和返回数据（调试时启用）
		elastic.SetTraceLog(log.New(os.Stdout, "", log.LstdFlags)), // 这一 必须的
	)

	if err != nil {
		// Handle error
		panic(err)
	}
	_ = client

}

package main

import (
	"fmt"
	"github.com/small-ek/antgo/container/queue"
	"go.uber.org/zap"
	"time"
)

func main() {
	// 初始化配置
	cfg := queue.ClientConfig{
		Addr:        "127.0.0.1:6379",
		Password:    "",
		DB:          1,
		PoolSize:    100,
		DialTimeout: 10 * time.Second,
	}

	// 创建客户端（带日志配置）
	client := queue.NewClient(cfg,
		queue.WithLogger(zap.NewExample()),
	)

	// 提交高优先级任务
	_, err := client.Enqueue("order:process", map[string]interface{}{
		"order_id": "12345",
		"amount":   199.9,
	},
		queue.WithQueue("critical"),
		queue.WithMaxRetry(3),
		queue.WithUnique(30*time.Minute),
		queue.WithDeadline(time.Now().Add(1*time.Hour)),
	)

	if err != nil {
		// 错误处理逻辑
	}

	// 获取队列统计信息
	info, _ := client.GetQueueInfo("critical")
	fmt.Printf("Pending tasks: %d\n", info.Size)

	// 应用退出时清理
	defer client.Close()
}

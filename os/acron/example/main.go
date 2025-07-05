package main

import (
	"context"
	"github.com/small-ek/antgo/frame/ant"
	"github.com/small-ek/antgo/os/acron"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
	"time"
)

// 定义每秒执行的任务
type SecondJob struct {
	logger *zap.Logger
}

func (j *SecondJob) Run(ctx context.Context) error {

	alog.Debug(ctx, "每一秒任务执行", zap.Time("timestamp", time.Now()))

	return nil
}

// 定义每两秒执行的任务
type TwoSecondJob struct {
	logger *zap.Logger
}

func (j *TwoSecondJob) Run(ctx context.Context) error {

	time.Sleep(1 * time.Second)
	alog.Debug(ctx, "每两秒任务执行", zap.Time("timestamp", time.Now()))
	return nil
}

// 定义每两秒执行的任务
type ThreeSecondJob struct {
	logger *zap.Logger
}

func (j *ThreeSecondJob) Run(ctx context.Context) error {

	//time.Sleep(10 * time.Second)
	alog.Debug(ctx, "每三秒任务执行", zap.Time("timestamp", time.Now()))
	return nil
}

func main() {
	// 初始化日志
	ant.New("D:\\Work\\GoApp\\src\\loan-link\\config\\config_dev.toml")
	logger := ant.Log()
	// 创建根上下文
	ctx := context.Background()

	// 创建定时任务管理器 (默认超时设为3秒)
	cronManager := acron.New(ctx, logger, 3*time.Second)

	// 添加每秒执行的任务
	if err1 := cronManager.AddJob(
		"secondly_task",    // 任务ID
		"* * * * * *",      // Cron表达式: 每秒执行
		&SecondJob{logger}, // 任务实例
	); err1 != nil {
		logger.Fatal("添加每秒任务失败", zap.Error(err1))
	}

	// 添加每两秒执行的任务
	if err2 := cronManager.AddJob(
		"two_second_task",     // 任务ID
		"*/2 * * * * *",       // Cron表达式: 每两秒执行
		&TwoSecondJob{logger}, // 任务实例
	); err2 != nil {
		logger.Fatal("添加每两秒任务失败", zap.Error(err2))
	}

	// 添加每两秒执行的任务
	if err2 := cronManager.AddJobWithTimeout(
		"three_second_task",     // 任务ID
		"*/3 * * * * *",         // Cron表达式: 每两秒执行
		&ThreeSecondJob{logger}, // 任务实例
		5*time.Second,
	); err2 != nil {
		logger.Fatal("添加每三秒任务失败", zap.Error(err2))
	}

	// 启动定时器
	cronManager.Start()
	logger.Info("定时任务已启动")

	// 可选：添加健康检查
	cronManager.StartHealthCheck(10 * time.Second)

	// 模拟程序运行
	time.Sleep(30 * time.Second)

	// 优雅停止
	stopCtx := cronManager.Stop()
	select {
	case <-stopCtx.Done():
		logger.Info("所有任务已完成")
	case <-time.After(5 * time.Second):
		logger.Warn("强制停止，部分任务可能未完成")
	}
}

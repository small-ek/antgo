package acron

import (
	"github.com/robfig/cron/v3"
	"testing"
)

// 测试 Crontab 管理器的功能 / Test Crontab manager functionality
func TestCrontab(t *testing.T) {
	// 创建新的 Crontab 实例 / Create a new Crontab instance
	crontab := New()

	// 定义一个简单的任务函数 / Define a simple job function
	job := cron.NewChain().Then(cron.FuncJob(func() {}))

	// 测试添加任务 / Test adding a task
	t.Run("AddByID", func(t *testing.T) {
		err := crontab.AddByID("task1", "* * * * *", job)
		if err != nil {
			t.Fatalf("expected no error, got %v", err)
		}

		// 确认任务是否被添加 / Ensure the task was added
		if !crontab.IsExists("task1") {
			t.Fatal("task1 should exist")
		}
	})

	// 测试任务ID已存在 / Test adding a task with an existing ID
	t.Run("AddByID-Exist", func(t *testing.T) {
		err := crontab.AddByID("task1", "* * * * *", job)
		if err == nil {
			t.Fatal("expected error, got nil")
		}
	})

	// 测试删除任务 / Test deleting a task
	t.Run("DelByID", func(t *testing.T) {
		crontab.DelByID("task1")

		// 确认任务是否已删除 / Ensure the task was deleted
		if crontab.IsExists("task1") {
			t.Fatal("task1 should not exist")
		}
	})

	// 测试清理无效的任务ID / Test cleaning invalid task IDs
	t.Run("IDs", func(t *testing.T) {
		// 添加无效的任务ID并测试清理 / Add invalid task IDs and test cleanup
		invalidID := "task_invalid"
		crontab.AddByID(invalidID, "* * * * *", job)
		crontab.DelByID(invalidID)

		// 获取有效的任务ID列表 / Get the list of valid task IDs
		validIDs := crontab.IDs()
		if len(validIDs) > 0 {
			t.Fatalf("expected no valid tasks, got %v", validIDs)
		}
	})

	// 测试任务启动与停止 / Test starting and stopping the cron engine
	t.Run("StartStop", func(t *testing.T) {
		// 启动 crontab 引擎 / Start the crontab engine
		crontab.Start()

		// 确认 Cron 引擎已启动 / Ensure the Cron engine has started
		if !crontab.IsRunning() {
			t.Fatal("crontab engine should be running")
		}

		// 停止 crontab 引擎 / Stop the crontab engine
		crontab.Stop()

		// 确认 Cron 引擎已停止 / Ensure the Cron engine has stopped
		if crontab.IsRunning() {
			t.Fatal("crontab engine should be stopped")
		}
	})
}

// 测试根据函数添加 cron 任务 / Test adding a cron job by function
func TestAddByFunc(t *testing.T) {
	crontab := New()

	// 定义一个简单的任务函数 / Define a simple job function
	jobFunc := func() {
		// 这里可以做一些简单的日志打印，检查是否执行 / You can log something here to check if it runs
	}

	// 测试添加任务 / Test adding the task by function
	err := crontab.AddByFunc("task2", "* * * * *", jobFunc)
	if err != nil {
		t.Fatalf("expected no error, got %v", err)
	}

	// 确认任务是否被添加 / Ensure the task was added
	if !crontab.IsExists("task2") {
		t.Fatal("task2 should exist")
	}
}

// 测试添加无效 cron 任务 / Test adding an invalid cron job
func TestInvalidCronSpec(t *testing.T) {
	crontab := New()

	// 使用无效的 cron 表达式 / Use an invalid cron expression
	err := crontab.AddByID("task_invalid", "invalid cron spec", cron.NewChain().Then(cron.FuncJob(func() {})))
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

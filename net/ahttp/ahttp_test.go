// client_test.go
package ahttp

import (
	"net/http"
	"net/http/httptest"
	"runtime"
	"testing"
	"time"
)

// 测试基础请求
func TestBasicRequest(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := New(nil)
	resp, err := client.Request().Get(ts.URL)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("预期状态码200，实际得到 %d", resp.StatusCode())
	}
}

// 测试超时
func TestTimeout(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(200 * time.Millisecond)
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := New(&Config{Timeout: 100 * time.Millisecond})
	_, err := client.Request().Get(ts.URL)
	if err == nil {
		t.Fatal("预期超时错误，但没有收到错误")
	}

	// 简单检查错误信息
	if err.Error() != "Get \""+ts.URL+"\": context deadline exceeded" {
		t.Errorf("收到意外错误信息: %v", err)
	}
}

// 测试配置初始化
func TestConfigInit(t *testing.T) {
	defaultConfig := New(nil).config
	expectedIdle := runtime.GOMAXPROCS(0) * 200
	if defaultConfig.MaxIdleConnections != expectedIdle {
		t.Errorf("预期MaxIdleConnections=%d，实际得到%d", expectedIdle, defaultConfig.MaxIdleConnections)
	}
}

// 测试重试机制
func TestRetry(t *testing.T) {
	var attempt int
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		attempt++
		if attempt < 3 {
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
		w.WriteHeader(http.StatusOK)
	}))
	defer ts.Close()

	client := New(&Config{RetryAttempts: 3})
	resp, err := client.Request().Get(ts.URL)
	if err != nil {
		t.Fatalf("请求失败: %v", err)
	}

	if attempt != 3 {
		t.Errorf("预期3次尝试，实际尝试%d次", attempt)
	}

	if resp.StatusCode() != http.StatusOK {
		t.Errorf("最终状态码应为200，实际得到 %d", resp.StatusCode())
	}
}

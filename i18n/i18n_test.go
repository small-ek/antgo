// i18n/i18n_test.go
package i18n

import (
	"fmt"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
)

// setupTestEnvironment 创建测试环境
func setupTestEnvironment(t *testing.T) (string, func()) {
	tmpDir, err := os.MkdirTemp("", "i18n_test")
	if err != nil {
		t.Fatalf("Failed to create temp dir: %v", err)
	}

	// 创建测试翻译文件
	createTestFiles(t, tmpDir)

	// 返回清理函数
	return tmpDir, func() {
		//if err := os.RemoveAll(tmpDir); err != nil {
		//	t.Errorf("Failed to clean up temp dir: %v", err)
		//}
	}
}

// createTestFiles 创建测试翻译文件
func createTestFiles(t *testing.T, basePath string) {
	files := []struct {
		file    string
		content string
	}{
		{
			file: "en.toml",
			content: `
[common]
hello = "Hello, %s!"
welcome = "Welcome"

[plural]
apple = "apple"
apple.plural = "apples"

[date]
format = "2006-01-02"
`,
		},
		{
			file: "en.json",
			content: `{
"nested": {
	"key": "Nested Value"
}
}`,
		},
		{
			file: "zh-CN.yaml",
			content: `
common:
  hello: "你好，%s！"
  welcome: "欢迎"

plural:
  apple: "苹果"
  apple.plural: "多个苹果"

date:
  format: "2006年01月02日"
`,
		},
	}

	for _, f := range files {
		path := filepath.Join(basePath, f.file)
		if err := os.WriteFile(path, []byte(f.content), 0644); err != nil {
			t.Fatalf("Failed to create test file %s: %v", path, err)
		}
	}
}

// TestInit 测试初始化功能
func TestInit(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()
	fmt.Println(tmpDir)
	t.Run("正常初始化", func(t *testing.T) {
		config := Config{
			DefaultLanguage:    "en",
			FallbackLanguage:   "zh-CN",
			TranslationPath:    tmpDir,
			SupportedLanguages: []string{"en", "zh-CN"},
			EnableCache:        true,
			CacheSize:          500,
			DateFormats: map[string]string{
				"en":    time.RFC3339,
				"zh-CN": "2006年01月02日",
			},
		}

		if err := Init(config); err != nil {
			t.Fatalf("Init failed: %v", err)
		}

		if languageBundles["en"] == nil {
			t.Error("English bundle not loaded")
		}
		if languageBundles["zh-CN"] == nil {
			t.Error("Chinese bundle not loaded")
		}
	})

	t.Run("默认语言缺失", func(t *testing.T) {
		config := Config{
			DefaultLanguage:    "fr",
			TranslationPath:    tmpDir,
			SupportedLanguages: []string{"en"},
		}

		err := Init(config)
		if err == nil {
			t.Error("Expected error for missing default language, got nil")
		}
	})
}

// TestTranslation 测试翻译功能
func TestTranslation(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	config := Config{
		DefaultLanguage:    "en",
		FallbackLanguage:   "zh-CN",
		TranslationPath:    tmpDir,
		SupportedLanguages: []string{"en", "zh-CN"},
		EnableCache:        true,
		DateFormats: map[string]string{
			"en":    "2006-01-02",
			"zh-CN": "2006年01月02日",
		},
	}
	if err := Init(config); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	router := gin.Default()
	router.Use(Middleware())

	t.Run("基本翻译", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?lang=en", nil)
		router.ServeHTTP(w, req)

		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		got := T(ctx, "common.hello", "World")
		if got != "Hello, World!" {
			t.Errorf("Expected 'Hello, World!', got '%s'", got)
		}

		got = T(ctx, "common.welcome")
		if got != "Welcome" {
			t.Errorf("Expected 'Welcome', got '%s'", got)
		}
	})

	t.Run("嵌套键解析", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?lang=en", nil)
		router.ServeHTTP(w, req)
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		got := T(ctx, "nested.key")
		if got != "Nested Value" {
			t.Errorf("Expected 'Nested Value', got '%s'", got)
		}
	})

	t.Run("复数形式", func(t *testing.T) {
		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?lang=en", nil)
		router.ServeHTTP(w, req)
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		got := TPlural(ctx, 1, "plural.apple")
		if got != "1 apple" {
			t.Errorf("Expected '1 apple', got '%s'", got)
		}

		got = TPlural(ctx, 3, "plural.apple")
		if got != "3 apples" {
			t.Errorf("Expected '3 apples', got '%s'", got)
		}
	})

	t.Run("日期格式化", func(t *testing.T) {
		testDate := time.Date(2023, 5, 15, 0, 0, 0, 0, time.UTC)

		w := httptest.NewRecorder()
		req, _ := http.NewRequest("GET", "/?lang=en", nil)
		router.ServeHTTP(w, req)
		ctx, _ := gin.CreateTestContext(w)
		ctx.Request = req

		got := TDate(ctx, testDate)
		if got != "2023-05-15" {
			t.Errorf("Expected '2023-05-15', got '%s'", got)
		}
	})
}

// TestMiddleware 测试中间件功能
func TestMiddleware(t *testing.T) {
	tmpDir, cleanup := setupTestEnvironment(t)
	defer cleanup()

	config := Config{
		DefaultLanguage:    "en",
		TranslationPath:    tmpDir,
		SupportedLanguages: []string{"en", "zh-CN"},
	}
	if err := Init(config); err != nil {
		t.Fatalf("Init failed: %v", err)
	}

	router := gin.Default()
	router.Use(Middleware())

	tests := []struct {
		name           string
		setupRequest   func(*http.Request)
		expectedBundle string
	}{
		{
			name: "URL参数优先级",
			setupRequest: func(r *http.Request) {
				q := r.URL.Query()
				q.Add("lang", "zh-CN")
				r.URL.RawQuery = q.Encode()
			},
			expectedBundle: "zh-CN",
		},
		{
			name: "Cookie优先级",
			setupRequest: func(r *http.Request) {
				r.AddCookie(&http.Cookie{Name: "lang", Value: "zh-CN"})
			},
			expectedBundle: "zh-CN",
		},
		{
			name: "Accept-Language头",
			setupRequest: func(r *http.Request) {
				r.Header.Set("Accept-Language", "zh-CN;q=0.9,en;q=0.8")
			},
			expectedBundle: "zh-CN",
		},
		{
			name:           "默认语言回退",
			setupRequest:   func(r *http.Request) {},
			expectedBundle: "en",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			w := httptest.NewRecorder()
			req, _ := http.NewRequest("GET", "/", nil)
			tt.setupRequest(req)

			router.ServeHTTP(w, req)

			ctx, _ := gin.CreateTestContext(w)
			ctx.Request = req
			bundle := ctx.MustGet(contextKey).(*Bundle)
			if bundle.languageTag != tt.expectedBundle {
				t.Errorf("Expected bundle language '%s', got '%s'", tt.expectedBundle, bundle.languageTag)
			}
		})
	}
}

package config

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

// createTempFile 创建临时配置文件，并写入指定内容
func createTempFile(content string, ext string) (string, error) {
	// 使用系统临时目录，生成指定后缀的文件名
	f, err := os.CreateTemp("", "configtest_*."+ext)
	if err != nil {
		return "", err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return "", err
	}
	return f.Name(), nil
}

// TestNewAndGetFunctions 测试通过 New 方法加载配置，并使用 Get 系列函数获取值
func TestNewAndGetFunctions(t *testing.T) {
	// 重置全局配置变量，避免单例影响
	Config = nil

	// 创建一个临时 TOML 配置文件
	// 文件内容包含数字、字符串和布尔值
	content := `
port = 8080
name = "testApp"
enabled = true
`
	filename, err := createTempFile(content, "toml")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	// 测试结束后删除临时文件
	defer os.Remove(filename)

	// 初始化配置（单例）并注册读取配置
	cfg := New(filename)
	if err = cfg.Register(); err != nil {
		t.Fatalf("注册配置失败: %v", err)
	}

	// 测试各个 Get 函数
	if port := GetInt("port"); port != 8080 {
		t.Errorf("预期 port 为 %d, 实际得到 %d", 8080, port)
	}
	if name := GetString("name"); name != "testApp" {
		t.Errorf("预期 name 为 %s, 实际得到 %s", "testApp", name)
	}
	if enabled := GetBool("enabled"); !enabled {
		t.Errorf("预期 enabled 为 %t, 实际得到 %t", true, enabled)
	}
}

// TestAddConfigFile 测试通过 AddConfigFile 方法合并新的配置文件
func TestAddConfigFile(t *testing.T) {
	// 重置全局配置变量
	Config = nil

	// 创建基础配置文件
	baseContent := `
port = 8080
name = "baseApp"
`
	baseFile, err := createTempFile(baseContent, "toml")
	if err != nil {
		t.Fatalf("创建基础配置文件失败: %v", err)
	}
	defer os.Remove(baseFile)

	// 创建需要合并的配置文件，覆盖 name，并增加 debug 字段
	addContent := `
name = "mergedApp"
debug = true
`
	addFile, err := createTempFile(addContent, "toml")
	if err != nil {
		t.Fatalf("创建合并配置文件失败: %v", err)
	}
	defer os.Remove(addFile)

	// 初始化配置使用基础文件
	cfg := New(baseFile)
	if err = cfg.Register(); err != nil {
		t.Fatalf("注册基础配置失败: %v", err)
	}

	// 合并新的配置文件
	if err = AddConfigFile(addFile); err != nil {
		t.Fatalf("添加配置文件失败: %v", err)
	}

	// 给 watcher 一点时间更新（实际测试中尽量避免依赖 sleep，但此处仅作简单示例）
	time.Sleep(100 * time.Millisecond)

	// 验证配置合并结果：port 保持不变、name 被覆盖、debug 增加
	if port := GetInt("port"); port != 8080 {
		t.Errorf("预期 port 为 %d, 实际得到 %d", 8080, port)
	}
	if name := GetString("name"); name != "mergedApp" {
		t.Errorf("预期 name 为 %s, 实际得到 %s", "mergedApp", name)
	}
	if debug := GetBool("debug"); !debug {
		t.Errorf("预期 debug 为 %t, 实际得到 %t", true, debug)
	}
}

// TestSetKey 测试全局配置中直接设置键值对
func TestSetKey(t *testing.T) {
	// 重置全局配置变量
	Config = nil

	// 创建一个空配置文件（内容为空也能成功读取）
	content := ""
	filename, err := createTempFile(content, "toml")
	if err != nil {
		t.Fatalf("创建临时文件失败: %v", err)
	}
	defer os.Remove(filename)

	cfg := New(filename)
	if err = cfg.Register(); err != nil {
		t.Fatalf("注册空配置失败: %v", err)
	}

	// 使用 SetKey 设置一个自定义键值对，并验证
	SetKey("custom.key", "customValue")
	if v := GetString("custom.key"); v != "customValue" {
		t.Errorf("预期 custom.key 为 %s, 实际得到 %s", "customValue", v)
	}
}

// TestConfigFileExtension 测试不同配置文件扩展名的处理逻辑
func TestConfigFileExtension(t *testing.T) {
	// 重置全局配置变量
	Config = nil

	// 这里创建一个包含 JSON 格式的配置文件
	jsonContent := `{"port":9090,"name":"jsonApp"}`
	// 注意这里需要去掉扩展名前面的点
	filename, err := createTempFile(jsonContent, "json")
	if err != nil {
		t.Fatalf("创建 JSON 配置文件失败: %v", err)
	}
	defer os.Remove(filename)

	// 由于 New 方法根据传入的路径来设置配置类型，
	// 这里我们检查文件后缀是否正确处理（去掉开头的 .）
	cfg := New(filename)
	if err = cfg.Register(); err != nil {
		t.Fatalf("注册 JSON 配置失败: %v", err)
	}

	// 验证 JSON 中的配置项
	if port := GetInt("port"); port != 9090 {
		t.Errorf("预期 port 为 %d, 实际得到 %d", 9090, port)
	}
	if name := GetString("name"); name != "jsonApp" {
		t.Errorf("预期 name 为 %s, 实际得到 %s", "jsonApp", name)
	}
}

// TestSetupConfigFileForRemote 模拟远程配置的初始化逻辑（不实际连接远程服务）
// 此测试主要检查 setupConfigFile 方法对不同参数长度的处理逻辑
func TestSetupConfigFileForRemote(t *testing.T) {
	// 重置全局配置变量
	Config = nil

	// 模拟传入3个参数：代表远程配置（无安全认证）
	// 注意：此处不会实际访问远程服务，因此仅测试初始化过程中是否 panic
	defer func() {
		if r := recover(); r != nil {
			// 如果 panic，说明初始化失败（期望在测试中不 panic）
			t.Errorf("setupConfigFile panic: %v", r)
		}
	}()

	// 此处传入的参数仅用于测试逻辑，实际不会去连接
	New("etcd", "localhost:2379", "config.toml")
}

// TestFilePath处理 检查文件路径及扩展名处理是否正确
func TestFilePathHandling(t *testing.T) {
	// 测试 filepath.Ext 与 strings.TrimPrefix 的配合
	filename := "config.toml"
	ext := strings.TrimPrefix(filepath.Ext(filename), ".")
	if ext != "toml" {
		t.Errorf("预期扩展名为 %s, 实际得到 %s", "toml", ext)
	}
}

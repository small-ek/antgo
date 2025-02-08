package plugins

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

// mockPlugin 是一个简单的实现PluginFunc接口的假插件
type mockPlugin struct{}

// Before 模拟插件执行前的逻辑
func (m *mockPlugin) Before() interface{} {
	return "Before logic executed"
}

// After 模拟插件执行后的逻辑
func (m *mockPlugin) After(data ...interface{}) interface{} {
	return "After logic executed"
}

func TestPluginManager_Register(t *testing.T) {
	manager := New()

	// 测试插件注册
	err := manager.Register("plugin1", &mockPlugin{})
	assert.NoError(t, err, "Plugin registration should succeed")

	// 测试重复注册
	err = manager.Register("plugin1", &mockPlugin{})
	assert.Error(t, err, "Plugin registration should fail if the plugin is already registered")
	assert.Equal(t, "plugin already registered", err.Error(), "Error message should match")
}

func TestPluginManager_Uninstall(t *testing.T) {
	manager := New()

	// 注册插件
	err := manager.Register("plugin1", &mockPlugin{})
	assert.NoError(t, err, "Plugin registration should succeed")

	// 卸载插件
	err = manager.Uninstall("plugin1")
	assert.NoError(t, err, "Plugin uninstallation should succeed")

	// 尝试卸载不存在的插件
	err = manager.Uninstall("plugin1")
	assert.Error(t, err, "Plugin uninstallation should fail if the plugin doesn't exist")
	assert.Equal(t, "plugin not found", err.Error(), "Error message should match")
}

func TestPluginManager_List(t *testing.T) {
	manager := New()

	// 确保插件列表为空
	plugins := List()
	assert.Empty(t, plugins, "Plugin list should be empty initially")

	// 注册插件
	err := manager.Register("plugin1", &mockPlugin{})
	assert.NoError(t, err, "Plugin registration should succeed")

	// 确保插件列表不为空
	plugins = List()
	assert.NotEmpty(t, plugins, "Plugin list should not be empty after registration")
	assert.Len(t, plugins, 1, "Plugin list should contain exactly 1 plugin")
	assert.Contains(t, plugins, "plugin1", "Plugin list should contain the registered plugin")
}

func TestPluginManager_Singleton(t *testing.T) {
	// 确保返回的实例是单例
	manager1 := New()
	manager2 := New()
	assert.Same(t, manager1, manager2, "New() should return the same instance every time")
}

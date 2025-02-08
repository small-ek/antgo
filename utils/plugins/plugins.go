package plugins

import (
	"errors"
	"sync"
)

// PluginFunc defines the interface that all plugins must implement.
// 插件接口，所有插件必须实现该接口
type PluginFunc interface {
	// Before method will be executed before the main logic.
	// Before方法将在主要逻辑之前执行
	Before() interface{}

	// After method will be executed after the main logic.
	// After方法将在主要逻辑之后执行
	After(data ...interface{}) interface{}
}

// PluginManager is a struct that manages plugins and provides methods for registering and uninstalling them.
// PluginManager结构体管理插件，并提供注册和卸载插件的方法
type PluginManager struct {
	// PluginList is a map that holds all registered plugins by name.
	// PluginList是一个存储所有已注册插件的映射，以插件名称为键
	PluginList map[string]PluginFunc
	// mu is a mutex for managing concurrent access to PluginList.
	// mu是一个互斥锁，用于管理对PluginList的并发访问
	mu sync.RWMutex
}

var once sync.Once
var managerInstance *PluginManager

// New initializes and returns the singleton instance of PluginManager.
// New方法初始化并返回PluginManager的单例实例
func New() *PluginManager {
	once.Do(func() {
		managerInstance = &PluginManager{
			PluginList: make(map[string]PluginFunc),
		}
	})
	return managerInstance
}

// List returns the list of all registered plugins.
// List方法返回所有已注册插件的列表
func List() map[string]PluginFunc {
	managerInstance.mu.RLock()
	defer managerInstance.mu.RUnlock()
	return managerInstance.PluginList
}

// Register adds a new plugin to the manager. Returns an error if the plugin is already registered.
// Register方法将新插件添加到管理器中。如果插件已注册，则返回错误
func (p *PluginManager) Register(name string, plugin PluginFunc) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if the plugin is already registered.
	// 检查插件是否已注册
	if _, exists := p.PluginList[name]; exists {
		return errors.New("plugin already registered")
	}

	// Register the plugin.
	// 注册插件
	p.PluginList[name] = plugin
	return nil
}

// Uninstall removes a plugin by name. Returns an error if the plugin does not exist.
// Uninstall方法根据名称移除插件。如果插件不存在，则返回错误
func (p *PluginManager) Uninstall(name string) error {
	p.mu.Lock()
	defer p.mu.Unlock()

	// Check if the plugin exists.
	// 检查插件是否存在
	if _, exists := p.PluginList[name]; !exists {
		return errors.New("plugin not found")
	}

	// Remove the plugin.
	// 移除插件
	delete(p.PluginList, name)
	return nil
}

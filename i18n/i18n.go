package i18n

import (
	"encoding/json"
	"fmt"
	"github.com/BurntSushi/toml"
	"github.com/gin-gonic/gin"
	"github.com/small-ek/antgo/utils/conv"
	"golang.org/x/text/language"
	"gopkg.in/yaml.v3"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/****************************** 常量定义 Constants ******************************/

// 默认配置参数
const (
	DefaultLanguageCode   = "zh-CN"      // 默认语言代码 | Default language code
	ContextKeyLanguage    = "i18nBundle" // 上下文存储键名 | Context storage key name
	MaxCacheEntries       = 1000         // 最大缓存条目数 | Maximum cache entries
	DefaultDateTimeLayout = time.RFC3339 // 默认日期时间格式 | Default datetime format
)

/****************************** 全局变量 Global Variables ******************************/

var (
	// 语言包映射表 | Language bundles map
	languageBundleMap = make(map[string]*LanguageBundle)

	// 映射表读写锁 | Map read-write mutex
	bundleMapMutex sync.RWMutex

	// 语言匹配器 | Language matcher
	languagePriorityMatcher language.Matcher

	// 全局配置 | Global configuration
	globalConfiguration Config

	// 翻译缓存 | Translation cache
	cachedTranslations *translationCache
)

/****************************** 结构体定义 Struct Definitions ******************************/

// Config 国际化配置结构
// Internationalization configuration structure
type Config struct {
	DefaultLang      string            `comment:"默认语言代码（如'en'）| Default language code"`
	FallbackLang     string            `comment:"备用回退语言代码 | Fallback language code"`
	TranslationsDir  string            `comment:"翻译文件目录路径 | Translation files directory path"`
	SupportedLangs   []string          `comment:"支持的语言列表 | Supported languages list"`
	CacheEnabled     bool              `comment:"是否启用翻译缓存 | Enable translation cache"`
	MaxCacheSize     int               `comment:"缓存最大容量 | Maximum cache size"`
	CustomPluralRule PluralRuleHandler `comment:"自定义复数规则处理器 | Custom plural rule handler"`
	DateTimeLayouts  map[string]string `comment:"各语言日期时间格式 | Per-language datetime formats"`
}

// PluralRuleHandler 复数规则处理函数类型
// Plural rule handler function type
type PluralRuleHandler func(lang string, n int, key string, args ...interface{}) string

// LanguageBundle 语言包数据结构
// Language bundle data structure
type LanguageBundle struct {
	langTag          string                 // 语言标签 | Language tag
	flatTranslations map[string]interface{} // 扁平化翻译映射 | Flattened translation map
	pluralRuleFunc   func(n int) int        // 复数规则函数 | Plural rule function
	datetimeLayout   string                 // 日期时间格式 | Datetime format
	numberFormats    map[string]string      // 数字格式化配置 | Number formatting configs
}

// translationCache 线程安全的LRU翻译缓存
// Thread-safe LRU translation cache
type translationCache struct {
	sync.RWMutex
	entries     map[string]string // 缓存条目存储 | Cache entries storage
	accessOrder []string          // 访问顺序记录 | Access order records
	maxEntries  int               // 最大条目限制 | Maximum entries limit
}

/****************************** 初始化函数 Initialization Functions ******************************/

// Init 初始化国际化模块
// Initialize internationalization module
//
// 参数:
//
//	config - 配置参数结构体 | Configuration parameters struct
//
// 返回:
//
//	error - 初始化过程中遇到的错误 | Error during initialization
func New(config Config) error {
	// 设置配置默认值 | Set default configuration values
	setConfigDefaults(&config)

	globalConfiguration = config

	// 初始化翻译缓存 | Initialize translation cache
	initTranslationCache(config)

	// 初始化语言匹配器 | Initialize language matcher
	if err := initLanguageMatcher(config); err != nil {
		return err
	}

	// 加载所有支持语言的翻译包 | Load bundles for all supported languages
	if err := loadAllLanguageBundles(config); err != nil {
		return err
	}

	// 验证默认和备用语言包 | Validate default and fallback bundles
	if err := validateCoreBundles(config); err != nil {
		return err
	}

	return nil
}

// setConfigDefaults 设置配置默认值
// Set configuration default values
func setConfigDefaults(config *Config) {
	if config.DefaultLang == "" {
		config.DefaultLang = DefaultLanguageCode
	}
	if config.MaxCacheSize == 0 {
		config.MaxCacheSize = MaxCacheEntries
	}
}

// initTranslationCache 初始化翻译缓存
// Initialize translation cache
func initTranslationCache(config Config) {
	cachedTranslations = &translationCache{
		entries:     make(map[string]string),
		accessOrder: make([]string, 0, config.MaxCacheSize),
		maxEntries:  config.MaxCacheSize,
	}
}

// initLanguageMatcher 初始化语言匹配器
// Initialize language matcher
func initLanguageMatcher(config Config) error {
	tags := make([]language.Tag, 0, len(config.SupportedLangs))
	for _, code := range config.SupportedLangs {
		tags = append(tags, language.Make(code))
	}
	languagePriorityMatcher = language.NewMatcher(tags)
	return nil
}

// loadAllLanguageBundles 加载所有语言包
// Load all language bundles
func loadAllLanguageBundles(config Config) error {
	for _, langCode := range config.SupportedLangs {
		bundle, err := loadLanguageBundle(langCode)
		if err != nil {
			return fmt.Errorf("加载语言包失败 [%s]: %w", langCode, err)
		}
		languageBundleMap[langCode] = bundle
	}
	return nil
}

// validateCoreBundles 验证核心语言包
// Validate core language bundles
func validateCoreBundles(config Config) error {
	if _, exists := languageBundleMap[config.DefaultLang]; !exists {
		return fmt.Errorf("默认语言包缺失 [%s]", config.DefaultLang)
	}
	if config.FallbackLang != "" {
		if _, exists := languageBundleMap[config.FallbackLang]; !exists {
			return fmt.Errorf("备用语言包缺失 [%s]", config.FallbackLang)
		}
	}
	return nil
}

/****************************** 语言包加载 Language Bundle Loading ******************************/

// loadLanguageBundle 加载指定语言包
// Load specified language bundle
func loadLanguageBundle(langCode string) (*LanguageBundle, error) {
	translations := make(map[string]interface{})

	// 遍历翻译目录 | Walk through translation directory
	err := filepath.WalkDir(globalConfiguration.TranslationsDir, func(path string, entry fs.DirEntry, err error) error {
		if err != nil || entry.IsDir() {
			return err
		}

		// 解析文件信息 | Parse file information
		ext := strings.ToLower(filepath.Ext(path))
		baseName := strings.TrimSuffix(filepath.Base(path), ext)

		// 匹配目标语言 | Match target language
		if !strings.HasPrefix(baseName, langCode) {
			return nil
		}

		// 解析翻译文件 | Parse translation file
		if err := parseTranslationFile(path, ext, translations); err != nil {
			return fmt.Errorf("文件解析失败 [%s]: %w", path, err)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("目录遍历错误: %w", err)
	}

	return createLanguageBundle(langCode, translations), nil
}

// parseTranslationFile 解析翻译文件
// Parse translation file
func parseTranslationFile(path, ext string, translations map[string]interface{}) error {
	file, err := os.Open(path)
	if err != nil {
		return err
	}
	defer file.Close()

	switch strings.TrimPrefix(ext, ".") {
	case "toml":
		_, err = toml.NewDecoder(file).Decode(&translations)
	case "yaml", "yml":
		err = yaml.NewDecoder(file).Decode(&translations)
	case "json":
		err = json.NewDecoder(file).Decode(&translations)
	default:
		return fmt.Errorf("不支持的格式类型: %s", ext)
	}

	if err != nil {
		return fmt.Errorf("解码错误: %w", err)
	}
	return nil
}

// createLanguageBundle 创建语言包实例
// Create language bundle instance
func createLanguageBundle(langCode string, data map[string]interface{}) *LanguageBundle {
	return &LanguageBundle{
		langTag:          langCode,
		flatTranslations: flattenTranslations(data),
		datetimeLayout:   resolveDatetimeLayout(langCode),
	}
}

// resolveDatetimeLayout 解析日期时间格式
// Resolve datetime format
func resolveDatetimeLayout(langCode string) string {
	if layout, exists := globalConfiguration.DateTimeLayouts[langCode]; exists {
		return layout
	}
	return DefaultDateTimeLayout
}

// flattenTranslations 扁平化嵌套翻译结构
// Flatten nested translation structure
func flattenTranslations(nested map[string]interface{}) map[string]interface{} {
	flat := make(map[string]interface{})
	for key, value := range nested {
		switch v := value.(type) {
		case map[string]interface{}:
			sub := flattenTranslations(v)
			for sk, sv := range sub {
				flat[key+"."+sk] = sv
			}
		default:
			flat[key] = v
		}
	}
	return flat
}

/****************************** Gin中间件 Gin Middleware ******************************/

// Middleware 语言包注入中间件
// Language bundle injection middleware
func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 检测客户端语言 | Detect client language
		lang := detectClientLanguage(c)

		// 获取对应语言包 | Get corresponding bundle
		bundle := getBundleByLanguage(lang)

		// 注入上下文 | Inject into context
		c.Set(ContextKeyLanguage, bundle)
		c.Next()
	}
}

/****************************** 语言检测 Language Detection ******************************/

// detectClientLanguage 检测客户端首选语言
// Detect client preferred language
func detectClientLanguage(c *gin.Context) string {
	// 1. 检查URL参数 | Check URL parameter

	if lang := c.Query("lang"); isValidLanguage(lang) {
		return lang
	}

	// 2. 检查Cookie设置 | Check cookie setting
	if lang, err := c.Cookie("lang"); err == nil && isValidLanguage(lang) {
		return lang
	}

	// 3. 解析Accept-Language头 | Parse Accept-Language header
	acceptLang := c.GetHeader("Accept-Language")

	matchedTag, _ := language.MatchStrings(languagePriorityMatcher, acceptLang)
	baseTag, _ := matchedTag.Base()

	if matchedTag.String() != "" && isValidLanguage(matchedTag.String()) {
		return matchedTag.String()
	}
	if baseTag.String() != "" && isValidLanguage(baseTag.String()) {
		return baseTag.String()
	}

	// 4. 返回默认语言 | Return default language
	return globalConfiguration.DefaultLang
}

// isValidLanguage 验证语言有效性
// Validate language validity
func isValidLanguage(code string) bool {
	for _, lang := range globalConfiguration.SupportedLangs {
		if lang == code {
			return true
		}
	}
	return false
}

/****************************** 语言包获取 Bundle Retrieval ******************************/

// getBundleByLanguage 获取对应语言包
// Get corresponding language bundle
func getBundleByLanguage(langCode string) *LanguageBundle {
	bundleMapMutex.RLock()
	defer bundleMapMutex.RUnlock()

	if bundle, exists := languageBundleMap[langCode]; exists {
		return bundle
	}

	// 回退到备用语言 | Fallback to secondary language
	if globalConfiguration.FallbackLang != "" {
		if bundle, exists := languageBundleMap[globalConfiguration.FallbackLang]; exists {
			return bundle
		}
	}

	// 回退到默认语言 | Fallback to default
	return languageBundleMap[globalConfiguration.DefaultLang]
}

/****************************** 翻译核心功能 Translation Core ******************************/

// T 获取翻译文本
// Get translated text

func T(c *gin.Context, key string, args ...interface{}) string {
	// 尝试缓存获取 | Try cache retrieval
	fmt.Println(globalConfiguration.CacheEnabled)
	if globalConfiguration.CacheEnabled {
		if cached := getCachedTranslation(c, key, args...); cached != "" {
			return cached
		}
	}

	// 执行翻译 | Perform translation
	bundle := c.MustGet(ContextKeyLanguage).(*LanguageBundle)
	result := bundle.translate(key, args...)

	// 缓存结果 | Cache result
	if globalConfiguration.CacheEnabled {
		cacheTranslation(c, key, result, args...)
	}
	return result
}

// TPlural 获取复数形式翻译
// Get plural-form translation
func TPlural(c *gin.Context, count int, key string, args ...interface{}) string {
	bundle := c.MustGet(ContextKeyLanguage).(*LanguageBundle)
	return bundle.pluralize(count, key, args...)
}

// TDate 本地化日期时间格式化
// Localized datetime formatting
func TDate(c *gin.Context, t time.Time) string {
	bundle := c.MustGet(ContextKeyLanguage).(*LanguageBundle)
	return t.Format(bundle.datetimeLayout)
}

/****************************** 语言包方法 Bundle Methods ******************************/

// translate 执行翻译逻辑
// Perform translation logic
func (b *LanguageBundle) translate(key string, args ...interface{}) string {
	// 当前语言包查找 | Lookup in current bundle
	if val, exists := b.flatTranslations[key]; exists {
		return formatString(conv.String(val), args)
	}

	// 回退语言处理 | Fallback handling
	if fbBundle := getBundleByLanguage(globalConfiguration.FallbackLang); fbBundle != nil {
		return fbBundle.translate(key, args...)
	}

	// 默认语言处理 | Default language handling
	if defBundle := getBundleByLanguage(globalConfiguration.DefaultLang); defBundle != nil {
		return defBundle.translate(key, args...)
	}

	return key // 最终退回键名 | Fallback to key
}

func (b *LanguageBundle) pluralize(n int, key string, args ...interface{}) string {
	// 优先使用自定义规则 | Priority to custom rules
	if globalConfiguration.CustomPluralRule != nil {
		return globalConfiguration.CustomPluralRule(b.langTag, n, key, args...)
	}

	// 默认复数规则 | Default plural rule
	pluralKey := key
	if n != 1 {
		pluralKey = key + ".plural"
	}

	// 提前判断是否存在对应的翻译
	if val, exists := b.flatTranslations[pluralKey]; exists {
		// 判断是否需要格式化
		if n == 1 {
			return val.(string) // 返回单数形式
		}

		// 使用 fmt.Sprintf 来格式化字符串，避免不必要的 conv.String 转换
		return fmt.Sprintf(val.(string), append([]interface{}{n}, args...)...)
	}

	// 如果找不到翻译，调用默认翻译
	return b.translate(key, args...)
}

// formatString 格式化字符串
// Format string with arguments
func formatString(base string, args []interface{}) string {
	if len(args) > 0 {
		return fmt.Sprintf(base, args...)
	}
	return base
}

/****************************** 缓存管理 Cache Management ******************************/

// getCachedTranslation 获取缓存翻译
// Get cached translation
func getCachedTranslation(c *gin.Context, key string, args ...interface{}) string {
	// 获取当前语言 | Get the current language from context
	lang := c.MustGet(ContextKeyLanguage).(*LanguageBundle).langTag
	cacheKey := generateCacheKey(lang, key, args...)
	cachedTranslations.RLock()
	defer cachedTranslations.RUnlock()
	return cachedTranslations.entries[cacheKey]
}

// cacheTranslation 缓存翻译结果
// Cache translation result
func cacheTranslation(c *gin.Context, key, value string, args ...interface{}) {
	// 获取当前语言 | Get the current language from context
	lang := c.MustGet(ContextKeyLanguage).(*LanguageBundle).langTag
	cacheKey := generateCacheKey(lang, key, args...)
	cachedTranslations.Lock()
	defer cachedTranslations.Unlock()

	// LRU淘汰机制 | LRU eviction mechanism
	if len(cachedTranslations.accessOrder) >= cachedTranslations.maxEntries {
		oldest := cachedTranslations.accessOrder[0]
		delete(cachedTranslations.entries, oldest)
		cachedTranslations.accessOrder = cachedTranslations.accessOrder[1:]
	}

	cachedTranslations.entries[cacheKey] = value
	cachedTranslations.accessOrder = append(cachedTranslations.accessOrder, cacheKey)
}

// generateCacheKey 生成缓存键
// Generate cache key
func generateCacheKey(lang, key string, args ...interface{}) string {
	// 在缓存键中包括语言 | Include language in the cache key
	return fmt.Sprintf("%s|%s|%v", lang, key, args)
}

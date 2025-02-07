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

// 常量定义
const (
	defaultLanguage   = "en"         // 默认语言
	contextKey        = "i18nBundle" // Gin上下文存储键
	maxCacheSize      = 1000         // 最大缓存条目数
	defaultDateFormat = time.RFC3339 // 默认日期格式
)

// 全局变量
var (
	languageBundles  = make(map[string]*Bundle) // 语言包集合
	bundlesLock      sync.RWMutex               // 语言包读写锁
	languageMatcher  language.Matcher           // 语言匹配器
	globalConfig     Config                     // 全局配置
	translationCache *lruCache                  // 翻译缓存
)

// Config 国际化配置结构
type Config struct {
	DefaultLanguage    string            // 默认语言代码，例如"en"
	FallbackLanguage   string            // 备用回退语言代码
	TranslationPath    string            // 翻译文件路径
	SupportedLanguages []string          // 支持的语言列表
	EnableCache        bool              // 是否启用翻译缓存
	CacheSize          int               // 缓存最大容量
	PluralRules        PluralRuleFunc    // 自定义复数规则处理函数
	DateFormats        map[string]string // 各语言日期格式配置
}

// PluralRuleFunc 复数规则处理函数类型定义
type PluralRuleFunc func(lang string, n int, key string, args ...interface{}) string

// Bundle 语言包结构
type Bundle struct {
	languageTag    string                 // 语言标签
	translations   map[string]interface{} // 扁平化翻译条目
	pluralRule     func(n int) int        // 复数形式计算函数
	dateTimeFormat string                 // 日期时间格式
	numberFormats  map[string]string      // 数字格式配置
}

// lruCache LRU缓存结构
type lruCache struct {
	sync.RWMutex
	cacheEntries map[string]string // 缓存条目存储
	cacheKeys    []string          // 缓存键顺序记录
	maxCapacity  int               // 最大缓存容量
}

// Init 初始化国际化配置
// Initialize internationalization configuration
func Init(config Config) error {
	// 配置默认值处理
	if config.DefaultLanguage == "" {
		config.DefaultLanguage = defaultLanguage
	}
	if config.CacheSize == 0 {
		config.CacheSize = maxCacheSize
	}

	globalConfig = config

	// 初始化缓存
	translationCache = &lruCache{
		cacheEntries: make(map[string]string),
		cacheKeys:    make([]string, 0, config.CacheSize),
		maxCapacity:  config.CacheSize,
	}

	// 初始化语言匹配器
	languageTags := make([]language.Tag, 0, len(config.SupportedLanguages))
	for _, langCode := range config.SupportedLanguages {
		languageTags = append(languageTags, language.Make(langCode))
	}
	languageMatcher = language.NewMatcher(languageTags)

	// 加载所有支持语言的翻译包
	for _, langCode := range config.SupportedLanguages {
		bundle, err := loadLanguageResources(langCode)
		if err != nil {
			if langCode == config.FallbackLanguage {
				return fmt.Errorf("fallback language %s load failed: %w", langCode, err)
			}
			continue
		}
		languageBundles[langCode] = bundle
	}
	
	// 校验默认语言包
	if _, exists := languageBundles[config.DefaultLanguage]; !exists {
		return fmt.Errorf("default language %s bundle not found", config.DefaultLanguage)
	}

	// 校验备用语言包
	if config.FallbackLanguage != "" {
		if _, exists := languageBundles[config.FallbackLanguage]; !exists {
			return fmt.Errorf("fallback language %s bundle not found", config.FallbackLanguage)
		}
	}

	return nil
}

// loadLanguageResources 加载指定语言的翻译资源
// Load translation resources for specified language
func loadLanguageResources(langCode string) (*Bundle, error) {
	translations := make(map[string]interface{})

	// 遍历翻译文件目录
	err := filepath.WalkDir(globalConfig.TranslationPath, func(filePath string, dirEntry fs.DirEntry, err error) error {
		if err != nil || dirEntry.IsDir() {
			return err
		}

		// 解析文件名和扩展名
		fileExt := filepath.Ext(filePath)
		fileName := strings.TrimSuffix(filepath.Base(filePath), fileExt)

		// 匹配语言代码
		if !strings.HasPrefix(fileName, langCode) {
			return nil
		}

		// 打开翻译文件
		file, openErr := os.Open(filePath)
		if openErr != nil {
			return openErr
		}
		defer file.Close()

		// 根据文件类型解析
		switch strings.ToLower(fileExt[1:]) {
		case "toml":
			if _, decodeErr := toml.NewDecoder(file).Decode(&translations); decodeErr != nil {
				return fmt.Errorf("TOML decode error: %w", decodeErr)
			}
		case "yaml", "yml":
			if decodeErr := yaml.NewDecoder(file).Decode(&translations); decodeErr != nil {
				return fmt.Errorf("YAML decode error: %w", decodeErr)
			}
		case "json":
			if decodeErr := json.NewDecoder(file).Decode(&translations); decodeErr != nil {
				return fmt.Errorf("JSON decode error: %w", decodeErr)
			}
		default:
			return fmt.Errorf("unsupported file format: %s", fileExt)
		}
		return nil
	})

	if err != nil {
		return nil, fmt.Errorf("directory traversal error: %w", err)
	}

	return &Bundle{
		languageTag:    langCode,
		translations:   flattenNestedTranslations(translations),
		dateTimeFormat: resolveDateTimeFormat(langCode),
	}, nil
}

// resolveDateTimeFormat 解析日期时间格式
func resolveDateTimeFormat(langCode string) string {
	if format, exists := globalConfig.DateFormats[langCode]; exists {
		return format
	}
	return defaultDateFormat
}

// flattenNestedTranslations 扁平化嵌套的翻译结构
func flattenNestedTranslations(nestedData map[string]interface{}) map[string]interface{} {
	flatData := make(map[string]interface{})
	for key, value := range nestedData {
		switch v := value.(type) {
		case map[string]interface{}:
			// 递归处理嵌套结构
			subItems := flattenNestedTranslations(v)
			for subKey, subValue := range subItems {
				flatData[key+"."+subKey] = subValue
			}
		default:
			flatData[key] = v
		}
	}
	return flatData
}

// Middleware Gin中间件，用于注入语言包到上下文
// Gin middleware for injecting language bundle into context
func Middleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		detectedLang := detectRequestLanguage(ctx)
		bundle := getLanguageBundle(detectedLang)
		ctx.Set(contextKey, bundle)
		ctx.Next()
	}
}

// detectRequestLanguage 检测请求的语言偏好
func detectRequestLanguage(ctx *gin.Context) string {
	// 1. 检查URL参数
	if langParam := ctx.Query("lang"); langParam != "" && isValidLanguageCode(langParam) {
		return langParam
	}

	// 2. 检查Cookie设置
	if langCookie, err := ctx.Cookie("lang"); err == nil && isValidLanguageCode(langCookie) {
		return langCookie
	}

	// 3. 解析Accept-Language头
	acceptLangHeader := ctx.GetHeader("Accept-Language")
	matchedTag, _ := language.MatchStrings(languageMatcher, acceptLangHeader)
	baseTag, _ := matchedTag.Base()
	if baseTag.String() != "" && isValidLanguageCode(baseTag.String()) {
		return baseTag.String()
	}

	// 4. 返回默认语言
	return globalConfig.DefaultLanguage
}

// isValidLanguageCode 验证语言代码有效性
func isValidLanguageCode(langCode string) bool {
	for _, supportedLang := range globalConfig.SupportedLanguages {
		if supportedLang == langCode {
			return true
		}
	}
	return false
}

// getLanguageBundle 获取对应的语言包
func getLanguageBundle(langCode string) *Bundle {
	bundlesLock.RLock()
	defer bundlesLock.RUnlock()

	if bundle, exists := languageBundles[langCode]; exists {
		return bundle
	}

	// 回退到备用语言
	if globalConfig.FallbackLanguage != "" {
		if fallbackBundle, exists := languageBundles[globalConfig.FallbackLanguage]; exists {
			return fallbackBundle
		}
	}

	// 最后返回默认语言包
	return languageBundles[globalConfig.DefaultLanguage]
}

// T 获取翻译文本（带格式化参数）
// Get translated text with formatting arguments
func T(ctx *gin.Context, key string, args ...interface{}) string {
	if globalConfig.EnableCache {
		if cached := getCachedTranslation(key, args...); cached != "" {
			return cached
		}
	}

	bundle := ctx.MustGet(contextKey).(*Bundle)
	translatedText := bundle.translateKey(key, args...)

	if globalConfig.EnableCache {
		cacheTranslation(key, translatedText, args...)
	}

	return translatedText
}

// TPlural 获取复数形式翻译
// Get plural-form translation
func TPlural(ctx *gin.Context, count int, key string, args ...interface{}) string {
	bundle := ctx.MustGet(contextKey).(*Bundle)
	return bundle.handlePluralization(count, key, args...)
}

// TDate 获取本地化日期格式
// Get localized date format
func TDate(ctx *gin.Context, timestamp time.Time) string {
	bundle := ctx.MustGet(contextKey).(*Bundle)
	return bundle.formatLocalizedDate(timestamp)
}

// translateKey 执行翻译逻辑
func (b *Bundle) translateKey(key string, args ...interface{}) string {
	translation, exists := b.translations[key]
	if !exists {
		// 尝试回退语言包
		if fallbackBundle := getLanguageBundle(globalConfig.FallbackLanguage); fallbackBundle != nil {
			return fallbackBundle.translateKey(key, args...)
		}
		return key // 最后返回键名本身
	}

	baseText := conv.String(translation)
	if len(args) > 0 {
		return fmt.Sprintf(baseText, args...)
	}
	return baseText
}

// handlePluralization 处理复数形式翻译
func (b *Bundle) handlePluralization(count int, key string, args ...interface{}) string {
	// 优先使用自定义复数规则
	if globalConfig.PluralRules != nil {
		return globalConfig.PluralRules(b.languageTag, count, key, args...)
	}

	// 默认复数规则处理
	pluralKey := key
	if count != 1 {
		pluralKey = key + ".plural"
	}

	if pluralText, exists := b.translations[pluralKey]; exists {
		return fmt.Sprintf(conv.String(pluralText), append([]interface{}{count}, args...)...)
	}
	return b.translateKey(key, args...)
}

// formatLocalizedDate 格式化本地化日期
func (b *Bundle) formatLocalizedDate(timestamp time.Time) string {
	return timestamp.Format(b.dateTimeFormat)
}

// 缓存相关操作
// getCachedTranslation 从缓存获取翻译
func getCachedTranslation(key string, args ...interface{}) string {
	cacheKey := generateCacheKey(key, args...)
	translationCache.RLock()
	defer translationCache.RUnlock()
	return translationCache.cacheEntries[cacheKey]
}

// cacheTranslation 缓存翻译结果
func cacheTranslation(key, value string, args ...interface{}) {
	cacheKey := generateCacheKey(key, args...)
	translationCache.Lock()
	defer translationCache.Unlock()

	// LRU淘汰策略
	if len(translationCache.cacheKeys) >= translationCache.maxCapacity {
		oldestKey := translationCache.cacheKeys[0]
		delete(translationCache.cacheEntries, oldestKey)
		translationCache.cacheKeys = translationCache.cacheKeys[1:]
	}

	translationCache.cacheEntries[cacheKey] = value
	translationCache.cacheKeys = append(translationCache.cacheKeys, cacheKey)
}

// generateCacheKey 生成缓存键
func generateCacheKey(key string, args ...interface{}) string {
	return fmt.Sprintf("%s|%v", key, args)
}

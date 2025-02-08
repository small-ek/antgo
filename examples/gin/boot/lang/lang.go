package lang

import (
	"fmt"
	"github.com/small-ek/antgo/i18n"
	"github.com/small-ek/antgo/os/alog"
	"go.uber.org/zap"
)

func Register() {
	dateTimeLayouts := map[string]string{
		"en":    "January 02, 2006",      // 英文格式：月 日, 年
		"zh-CN": "2006年01月02日",           // 中文格式：年月日
		"ja":    "2006年01月02日 15時04分05秒", // 日文格式：年月日 时分秒
		"fr":    "02/01/2006",            // 法文格式：日/月/年
		"de":    "02.01.2006",            // 德文格式：日.月.年
	}
	config := i18n.Config{
		DefaultLang:     "en",
		TranslationsDir: "./lang",
		SupportedLangs:  []string{"en", "zh-CN"},
		CacheEnabled:    true,
		MaxCacheSize:    1000,
		DateTimeLayouts: dateTimeLayouts,
	}
	if err := i18n.New(config); err != nil {
		alog.Write.Panic("Init failed: %v", zap.Error(err))
	}
}

func customPluralRule(lang string, n int, key string, args ...interface{}) string {
	switch lang {
	case "en": // 英文复数规则
		if n == 1 {
			return fmt.Sprintf("%s.singular", key) // 单数形式
		}
		return fmt.Sprintf("%s.plural", key) // 复数形式
	case "zh": // 中文没有复数形式
		return key
	default: // 默认规则
		return key
	}
}

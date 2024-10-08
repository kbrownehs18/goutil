package env

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/viper"
)

func InitConfig[T any](config *T, path ...string) {
	configName := "config"
	configType := "yaml" // 暂时只支持yaml
	configPath := "."

	if len(path) > 0 {
		configPath = filepath.Dir(path[0])
		configName = filepath.Base(path[0])
		fileExt := filepath.Ext(configName)
		configName = strings.TrimSuffix(configName, fileExt)
	}

	viper.AddConfigPath(configPath)
	viper.SetConfigType(configType)
	viper.SetConfigName(configName)

	if err := viper.ReadInConfig(); err != nil {
		panic(fmt.Errorf("init config error: %v %s", config, err))
	}

	// 替换占位符
	replacePlaceholders()

	if err := viper.Unmarshal(config); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %v %s", config, err))
	}
}

// replacePlaceholders 用于替换配置文件中的 ${VAR:DEFAULT} 占位符
func replacePlaceholders() {
	settings := viper.AllSettings()
	replacePlaceholdersRecursive(settings, "")
}

func replacePlaceholdersRecursive(settings map[string]interface{}, parentKey string) {
	for key, value := range settings {
		fullKey := key
		if parentKey != "" {
			fullKey = parentKey + "." + key
		}

		switch v := value.(type) {
		case string:
			// 对字符串类型进行占位符替换
			viper.Set(fullKey, parseEnvPlaceholder(v))
		case map[string]interface{}:
			// 递归处理嵌套的 map
			replacePlaceholdersRecursive(v, fullKey)
		case []interface{}:
			// 处理列表类型
			for i, elem := range v {
				elemKey := fmt.Sprintf("%s[%d]", fullKey, i)
				if elemMap, ok := elem.(map[string]interface{}); ok {
					// 如果列表中的元素是 map，则递归处理
					replacePlaceholdersRecursive(elemMap, elemKey)
				} else if elemStr, ok := elem.(string); ok {
					// 如果列表中的元素是字符串，则替换占位符
					v[i] = parseEnvPlaceholder(elemStr)
				}
			}
			viper.Set(fullKey, v) // 更新替换后的列表
		default:
			// 其他类型（如 int, bool, float 等）不做处理
			viper.Set(fullKey, value)
		}
	}
}

// parseEnvPlaceholder 解析 ${VAR:DEFAULT} 语法
func parseEnvPlaceholder(value string) string {
	if strings.Index(value, "${") == 0 && strings.LastIndex(value, "}") == (len(value)-1) {
		parts := strings.SplitN(strings.Trim(value, "${}"), ":", 2)
		n := len(parts)
		if n == 1 {
			return os.Getenv(parts[0])
		} else if n == 2 {
			envVar, defaultValue := parts[0], parts[1]
			if envVal := os.Getenv(envVar); envVal != "" {
				// 如果获取到环境变量，那么直接返回
				return envVal
			}
			return defaultValue
		}
	}

	return value
}

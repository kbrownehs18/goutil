package env

import (
	"fmt"
	f "github.com/kbrownehs18/goutil/file"
	"github.com/spf13/viper"
	"os"
	"strings"
)

func InitConfig[T any](config *T, path ...string) {
	configName := "config"
	configType := "yaml" // 暂时只支持yaml
	configPath := "."

	fileInfo := f.Info{Path: configPath, Type: configType, Name: configName}

	if len(path) > 0 {
		fileInfo = f.ParseFilePath(path[0])
	}

	v := NewViper(fileInfo)

	if err := v.ReadInConfig(); err != nil {
		panic(fmt.Errorf("init config error: %v %s", config, err))
	}

	// 替换占位符
	replacePlaceholders(v)

	if err := v.Unmarshal(config); err != nil {
		panic(fmt.Errorf("unable to decode into struct: %v %s", config, err))
	}
}

// replacePlaceholders 用于替换配置文件中的 ${VAR:DEFAULT} 占位符
func replacePlaceholders(v *viper.Viper) {
	settings := v.AllSettings()
	replacePlaceholdersRecursive(v, settings, "")
}

func replacePlaceholdersRecursive(v *viper.Viper, settings map[string]interface{}, parentKey string) {
	for key, value := range settings {
		fullKey := key
		if parentKey != "" {
			fullKey = parentKey + "." + key
		}

		switch val := value.(type) {
		case string:
			// 对字符串类型进行占位符替换
			v.Set(fullKey, parseEnvPlaceholder(val))
		case map[string]interface{}:
			// 递归处理嵌套的 map
			replacePlaceholdersRecursive(v, val, fullKey)
		case []interface{}:
			// 处理列表类型
			for i, elem := range val {
				elemKey := fmt.Sprintf("%s[%d]", fullKey, i)
				if elemMap, ok := elem.(map[string]interface{}); ok {
					// 如果列表中的元素是 map，则递归处理
					replacePlaceholdersRecursive(v, elemMap, elemKey)
				} else if elemStr, ok := elem.(string); ok {
					// 如果列表中的元素是字符串，则替换占位符
					val[i] = parseEnvPlaceholder(elemStr)
				}
			}
			v.Set(fullKey, val) // 更新替换后的列表
		default:
			// 其他类型（如 int, bool, float 等）不做处理
			v.Set(fullKey, value)
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

func NewViperFromPath(path string) *viper.Viper {
	fileInfo := f.ParseFilePath(path)
	return NewViper(fileInfo)
}

func NewViper(fileInfo f.Info) *viper.Viper {
	v := viper.New()

	v.AddConfigPath(fileInfo.Path)
	v.SetConfigType(fileInfo.Type)
	v.SetConfigName(fileInfo.Name)

	return v
}

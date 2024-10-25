package file

import (
	"os"
	"path/filepath"
	"strings"
)

// Info for viper file
type Info struct {
	Path string
	Type string
	Name string
}

func ParseFilePath(path string) Info {
	configPath := filepath.Dir(path)
	configName := filepath.Base(path)
	fileExt := filepath.Ext(configName)
	configName = strings.TrimSuffix(configName, fileExt)
	configType := strings.TrimPrefix(fileExt, ".")

	return Info{configPath, configType, configName}
}

// FileExists 判断是否存在文件
func FileExists(filePath string) bool {
	_, err := os.Stat(filePath)
	if os.IsNotExist(err) {
		return false // 文件不存在
	}
	return err == nil // 文件存在
}

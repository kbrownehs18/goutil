package file

import (
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

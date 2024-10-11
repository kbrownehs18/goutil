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

	return Info{configPath, fileExt, configName}
}

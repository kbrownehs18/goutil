package test

import (
	"testing"

	"github.com/kbrownehs18/goutil/env"
)

type Config struct {
	ServerName string `mapstructure:"server_name"`
	App        App
}

type App struct {
	Name      string
	Age       int
	Addresses []string
	Old       bool
	Weight    float32
	Server    Server
}

type Server struct {
	Port int
}

func TestInitConfig(t *testing.T) {
	config := &Config{}
	env.InitConfig[Config](config)

	t.Logf("config: %v", config)
}

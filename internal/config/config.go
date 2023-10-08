package config

import (
	"errors"
	"fmt"
	"os"
	"time"

	"github.com/knadh/koanf"
	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
)

var k = koanf.New(".")

type Config struct {
	Env        string `koanf:"env"`
	LogsLevel  uint8  `koanf:"logs_level"`
	HTTPServer `koanf:"http_server"`
}

type HTTPServer struct {
	Address       string        `koanf:"address"`
	Timeout       time.Duration `koanf:"timeout"`
	IdleTimeout   time.Duration `koanf:"idle_timeout"`
	Mode          string        `koanf:"mode"`
	TemplatesPath string        `koanf:"templates_path"`
}

func New() (*Config, error) {
	configPath, configExists := os.LookupEnv("CONFIG_PATH")
	if !configExists || configPath == "" {
		return nil, errors.New("CONFIG_PATH is not set")
	}

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	if err := k.Load(file.Provider(configPath), yaml.Parser()); err != nil {
		return nil, fmt.Errorf("error loading config: %v", err)
	}

	var cfg Config
	if err := k.Unmarshal("", &cfg); err != nil {
		return nil, fmt.Errorf("error parsing config: %v", err)
	}

	return &cfg, nil
}

package config

import (
	"errors"
	"fmt"
	"os"

	"gopkg.in/yaml.v3"
)

type Config struct {
	Flags
	Exclusions []string `yaml:"exclusions"`
}

type Flags struct {
	ConfigPath string
	Chart      string
	Version    string
	ValuesFile string
	Namespace  string
}

func Load(path string) (Config, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		if errors.Is(err, os.ErrNotExist) {
			return Config{}, nil
		}
		return Config{}, fmt.Errorf("failed to read config file: %w", err)
	}

	var config Config

	if err := yaml.Unmarshal(data, &config); err != nil {
		return Config{}, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	return config, nil
}

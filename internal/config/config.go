package config

import (
	"fmt"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type Config struct {
	TargetURL    string            `yaml:"target_url"`
	Method       string            `yaml:"method"`
	Headers      map[string]string `yaml:"headers"`
	Concurrency  uint              `yaml:"concurrency"`
	RequestCount *uint             `yaml:"request_count"`
	Duration     *time.Duration    `yaml:"duration"`
	HTTPTimeout  time.Duration     `yaml:"http_timeout"`
}

func LoadConfig(filePath string) (c Config, e error) {
	file, err := os.Open(filePath)
	if err != nil {
		return c, fmt.Errorf("failed to open config file: %w", err)
	}
	defer file.Close()

	var cfg Config
	decoder := yaml.NewDecoder(file)
	if err := decoder.Decode(&cfg); err != nil {
		return c, fmt.Errorf("failed to decode config: %w", err)
	}

	return cfg, nil
}

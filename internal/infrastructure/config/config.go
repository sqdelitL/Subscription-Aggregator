package config

import (
	"errors"
	"fmt"
	"os"

	"github.com/caarlos0/env/v6"
	"gopkg.in/yaml.v3"
)

type Config struct {
	DataBase struct {
		Name     string `env:"DB,required"`
		UserName string `env:"DB_USER_NAME,required"`
		Password string `env:"DB_PASS,required"`
		Host     string `env:"DB_HOST,required"`
		Port     string `env:"DB_PORT,required"`
		Retry    struct {
			RetryDelaySeconds        int64   `yaml:"retry_delay_seconds"`
			MaxRetries               int     `yaml:"max_retries"`
			BackoffMultiplierSeconds float64 `yaml:"backoff_multiplier_seconds"`
			MaxIntervalSeconds       int64   `yaml:"max_interval_seconds"`
		} `yaml:"retry"`
		MaxIdleConns    int `yaml:"max_idle_conns"`
		MaxOpenConns    int `yaml:"max_open_conns"`
		ConnMaxLifetime int `yaml:"conn_max_lifetime"`
	} `yaml:"data_base"`

	ServerPort string `yaml:"server_port"`
	LogLevel   string `yaml:"log_level"`
}

func New() (*Config, error) {
	var data []byte
	var err error

	configPath := os.Getenv("CONFIG_PATH")
	if configPath == "" {
		return nil, errors.New("CONFIG_PATH environment variable not set")
	}
	data, err = os.ReadFile(configPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read config file %s: %w", configPath, err)
	}

	var config Config
	if err := yaml.Unmarshal(data, &config); err != nil {
		return nil, fmt.Errorf("failed to unmarshal config: %w", err)
	}

	if err := env.Parse(&config); err != nil {
		return nil, fmt.Errorf("failed to parse env: %w", err)
	}

	return &config, nil
}

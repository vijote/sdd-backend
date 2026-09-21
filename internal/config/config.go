package config

import (
	"github.com/kelseyhightower/envconfig"
)

// Config holds application configuration loaded from environment variables.
type Config struct {
	Port     int    `envconfig:"PORT" default:"8080"`
	LogLevel string `envconfig:"LOG_LEVEL" default:"info"`
}

// Load reads configuration from the environment, applying struct-tag defaults.
func Load() (*Config, error) {
	var cfg Config
	if err := envconfig.Process("", &cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

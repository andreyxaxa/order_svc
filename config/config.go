package config

import (
	"fmt"

	"github.com/caarlos0/env/v11"
)

type (
	Config struct {
	}

	HTTP struct {
		Port string `env:"HTTP_PORT,required"`
	}

	Log struct {
		Level string `env:"LOG_LEVEL,required"`
	}

	PG struct {
		PoolMax int    `env:"PG_POOL_MAX,required"`
		URL     string `env:"PG_URL,required"`
	}
)

func NewConfig() (*Config, error) {
	cfg := &Config{}
	if err := env.Parse(cfg); err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	return cfg, nil
}

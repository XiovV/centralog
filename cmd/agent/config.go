package main

import (
	"errors"
	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port                 string   `env:"PORT" env-default:"8080"`
	Environment          string   `env:"ENV" env-default:"PRODUCTION"`
	MonitorAllContainers bool     `env:"MONITOR_ALL_CONTAINERS" env-default:"false"`
	Containers           []string `env:"CONTAINERS" env-default:""`
}

func NewConfig() (*Config, error) {
	var cfg Config

	err := cleanenv.ReadEnv(&cfg)
	if err != nil {
		return nil, err
	}

	if !cfg.MonitorAllContainers && len(cfg.Containers) == 0 {
		return nil, errors.New("MONITOR_ALL_CONTAINERS or CONTAINERS environment variables should be specified")
	}

	if cfg.MonitorAllContainers && len(cfg.Containers) > 0 {
		return nil, errors.New("MONITOR_ALL_CONTAINERS and CONTAINERS environment variables should not be set simultaneously")
	}

	return &cfg, err
}

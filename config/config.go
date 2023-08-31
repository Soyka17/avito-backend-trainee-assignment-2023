package config

import (
	"fmt"
	"github.com/ilyakaznacheev/cleanenv"
)

type (
	// Config -.
	Config struct {
		HTTP     `yaml:"http"`
		Postgres `yaml:"postgres"`
		Log      `yaml:"logger"`
	}
	// HTTP -.
	HTTP struct {
		Port    string `env-required:"true" yaml:"port" env:"HTTP_PORT"`
		SiteUrl string `env-required:"true" yaml:"site_url" env:"HTTP_SITE_URL"`
	}
	// Postgres -.
	Postgres struct {
		Url string `env-required:"true" yaml:"url" env:"DATABASE_URL"`
	}

	// Log -.
	Log struct {
		Level string `env-required:"true" yaml:"log_level"   env:"LOG_LEVEL"`
	}
)

// NewConfig returns app config.
func NewConfig() (*Config, error) {
	cfg := &Config{}

	err := cleanenv.ReadConfig("./config/config.yml", cfg)
	if err != nil {
		return nil, fmt.Errorf("config error: %w", err)
	}

	err = cleanenv.ReadEnv(cfg)
	if err != nil {
		return nil, err
	}

	return cfg, nil
}

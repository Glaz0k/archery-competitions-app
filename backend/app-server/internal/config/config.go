package config

import (
	"app-server/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port      int             `yaml:"PORT" env:"PORT"`
	Postgres  postgres.Config `yaml:"POSTGRES"`
	SecretKey string          `yaml:"SECRET_KEY" env:"SECRET_KEY"`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadEnv(&cfg); err != nil {
		return nil, err
	}
	return &cfg, nil
}

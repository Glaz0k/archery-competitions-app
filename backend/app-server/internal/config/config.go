package config

import (
	"app-server/pkg/postgres"

	"github.com/ilyakaznacheev/cleanenv"
)

type Config struct {
	Port     int             `yaml:"PORT"`
	Postgres postgres.Config `yaml:"POSTGRES"`
}

func New() (*Config, error) {
	var cfg Config
	if err := cleanenv.ReadConfig("./config/config.yaml", &cfg); err != nil {
		return nil, err
	}

	return &cfg, nil
}

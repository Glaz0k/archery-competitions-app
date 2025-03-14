package postgres

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" `
	Port     int    `yaml:"POSTGRES_PORT" `
	User     string `yaml:"POSTGRES_USER" `
	Password string `yaml:"POSTGRES_PASSWORD" `
	Database string `yaml:"POSTGRES_DB"`
}

func New(config Config) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)
	conn, err := pgx.Connect(context.Background(), connString)

	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}
	return conn, nil
}

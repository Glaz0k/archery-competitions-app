package postgres

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" envDefault:"5434"`
	User     string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" envDefault:"root"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD" `
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" envDefault:"BowCompetitions"`
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

package postgres

import (
	"context"
	"fmt"

	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5/pgxpool"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" env:"POSTGRES_HOST" envDefault:"localhost"`
	Port     int    `yaml:"POSTGRES_PORT" env:"POSTGRES_PORT" envDefault:"5434"`
	User     string `yaml:"POSTGRES_USER" env:"POSTGRES_USER" envDefault:"root"`
	Password string `yaml:"POSTGRES_PASSWORD" env:"POSTGRES_PASSWORD"`
	Database string `yaml:"POSTGRES_DB" env:"POSTGRES_DB" envDefault:"BowCompetitions"`
	PoolSize int32  `yaml:"POSTGRES_POOL_SIZE" env:"POSTGRES_POOL_SIZE" envDefault:"10"`
}

func New(config Config) (*pgxpool.Pool, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
	)

	poolConfig, err := pgxpool.ParseConfig(connString)
	if err != nil {
		return nil, fmt.Errorf("failed to parse pool config: %w", err)
	}

	poolConfig.MaxConns = config.PoolSize

	pool, err := pgxpool.NewWithConfig(context.Background(), poolConfig)
	if err != nil {
		return nil, fmt.Errorf("could not create connection pool: %w", err)
	}

	err = pool.Ping(context.Background())
	if err != nil {
		return nil, fmt.Errorf("database connection failed: %w", err)
	}

	return pool, nil
}

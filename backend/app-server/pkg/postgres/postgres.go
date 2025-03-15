package postgres

import (
	"context"
	"fmt"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"
	"github.com/jackc/pgx/v5"

	"app-server/pkg/logger"
)

type Config struct {
	Host     string `yaml:"POSTGRES_HOST" `
	Port     int    `yaml:"POSTGRES_PORT" `
	User     string `yaml:"POSTGRES_USER" `
	Password string `yaml:"POSTGRES_PASSWORD" `
	Database string `yaml:"POSTGRES_DB"`
	//PoolMaxConns int    `yaml:"POSTGRES_POOL_MAX_CONNS"`
	//PoolMinConns int    `yaml:"POSTGRES_POOL_MIN_CONNS"`
}

func New(config Config) (*pgx.Conn, error) {
	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable", //&pool_max_conns=%d&pool_min_conns=%d",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		//config.PoolMaxConns,
		//config.PoolMinConns,
	)
	conn, err := pgx.Connect(context.Background(), connString)

	if err != nil {
		return nil, fmt.Errorf("could not connect to postgres: %w", err)
	}

	err = CreateMigration(connString)
	if err != nil {
		return nil, fmt.Errorf("could not create migrations: %w", err)
	}
	return conn, nil
}

func CreateMigration(connString string) error {
	m, err := migrate.New("file://migrations", connString)
	if err != nil {
		return err
	}

	version, dirty, err := m.Version()
	if err != nil && err != migrate.ErrNilVersion {
		return err
	}

	if dirty {
		logger.Logger.Fatal(fmt.Sprintf("Database is dirty at version %d. Forcing version.", version))
		err = m.Force(int(version))
		if err != nil {
			return fmt.Errorf("could not force version: %w", err)
		}
	}

	err = m.Up()
	if err != nil && err != migrate.ErrNoChange {
		return err
	}
	return nil
}

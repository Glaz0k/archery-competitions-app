package db

import (
	"errors"
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	"app-server/pkg/logger"
)

func CreateMigration(connString string) error {
	m, err := migrate.New("file://./cmd/migrate/migrations", connString)
	if err != nil {
		return fmt.Errorf("could not create migration: %w", err)
	}

	version, dirty, err := m.Version()
	if err != nil && !errors.Is(err, migrate.ErrNilVersion) {
		return fmt.Errorf("could not get migration version: %w", err)
	}

	if dirty {
		logger.Logger.Error(fmt.Sprintf("Database is dirty at version %d. Forcing version.", version))
		err = m.Force(int(version))
		if err != nil {
			return fmt.Errorf("could not force version %d: %w", version, err)
		}
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("could not apply migrations: %w", err)
	}
	return nil
}

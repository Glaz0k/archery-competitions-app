package db

import (
	"fmt"

	"github.com/golang-migrate/migrate/v4"

	"app-server/pkg/logger"
)

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

package migrate

import (
	"fmt"

	"app-server/internal/config"
	"app-server/pkg/logger"
	"app-server/pkg/postgres"
)

func main() {
	cfg, err := config.New()
	if err != nil {
		panic(err)
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable",
		cfg.Postgres.User,
		cfg.Postgres.Password,
		cfg.Postgres.Host,
		cfg.Postgres.Port,
		cfg.Postgres.Database,
	)

	err = postgres.CreateMigration(connString)
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to apply migrations: %v", err))
	}

	logger.Logger.Info("Migrations applied successfully")
}

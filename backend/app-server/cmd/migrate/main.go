package migrate

import (
	"fmt"

<<<<<<< HEAD
	"app-server/cmd/migrate/db"
	"app-server/internal/config"
	"app-server/pkg/logger"
=======
	"app-server/internal/config"
	"app-server/pkg/logger"
	"app-server/pkg/postgres"
>>>>>>> a073ba83ee8dde11bdcf876b336b1a2cf419b4ce
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

<<<<<<< HEAD
	err = db.CreateMigration(connString)
=======
	err = postgres.CreateMigration(connString)
>>>>>>> a073ba83ee8dde11bdcf876b336b1a2cf419b4ce
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to apply migrations: %v", err))
	}

	logger.Logger.Info("Migrations applied successfully")
}

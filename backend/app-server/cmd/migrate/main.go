package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"

	"go.uber.org/zap"

	"app-server/cmd/migrate/db"
	"app-server/internal/config"
	"app-server/pkg/logger"
	"app-server/pkg/postgres"
)

func main() {
	ctx := context.Background()
	ctx, stop := signal.NotifyContext(ctx, os.Interrupt)
	defer stop()

	err := logger.New()
	if err != nil {
		panic(err)
	}

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

	_, err = postgres.New(cfg.Postgres)
	if err != nil {
		logger.Logger.Fatal("failed to connect to database", zap.Error(err))
	}

	err = db.CreateMigration(connString)
	if err != nil {
		logger.Logger.Fatal(fmt.Sprintf("Failed to apply migrations: %v", err))
	}

	logger.Logger.Info("Migrations applied successfully")
}

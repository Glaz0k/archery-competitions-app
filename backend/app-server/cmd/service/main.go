package main

import (
	"context"
	"os"
	"os/signal"

	"app-server/internal/config"
	"app-server/internal/server"
	"app-server/internal/server/handlers"
	"app-server/pkg/logger"
	"app-server/pkg/postgres"

	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
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
		logger.Logger.Fatal("failed to read config", zap.Error(err))
	}

	conn, err := postgres.New(cfg.Postgres)
	if err != nil {
		logger.Logger.Fatal("failed to connect to database", zap.Error(err))
	}

	defer func(conn *pgx.Conn, ctx context.Context) {
		err := conn.Close(ctx)
		if err != nil {
			logger.Logger.Fatal("failed to close connection", zap.Error(err))
		}
	}(conn, context.Background())

	handlers.InitDB(conn)

	srv := server.New(*cfg, logger.Logger)

	go func() {
		logger.Logger.Fatal("server listen and serve:", zap.Error(srv.Run()))
	}()

	select {
	case <-ctx.Done():
		logger.Logger.Info("Server stopped")
	}
}

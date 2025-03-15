package postgres

import (
	"context"
	"fmt"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	config := Config{
		Host:         "localhost",
		Port:         5432,
		User:         "root",
		Password:     "root_password",
		Database:     "BowCompetiotions",
		PoolMaxConns: 10,
		PoolMinConns: 1,
	}

	conn, err := New(config)
	require.NoError(t, err, "Failed to connect to database")
	defer conn.Close(context.Background())

	var tableExists bool
	err = conn.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'cups'
		);
	`).Scan(&tableExists)
	require.NoError(t, err, "Failed to check if table exists")
	require.True(t, tableExists, "Table 'users' should exist")
}

func TestCreateMigration(t *testing.T) {
	config := Config{
		Host:         "localhost",
		Port:         5432,
		User:         "root",
		Password:     "root_password",
		Database:     "BowCompetiotions",
		PoolMaxConns: 10,
		PoolMinConns: 1,
	}

	connString := fmt.Sprintf("postgres://%s:%s@%s:%d/%s?sslmode=disable&pool_max_conns=%d&pool_min_conns=%d",
		config.User,
		config.Password,
		config.Host,
		config.Port,
		config.Database,
		config.PoolMaxConns,
		config.PoolMinConns,
	)

	err := CreateMigration(connString)
	require.NoError(t, err, "Failed to apply migrations")

	conn, err := pgx.Connect(context.Background(), connString)
	require.NoError(t, err, "Failed to connect to database")
	defer conn.Close(context.Background())

	var tableExists bool
	err = conn.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'cups'
		);
	`).Scan(&tableExists)
	require.NoError(t, err, "Failed to check if table exists")
	require.True(t, tableExists, "Table 'users' should exist")
}

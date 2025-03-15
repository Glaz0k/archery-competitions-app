package postgres

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	config := Config{
		Host:     "localhost",
		Port:     5432,
		User:     "root",
		Password: "root_password",
		Database: "postgres",
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
	require.False(t, tableExists, "Table 'cups' should exist")

	err = conn.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'users'
		);
	`).Scan(&tableExists)
	require.NoError(t, err, "Failed to check if table exists")
	require.False(t, tableExists, "Table 'users' should not exist")

	err = conn.QueryRow(context.Background(), `
		SELECT EXISTS (
			SELECT FROM information_schema.tables 
			WHERE table_name = 'shots'
		);
	`).Scan(&tableExists)
	require.NoError(t, err, "Failed to check if table exists")
	require.False(t, tableExists, "Table 'shots' should exist")
}

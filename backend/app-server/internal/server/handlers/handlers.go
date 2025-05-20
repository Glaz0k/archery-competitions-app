package handlers

import (
	"github.com/jackc/pgx/v5/pgxpool"
)

var dbPool *pgxpool.Pool

func InitDB(pool *pgxpool.Pool) {
	dbPool = pool
}

package handlers

import (
	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func InitDB(dbConn *pgx.Conn) {
	conn = dbConn
}

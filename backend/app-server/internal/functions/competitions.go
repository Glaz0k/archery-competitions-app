package functions

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func InitDB(dbConn *pgx.Conn) {
	conn = dbConn
}

func AddCup(w http.ResponseWriter, r *http.Request) {
	var cup models.Cup
	err := json.NewDecoder(r.Body).Decode(&cup)
	if err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO cups (title, address, season) VALUES ($1, $2, $3)", cup.Title, cup.Address, cup.Season)
	if err != nil {
		log.Fatalf("Unable to insert data: %v\n", err)
	}
}

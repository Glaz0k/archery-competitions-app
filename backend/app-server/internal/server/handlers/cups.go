package handlers

import (
	"app-server/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

func CreateCup(w http.ResponseWriter, r *http.Request) {
	var cup models.Cup
	err := json.NewDecoder(r.Body).Decode(&cup)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO cups (title, address, season) VALUES ($1, $2, $3)", cup.Title, cup.Address, cup.Season)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetCup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cupID := vars["cup_id"]

	if cupID == "" {
		http.Error(w, "cup_id is required", http.StatusBadRequest)
		return
	}

	var cup models.Cup
	query := `SELECT id, title, address, season FROM cups WHERE id = $1`
	err := conn.QueryRow(context.Background(), query, cupID).Scan(
		&cup.ID,
		&cup.Title,
		&cup.Address,
		&cup.Season,
	)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "cup not found", http.StatusNotFound)
		} else {
			http.Error(w, fmt.Sprintf("unable to get data: %v", err), http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(cup)
}

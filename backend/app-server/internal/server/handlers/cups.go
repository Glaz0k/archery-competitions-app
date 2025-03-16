package handlers

import (
	"app-server/internal/models"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

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
		http.Error(w, fmt.Sprintf("unable to insert data: %v\n", err), http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetCup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cupId := vars["cup_id"]
	cupID, err := strconv.Atoi(cupId)
	if err != nil {
		http.Error(w, "invalid cup_id", http.StatusBadRequest)
	}

	var cup models.Cup
	query := `SELECT id, title, address, season FROM cups WHERE id = $1`
	err = conn.QueryRow(context.Background(), query, cupID).Scan(
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

func EditCup(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	cupID := vars["cup_id"]
	var cup models.Cup
	err := json.NewDecoder(r.Body).Decode(&cup)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	query := `UPDATE cups SET title = $1, address = $2, season = $3 WHERE id = $4`
	_, err = conn.Exec(context.Background(), query, cup.Title, cup.Address, cup.Season, cupID)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to update data: %v\n", err), http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

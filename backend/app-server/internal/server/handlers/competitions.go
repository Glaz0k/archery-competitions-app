package handlers

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

func CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var competition models.Competition
	err := json.NewDecoder(r.Body).Decode(&competition)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	cupId := mux.Vars(r)["cup_id"]
	cupID, err := strconv.Atoi(cupId)
	if err != nil {
		http.Error(w, "invalid cup_id", http.StatusBadRequest)
	}
	competition.CupID = cupID
	var exists bool
	queryCheck := `SELECT EXISTS(SELECT 1 FROM competitions WHERE cup_id = $1 AND stage = $2)`
	err = conn.QueryRow(context.Background(), queryCheck, competition.CupID, competition.Stage).Scan(&exists)
	if err != nil {
		http.Error(w, fmt.Sprintf("unable to check data existence: %v", err), http.StatusInternalServerError)
		return
	}

	if exists {
		http.Error(w, "competition already exists", http.StatusConflict)
		return
	}
	query := "INSERT INTO competitions (cup_id, stage, start_date, end_date, is_ended) VALUES ($1, $2, $3, $4, $5)"

	_, err = conn.Exec(context.Background(), query, competition.CupID, competition.Stage,
		competition.StartDate, competition.EndDate, competition.IsEnded)

	if err != nil {
		http.Error(w, "unable to insert data: %v\n", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func EditCompetition(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	competitionID := vars["competition_id"]
	var updateData struct {
		StartDate string `json:"start_date"`
		EndDate   string `json:"end_date"`
	}
	err := json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	query := `
        UPDATE competitions
        SET start_date = $1, end_date = $2
        WHERE id = $3
        RETURNING id, cup_id, stage, start_date, end_date, is_ended
    `
	var competition models.Competition
	err = conn.QueryRow(context.Background(), query, updateData.StartDate, updateData.EndDate, competitionID).Scan(
		&competition.ID,
		&competition.CupID,
		&competition.Stage,
		&competition.StartDate,
		&competition.EndDate,
		&competition.IsEnded,
	)

	if err != nil {
		http.Error(w, "unable to update data", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(competition)
}

package handlers

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var competition models.Competition
	err := json.NewDecoder(r.Body).Decode(&competition)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO competitions (cup_id, stage, start_date, end_date, is_ended) VALUES ($1, $2, $3, $4, $5)", competition.CupID, competition.Stage, competition.StartDate, competition.EndDate, competition.IsEnded)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

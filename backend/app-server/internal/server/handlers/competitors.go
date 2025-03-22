package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"net/http"

	"app-server/internal/models"
	"app-server/pkg/tools"
)

func CreateCompetitors(w http.ResponseWriter, r *http.Request) {
	var competitor models.Competitor
	err := json.NewDecoder(r.Body).Decode(&competitor)
	if err != nil {
		http.Error(w, "invalid decode", http.StatusBadRequest)
		return
	}

	var ok bool
	competitor.ID, ok = r.Context().Value("user_id").(int)
	if !ok {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	query := `INSERT INTO competitors (id, full_name, birth_date, identity, bow, rank, region, federation, club)
VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
ON CONFLICT (full_name, birth_date)
DO UPDATE SET
    id = EXCLUDED.id,
    identity = EXCLUDED.identity,
    bow = EXCLUDED.bow,
    rank = EXCLUDED.rank,
    region = EXCLUDED.region,
    federation = EXCLUDED.federation,
    club = EXCLUDED.club;`
	_, err = conn.Exec(context.Background(), query, competitor.ID, competitor.FullName, competitor.BirthDate, competitor.Identity, competitor.Bow, competitor.Rank, competitor.Region, competitor.Federation, competitor.Club)
	if err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func GetCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := tools.ParseParamToInt(r, "competitor_id")
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	var competitor models.Competitor
	query := `SELECT * FROM competitors WHERE id = $1`
	err = conn.QueryRow(context.Background(), query, competitorID).Scan(&competitor.ID, &competitor.FullName, &competitor.BirthDate, &competitor.Identity, &competitor.Bow, &competitor.Rank, &competitor.Region, &competitor.Federation, &competitor.Club)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "competitor not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(competitor); err != nil {
		http.Error(w, "database error", http.StatusInternalServerError)
		return
	}
}

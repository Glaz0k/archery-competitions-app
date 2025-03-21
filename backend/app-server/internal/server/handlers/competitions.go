package handlers

import (
	"app-server/internal/models"
	"app-server/pkg/tools"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"time"

	"github.com/gorilla/mux"
)

func CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var competition models.Competition
	err := json.NewDecoder(r.Body).Decode(&competition)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	cupID, err := tools.ParseParamToInt(r, "cup_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID ENDPOINT"})
	}
	competition.CupID = cupID
	if competition.EndDate.Before(time.Now()) {
		competition.IsEnded = true
	} else {
		competition.IsEnded = false
	}
	var exists bool
	queryCheck := `SELECT EXISTS(SELECT 1 FROM competitions WHERE cup_id = $1 AND stage = $2)`
	err = conn.QueryRow(context.Background(), queryCheck, competition.CupID, competition.Stage).Scan(&exists)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to check data existence: %v", err))
		return
	}
	if exists {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "EXISTS"})
		return
	}
	query := "INSERT INTO competitions (cup_id, stage, start_date, end_date, is_ended) VALUES ($1, $2, $3, $4, $5)"

	_, err = conn.Exec(context.Background(), query, competition.CupID, competition.Stage,
		competition.StartDate, competition.EndDate, competition.IsEnded)

	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusCreated, competition)
}

func GetAllCompetitions(w http.ResponseWriter, r *http.Request) {
	cupID, err := tools.ParseParamToInt(r, "cup_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID ENDPOINT"})
	}

	query := `SELECT id, cup_id, stage, start_date, end_date, is_ended FROM competitions WHERE cup_id = $1`
	rows, err := conn.Query(context.Background(), query, cupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()
	var competitions []models.Competition

	for rows.Next() {
		var competition models.Competition
		err = rows.Scan(
			&competition.ID,
			&competition.CupID,
			&competition.Stage,
			&competition.StartDate,
			&competition.EndDate,
			&competition.IsEnded,
		)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		competitions = append(competitions, competition)
	}

	if err = rows.Err(); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, competitions)
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
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	checkQuery := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, checkQuery, competitionID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "COMPETITION NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
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

	tools.WriteJSON(w, http.StatusOK, competition)
}

func EndCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID ENDPOINT"})
		return
	}

	checkQuery := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, checkQuery, competitionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "COMPETITION NOT FOUND"})
		return
	}
	var competition models.Competition

	query := `SELECT id, cup_id, stage, start_date, end_date, is_ended FROM competitions WHERE id = $1`
	err = conn.QueryRow(context.Background(), query, competitionID).Scan(&competition.ID, &competition.CupID, &competition.Stage, &competition.StartDate, &competition.EndDate, &competition.IsEnded)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if competition.IsEnded {
		tools.WriteJSON(w, http.StatusOK, competition)
		return
	}

	var totalCount, finalsEndCount int
	totalQuery := `SELECT COUNT(*) FROM individual_groups WHERE competition_id = $1`
	err = conn.QueryRow(context.Background(), totalQuery, competitionID).Scan(&totalCount)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	finalsEndQuery := `SELECT COUNT(*) FROM individual_groups WHERE competition_id = $1 AND state = 'finals_end'`
	err = conn.QueryRow(context.Background(), finalsEndQuery, competitionID).Scan(&finalsEndCount)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if totalCount != finalsEndCount {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	competition.IsEnded = true
	query = `UPDATE competitions SET is_ended = true WHERE id = $1`
	_, err = conn.Exec(context.Background(), query, competitionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, competition)
}

func GetCompetitorsFromCompetitionUser(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID ENDPOINT"})
		return
	}

	userID, ok := r.Context().Value("user_id").(uint)
	if !ok {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "UserID not found"})
		return
	}

	var registered bool
	queryCheck := `SELECT EXISTS (SELECT 1 FROM competitor_competition_details WHERE competition_id = $1 AND competitor_id = $2)`

	err = conn.QueryRow(context.Background(), queryCheck, competitionID, userID).Scan(&registered)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to check registration: %v", err))
		return
	}
	if !registered {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	//TODO
}

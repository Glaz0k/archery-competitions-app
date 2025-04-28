package handlers

import (
	"app-server/internal/models"
	"app-server/pkg/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/jackc/pgx/v5"
)

func CreateCup(w http.ResponseWriter, r *http.Request) {
	var cup models.Cup
	err := json.NewDecoder(r.Body).Decode(&cup)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}

	checkQuery := "SELECT id FROM cups WHERE title = $1 AND address = $2 AND season = $3"
	exists, err := tools.ExistsInDB(context.Background(), conn, checkQuery, cup.Title, cup.Address, cup.Season)
	if exists {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ALREADY EXISTS"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	err = conn.QueryRow(context.Background(),
		"INSERT INTO cups (title, address, season) VALUES ($1, $2, $3) RETURNING id",
		cup.Title, cup.Address, cup.Season).Scan(&cup.ID)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	tools.WriteJSON(w, http.StatusCreated, cup)
}

func GetCup(w http.ResponseWriter, r *http.Request) {
	cupID, err := tools.ParseParamToInt(r, "cup_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}
	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "BAD ACTION"})
		return
	}
	if role == "user" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil {
			tools.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "BAD ACTION"})
			return
		}

		checkQuery := `
            SELECT EXISTS (
                SELECT 1 
                FROM competitor_competition_details ccd
                JOIN competitions comp ON ccd.competition_id = comp.id
                WHERE ccd.competitor_id = $1 AND comp.cup_id = $2
            )`
		var isParticipant bool
		err = conn.QueryRow(context.Background(), checkQuery, userID, cupID).Scan(&isParticipant)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !isParticipant {
			tools.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "BAD ACTION"})
			return
		}
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
		if errors.Is(err, pgx.ErrNoRows) {
			tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		} else {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		}
		return
	}

	tools.WriteJSON(w, http.StatusOK, cup)
}

func GetAllCups(w http.ResponseWriter, r *http.Request) {
	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "BAD ACTION"})
		return
	}
	var query string
	var rows pgx.Rows
	if role == "user" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil {
			tools.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "BAD ACTION"})
			return
		}
		query = `
        SELECT 
            c.id, 
            c.title, 
            c.address, 
            c.season
        FROM 
            cups c
        JOIN 
            competitions comp ON c.id = comp.cup_id
        JOIN 
            competitor_competition_details ccd ON comp.id = ccd.competition_id
        WHERE 
            ccd.competitor_id = $1
        ORDER BY 
            c.season DESC, c.id`
		rows, err = conn.Query(context.Background(), query, userID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
	}
	if role == "admin" {
		query = `SELECT id, title, address, season FROM cups`
		rows, err = conn.Query(context.Background(), query)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
	}
	defer rows.Close()
	var cups []models.Cup

	for rows.Next() {
		var cup models.Cup
		err = rows.Scan(
			&cup.ID,
			&cup.Title,
			&cup.Address,
			&cup.Season,
		)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		cups = append(cups, cup)
	}

	if err = rows.Err(); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	tools.WriteJSON(w, http.StatusOK, cups)
}

func EditCup(w http.ResponseWriter, r *http.Request) {
	cupID, err := tools.ParseParamToInt(r, "cup_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID ENDPOINT"})
		return
	}
	var cup models.Cup
	err = json.NewDecoder(r.Body).Decode(&cup)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	checkQuery := `SELECT id FROM cups WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, checkQuery, cupID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "CUP NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	query := `UPDATE cups SET title = $1, address = $2, season = $3 WHERE id = $4`
	_, err = conn.Exec(context.Background(), query, cup.Title, cup.Address, cup.Season, cupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, cup)
}

func CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var competition models.Competition
	err := json.NewDecoder(r.Body).Decode(&competition)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
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
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("DATABSE ERROR: %v", err))
		return
	}
	if exists {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "EXISTS"})
		return
	}
	query := "INSERT INTO competitions (cup_id, stage, start_date, end_date, is_ended) VALUES ($1, $2, $3, $4, $5) RETURNING id"

	err = conn.QueryRow(context.Background(), query, competition.CupID, competition.Stage,
		competition.StartDate, competition.EndDate, competition.IsEnded).Scan(&competition.ID)

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

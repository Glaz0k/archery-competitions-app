package handlers

import (
	"app-server/internal/models"
	"app-server/pkg/tools"
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func CreateCup(w http.ResponseWriter, r *http.Request) {
	var cup models.Cup
	err := json.NewDecoder(r.Body).Decode(&cup)
	if err != nil {
		err = tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	checkQuery := "SELECT id FROM cups WHERE title = $1 AND address = $2 AND season = $3"
	exists, err := tools.ExistsInDB(context.Background(), conn, checkQuery, cup.Title, cup.Address, cup.Season)
	if exists {
		err = tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "CUP ALREADY EXISTS"})
		return
	}
	if err != nil {
		err = tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	err = conn.QueryRow(context.Background(),
		"INSERT INTO cups (title, address, season) VALUES ($1, $2, $3) RETURNING id",
		cup.Title, cup.Address, cup.Season).Scan(&cup.ID)
	if err != nil {
		err = tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	err = tools.WriteJSON(w, http.StatusCreated, cup)
}

func GetCup(w http.ResponseWriter, r *http.Request) {
	cupID, err := tools.ParseParamToInt(r, "cup_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
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
			tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "CUP NOT FOUND"})
		} else {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		}
		return
	}

	tools.WriteJSON(w, http.StatusOK, cup)
}

func GetAllCups(w http.ResponseWriter, r *http.Request) {
	query := `SELECT id, title, address, season FROM cups`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
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

	err = tools.WriteJSON(w, http.StatusOK, cups)
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

package handlers

import (
	"app-server/internal/models"
	"app-server/pkg/tools"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
)

func RegisterCompetitor(w http.ResponseWriter, r *http.Request) {
	var competitor models.Competitor
	err := json.NewDecoder(r.Body).Decode(&competitor)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return

	}
	userId, err := tools.GetUserIDFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	competitor.ID = userId
	exists, err := tools.ExistsInDB(context.Background(), conn, "SELECT id FROM competitors WHERE id = $1", competitor.ID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if exists {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "ALREADY EXISTS"})
		return
	}
	query := `INSERT INTO competitors (id, full_name, birth_date, identity, bow, rank, region, federation, club)
			  VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)`

	_, err = conn.Exec(context.Background(), query, competitor.ID, competitor.FullName, competitor.BirthDate,
		competitor.Identity, competitor.Bow, competitor.Rank, competitor.Region,
		competitor.Federation, competitor.Club)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusCreated, competitor)
}

func GetCompetitorsFromGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	rows, err := conn.Query(context.Background(), `
        SELECT c.id, c.full_name 
        FROM competitor_group_details cgd 
        JOIN competitors c ON cgd.competitor_id = c.id 
        WHERE cgd.group_id = $1`, groupId)

	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()

	var competitors []models.Competitor
	for rows.Next() {
		var competitor models.Competitor
		if err = rows.Scan(&competitor.ID, &competitor.FullName); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		competitors = append(competitors, competitor)
	}
	if err = rows.Err(); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if len(competitors) == 0 {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "COMPETITORS NOT FOUND"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, competitors)
}

func GetCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := tools.ParseParamToInt(r, "competitor_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	var competitor models.Competitor
	query := `SELECT id, full_name, birth_date, identity, bow, rank, region, federation, club FROM competitors WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, query, competitorID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	err = conn.QueryRow(context.Background(), query, competitorID).Scan(&competitor.ID, &competitor.FullName, &competitor.BirthDate,
		&competitor.Identity, &competitor.Bow, &competitor.Rank, &competitor.Region,
		&competitor.Federation, &competitor.Club)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, competitor)
}

func EditCompetitor(w http.ResponseWriter, r *http.Request, isAdmin bool) {
	competitorID, err := tools.ParseParamToInt(r, "competitor_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	if !isAdmin {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil {
			tools.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}
		if competitorID != userID {
			tools.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "PERMISSION DENIED"})
			return
		}
	}

	var competitor models.Competitor
	if err := json.NewDecoder(r.Body).Decode(&competitor); err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	query := `UPDATE competitors SET full_name = $1, birth_date = $2,
              identity = $3, bow = $4, rank = $5, region = $6,
              federation = $7, club = $8 WHERE id = $9`

	_, err = conn.Exec(context.Background(), query,
		competitor.FullName, competitor.BirthDate,
		competitor.Identity, competitor.Bow, competitor.Rank,
		competitor.Region, competitor.Federation,
		competitor.Club, competitorID)

	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, competitor)
}

func AdminEditCompetitor(w http.ResponseWriter, r *http.Request) {
	EditCompetitor(w, r, true)
}

func UserEditCompetitor(w http.ResponseWriter, r *http.Request) {
	EditCompetitor(w, r, false)
}

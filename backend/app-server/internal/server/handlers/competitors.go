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
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	competitor.ID = userId
	exists, err := tools.ExistsInDB(context.Background(), conn, "SELECT EXISTS(SELECT 1 FROM competitors WHERE id = $1)", competitor.ID)
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

func GetCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := tools.ParseParamToInt(r, "competitor_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	var competitor models.Competitor
	queryCheck := `SELECT EXISTS(SELECT 1 FROM competitors WHERE id = $1)`
	exists, err := tools.ExistsInDB(context.Background(), conn, queryCheck, competitorID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	query := `SELECT id, full_name, birth_date, identity, bow, rank, region, federation, club FROM competitors WHERE id = $1`
	err = conn.QueryRow(context.Background(), query, competitorID).Scan(&competitor.ID, &competitor.FullName, &competitor.BirthDate,
		&competitor.Identity, &competitor.Bow, &competitor.Rank, &competitor.Region,
		&competitor.Federation, &competitor.Club)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, competitor)
}

func EditCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := tools.ParseParamToInt(r, "competitor_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
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
		if competitorID != userID {
			tools.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	var competitor models.Competitor
	competitor.ID = competitorID
	if err = json.NewDecoder(r.Body).Decode(&competitor); err != nil {
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

func GetAllCompetitors(w http.ResponseWriter, r *http.Request) {
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	competitors := make([]models.Competitor, 0, 20)
	query := `SELECT id, full_name, birth_date, identity, bow, rank, region, federation, club FROM competitors`
	rows, err := conn.Query(context.Background(), query)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()
	for rows.Next() {
		var competitor models.Competitor
		err = rows.Scan(&competitor.ID, &competitor.FullName, &competitor.BirthDate,
			&competitor.Identity, &competitor.Bow, &competitor.Rank, &competitor.Region,
			&competitor.Federation, &competitor.Club)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		competitors = append(competitors, competitor)
	}
	tools.WriteJSON(w, http.StatusOK, competitors)
}

func DeleteCompetitor(w http.ResponseWriter, r *http.Request) {
	competitorID, err := tools.ParseParamToInt(r, "competitor_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()

	queryCheck := `SELECT EXISTS(SELECT 1 FROM competitors WHERE id = $1)`
	exists, err := tools.ExistsInDB(r.Context(), conn, queryCheck, competitorID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tx, err := conn.Begin(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE_ERROR"})
		return
	}
	defer tx.Rollback(r.Context())

	var isCompetitionEnded bool
	err = tx.QueryRow(r.Context(), `
		SELECT EXISTS(
			SELECT 1 FROM competitor_competition_details ccd
			JOIN competitions c ON ccd.competition_id = c.id
			WHERE ccd.competitor_id = $1 
			AND c.is_ended = true
		)`, competitorID).Scan(&isCompetitionEnded)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	if isCompetitionEnded {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	var isGroupActiveEnded bool
	err = tx.QueryRow(r.Context(), `SELECT EXISTS(
    SELECT 1 FROM competitor_group_details cgd
    JOIN individual_groups ig ON cgd.group_id = ig.id
    WHERE cgd.competitor_id = $1 
    AND ig.state != 'created'
)`, competitorID).Scan(&isGroupActiveEnded)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	if isGroupActiveEnded {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	_, err = tx.Exec(r.Context(), `
		DELETE FROM shoot_outs 
		WHERE place_id IN (
			SELECT id FROM sparring_places WHERE competitor_id = $1
		)`, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	_, err = tx.Exec(r.Context(), `
		DELETE FROM sparring_places 
		WHERE competitor_id = $1`, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	_, err = tx.Exec(r.Context(), `
		DELETE FROM qualification_rounds 
		WHERE section_id IN (
			SELECT id FROM qualification_sections WHERE competitor_id = $1
		)`, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	_, err = tx.Exec(r.Context(), `
		DELETE FROM qualification_sections 
		WHERE competitor_id = $1`, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	_, err = tx.Exec(r.Context(), `
		DELETE FROM competitor_group_details 
		WHERE competitor_id = $1`, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE_ERROR"})
		return
	}

	_, err = tx.Exec(r.Context(), `
		DELETE FROM competitor_competition_details 
		WHERE competitor_id = $1`, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	_, err = tx.Exec(r.Context(), `
		DELETE FROM competitors 
		WHERE id = $1`, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	err = tx.Commit(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

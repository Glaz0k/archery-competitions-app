package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"time"

	"app-server/internal/dto"
	"app-server/internal/models"
	"app-server/pkg/tools"

	"github.com/jackc/pgx/v5"
)

func EditCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	updateData := models.CompetitionUpdateData
	err = json.NewDecoder(r.Body).Decode(&updateData)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	queryCheck := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, queryCheck, competitionID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
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
		http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
		return
	}

	tools.WriteJSON(w, http.StatusOK, competition)
}

func DeleteCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	queryCheck := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, queryCheck, competitionID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer tx.Rollback(context.Background())

	err = deleteCompetition(tx, competitionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func deleteCompetition(tx pgx.Tx, competitionID int) error {
	var groups []int
	query := `SELECT id FROM individual_groups WHERE competition_id = $1`
	rows, err := tx.Query(context.Background(), query, competitionID)
	if err != nil {

	}
	defer rows.Close()

	for rows.Next() {
		var group int
		err = rows.Scan(&group)
		if err != nil {
			return fmt.Errorf("DATABASE ERROR: %v", err)
		}
		groups = append(groups, group)
	}

	for _, group := range groups {
		err = deleteAllGroupData(context.Background(), tx, group)
		if err != nil {
			return fmt.Errorf("DATABASE ERROR: %v", err)
		}
	}

	query = `DELETE FROM competitor_competition_details WHERE competition_id = $1`
	_, err = tx.Exec(context.Background(), query, competitionID)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}

	query = `DELETE FROM competitions WHERE id = $1`
	_, err = tx.Exec(context.Background(), query, competitionID)
	if err != nil {
		return fmt.Errorf("DATABASE ERROR: %v", err)
	}
	return nil
}

func EndCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	checkQuery := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, checkQuery, competitionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
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

func AddCompetitorCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	var competitorId dto.Competitor
	err = json.NewDecoder(r.Body).Decode(&competitorId)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	queryCheck := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, queryCheck, competitionID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var competitor models.Competitor
	query := `SELECT id, full_name, birth_date, identity, bow, rank, region, federation, club FROM competitors WHERE id = $1`
	err = conn.QueryRow(context.Background(), query, competitorId.CompetitorID).Scan(&competitor.ID, &competitor.FullName, &competitor.BirthDate,
		&competitor.Identity, &competitor.Bow, &competitor.Rank, &competitor.Region, &competitor.Federation, &competitor.Club)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
			return
		}
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var exist bool
	queryCheck = `SELECT EXISTS (SELECT 1 FROM competitor_competition_details WHERE competition_id = $1 AND competitor_id = $2)`
	err = conn.QueryRow(context.Background(), queryCheck, competitionID, competitorId.CompetitorID).Scan(&exist)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var competitionDetails dto.CompetitorCompetitionDetails
	if exist {
		query = `SELECT is_active, created_at FROM competitor_competition_details WHERE competition_id = $1`
		competitionDetails.CompetitionID = competitionID
		err = conn.QueryRow(context.Background(), query, competitionID).Scan(&competitionDetails.IsActive, &competitionDetails.CreatedAt)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
	} else {
		competitionDetails = dto.CompetitorCompetitionDetails{
			CompetitionID: competitionID,
			IsActive:      true,
			CreatedAt:     time.Now(),
		}
		query = `INSERT INTO competitor_competition_details (competition_id, competitor_id, is_active, created_at) VALUES ($1, $2, $3, $4)`
		_, err = conn.Exec(context.Background(), query, competitionID, competitor.ID, competitionDetails.IsActive, competitionDetails.CreatedAt)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
	}

	competitionDetails.Competitors = append(competitionDetails.Competitors, competitor)
	tools.WriteJSON(w, http.StatusOK, competitionDetails)
}

func GetCompetitorsFromCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}

	if role == "user" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		var registered bool
		queryCheck := `SELECT EXISTS (SELECT 1 FROM competitor_competition_details WHERE competition_id = $1 AND competitor_id = $2)`

		err = conn.QueryRow(context.Background(), queryCheck, competitionID, userID).Scan(&registered)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !registered {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	queryCheck := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, queryCheck, competitionID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	query := `SELECT is_active, created_at FROM competitor_competition_details WHERE competition_id = $1`
	var competitionDetails dto.CompetitorCompetitionDetails
	competitionDetails.CompetitionID = competitionID
	err = conn.QueryRow(context.Background(), query, competitionID).Scan(&competitionDetails.IsActive, &competitionDetails.CreatedAt)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	query = `SELECT competitor_id FROM competitor_competition_details WHERE competition_id = $1`
	rows, err := conn.Query(context.Background(), query, competitionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()

	var competitorIDs []int
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		competitorIDs = append(competitorIDs, id)
	}

	var competitor models.Competitor
	for _, competitorID := range competitorIDs {
		query = `SELECT id, full_name, birth_date, identity, bow, rank,region, federation, club FROM competitors WHERE id = $1`
		err = conn.QueryRow(context.Background(), query, competitorID).Scan(&competitor.ID, &competitor.FullName, &competitor.BirthDate, &competitor.Identity, &competitor.Bow, &competitor.Rank, &competitor.Region, &competitor.Federation, &competitor.Club)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		competitionDetails.Competitors = append(competitionDetails.Competitors, competitor)
	}
	tools.WriteJSON(w, http.StatusOK, competitionDetails)
}

func EditCompetitorStatus(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	competitorID, err := tools.ParseParamToInt(r, "competitor_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}

	if role == "user" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		if competitorID != userID {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	queryCheck := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, queryCheck, competitionID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	type RequestBody struct {
		IsActive bool `json:"is_active"`
	}
	var newStatus RequestBody
	err = json.NewDecoder(r.Body).Decode(&newStatus)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	var end bool
	queryCheck = `SELECT is_ended FROM competitions WHERE id = $1`
	err = conn.QueryRow(context.Background(), queryCheck, competitionID).Scan(&end)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if end {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	var registered bool
	queryCheck = `SELECT EXISTS (SELECT 1 FROM competitor_competition_details WHERE competition_id = $1 AND competitor_id = $2)`
	err = conn.QueryRow(context.Background(), queryCheck, competitionID, competitorID).Scan(&registered)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !registered {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	query := `SELECT ig.state
	FROM individual_groups ig
	JOIN competitor_group_details cgd ON ig.id = cgd.group_id
	WHERE ig.competition_id = $1
	AND cgd.competitor_id = $2`
	rows, err := conn.Query(context.Background(), query, competitionID, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()

	var status string
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if status != "created" {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	var competitorDetails dto.CompetitorCompetitionDetails
	competitorDetails.CompetitionID = competitionID
	query = "UPDATE competitor_competition_details SET is_active=$1 WHERE competition_id = $2 and competitor_id = $3 RETURNING is_active, created_at "
	err = conn.QueryRow(context.Background(), query, newStatus.IsActive, competitionID, competitorID).Scan(&competitorDetails.IsActive, &competitorDetails.CreatedAt)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var competitor models.Competitor
	query = `SELECT id, full_name, birth_date, identity, bow, rank,region, federation, club FROM competitors WHERE id = $1`
	err = conn.QueryRow(context.Background(), query, competitorID).Scan(&competitor.ID, &competitor.FullName, &competitor.BirthDate, &competitor.Identity, &competitor.Bow, &competitor.Rank, &competitor.Region, &competitor.Federation, &competitor.Club)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	competitorDetails.Competitors = append(competitorDetails.Competitors, competitor)
	tools.WriteJSON(w, http.StatusOK, competitorDetails)
}

func DeleteCompetitorCompetition(w http.ResponseWriter, r *http.Request) {
	competitionID, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	competitorID, err := tools.ParseParamToInt(r, "competitor_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	queryCheck := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, queryCheck, competitionID)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var registered bool
	queryCheck = `SELECT EXISTS (SELECT 1 FROM competitor_competition_details WHERE competition_id = $1 AND competitor_id = $2)`
	err = conn.QueryRow(context.Background(), queryCheck, competitionID, competitorID).Scan(&registered)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !registered {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	query := `SELECT ig.state
	FROM individual_groups ig
	JOIN competitor_group_details cgd ON ig.id = cgd.group_id
	WHERE ig.competition_id = $1
	AND cgd.competitor_id = $2`
	rows, err := conn.Query(context.Background(), query, competitionID, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()

	var status string
	for rows.Next() {
		err = rows.Scan(&status)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if status != "created" {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	query = `DELETE FROM competitor_competition_details WHERE competition_id = $1 AND competitor_id = $2`
	_, err = conn.Query(context.Background(), query, competitionID, competitorID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// TODO check propriety of enums
func CreateIndividualGroup(w http.ResponseWriter, r *http.Request) {
	competitionId, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	var is_ended bool
	queryCheck := `SELECT is_ended FROM competitions WHERE id = $1`
	err = conn.QueryRow(context.Background(), queryCheck, competitionId).Scan(&is_ended)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
			return
		}
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	if is_ended {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	var individualGroup models.IndividualGroup
	err = json.NewDecoder(r.Body).Decode(&individualGroup)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	var exist bool
	queryCheck = `SELECT EXISTS (SELECT 1 FROM individual_groups WHERE competition_id= $1 and bow = $2 and identity = $3)`
	err = conn.QueryRow(context.Background(), queryCheck, competitionId, individualGroup.Bow, individualGroup.Identity).Scan(&exist)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if exist {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "EXISTS"})
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer tx.Rollback(context.Background())

	individualGroup.CompetitionID = competitionId
	query := `INSERT INTO individual_groups (competition_id, bow, identity) VALUES ($1, $2, $3) RETURNING id, state`
	err = tx.QueryRow(context.Background(), query, individualGroup.CompetitionID, individualGroup.Bow, individualGroup.Identity).Scan(&individualGroup.ID, &individualGroup.State)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	query = `SELECT ccd.competitor_id 
		FROM competitor_competition_details ccd
        JOIN competitors c ON ccd.competitor_id = c.id
		WHERE ccd.competition_id = $1 
		AND ccd.is_active = $2
		AND c.identity = $3
		AND c.bow = $4`
	rows, err := tx.Query(context.Background(), query, competitionId, true, individualGroup.Identity, individualGroup.Bow)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()

	var competitorIDs []int
	var id int
	for rows.Next() {
		err = rows.Scan(&id)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		competitorIDs = append(competitorIDs, id)
	}

	query = `INSERT INTO competitor_group_details (group_id, competitor_id) VALUES ($1, $2)`
	for _, id = range competitorIDs {
		_, err = tx.Exec(context.Background(), query, individualGroup.ID, id)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	tools.WriteJSON(w, http.StatusCreated, individualGroup)
}

func GetIndividualGroupsFromCompetition(w http.ResponseWriter, r *http.Request) {
	competitionId, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	queryCheck := `SELECT id FROM competitions WHERE id = $1`
	exists, err := tools.ExistsInDB(context.Background(), conn, queryCheck, competitionId)
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	var rows pgx.Rows
	if role == "user" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		query := `SELECT ig.id, ig.competition_id, ig.bow, ig.identity, ig.state 
			FROM individual_groups ig 
			JOIN competitor_group_details cgd ON cgd.group_id = ig.id
			WHERE ig.competition_id = $1
			AND cgd.competitor_id = $2`

		rows, err = conn.Query(context.Background(), query, competitionId, userID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		defer rows.Close()
	} else {
		query := `SELECT id, competition_id, bow, identity, state FROM individual_groups WHERE competition_id = $1`
		rows, err = conn.Query(context.Background(), query, competitionId)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		defer rows.Close()
	}

	groups := make([]models.IndividualGroup, 0)
	var group models.IndividualGroup
	for rows.Next() {
		err = rows.Scan(&group.ID, &group.CompetitionID, &group.Bow, &group.Identity, &group.State)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		groups = append(groups, group)
	}
	tools.WriteJSON(w, http.StatusOK, groups)
}

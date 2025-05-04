package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"

	"app-server/internal/dto"
	"app-server/pkg/tools"

	"github.com/jackc/pgx/v5"

	"app-server/internal/models"
)

func GetIndividualGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()

	var individualGroup models.IndividualGroup
	err = conn.QueryRow(context.Background(), `SELECT * FROM individual_groups WHERE id = $1`, groupId).Scan(&individualGroup.ID, &individualGroup.CompetitionID, &individualGroup.Bow, &individualGroup.Identity, &individualGroup.State)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		} else {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		}
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
		err = conn.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM competitor_group_details 
			WHERE group_id = $1 AND competitor_id = $2)`, groupId, userID).Scan(&registered)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !registered {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	tools.WriteJSON(w, http.StatusOK, individualGroup)
}

func GetCompetitorsFromGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := tools.ParseParamToInt(r, "group_id")
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
	q := `
            SELECT EXISTS(
                SELECT 1 FROM individual_groups 
                WHERE id = $1
            )`
	var exists bool
	err = conn.QueryRow(context.Background(), q, groupId).Scan(&exists)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !exists {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
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
		checkQuery := `
            SELECT EXISTS(
                SELECT 1 FROM competitor_group_details 
                WHERE group_id = $1 AND competitor_id = $2
            )`
		err = conn.QueryRow(context.Background(), checkQuery, groupId, userID).Scan(&exists)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !exists {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	query := `
        SELECT c.id, c.full_name, c.birth_date, c.identity, c.bow, c.rank, c.region, c.federation, c.club
        FROM competitor_group_details cgd 
        JOIN competitors c ON cgd.competitor_id = c.id 
        WHERE cgd.group_id = $1`
	rows, err := conn.Query(context.Background(), query, groupId)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var cgd dto.CompetitorGroupDetail
	defer rows.Close()

	competitors := make([]models.Competitor, 0)
	for rows.Next() {
		var competitor models.Competitor
		if err = rows.Scan(&competitor.ID, &competitor.FullName, &competitor.BirthDate,
			&competitor.Identity, &competitor.Bow, &competitor.Rank, &competitor.Region,
			&competitor.Federation, &competitor.Club); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		competitors = append(competitors, competitor)
	}
	if err = rows.Err(); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	cgd.Competitors = competitors
	cgd.GroupID = groupId
	tools.WriteJSON(w, http.StatusOK, cgd)
}

func SyncIndividualGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()

	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	defer tx.Rollback(context.Background())

	var competitionID int
	err = tx.QueryRow(context.Background(), "SELECT competition_id FROM individual_groups WHERE id = $1", groupID).Scan(&competitionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "GROUP NOT FOUND"})
		return
	}

	var hasQualification bool
	err = tx.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM qualifications WHERE group_id = $1)", groupID).Scan(&hasQualification)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to check qualification: %v", err)})
		return
	}

	if hasQualification {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	_, err = tx.Exec(context.Background(), "DELETE FROM competitor_group_details WHERE group_id = $1", groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to delete old competitors: %v", err)})
		return
	}

	var competitorIDs []int
	rows, err := tx.Query(context.Background(), `
        INSERT INTO competitor_group_details (group_id, competitor_id)
        SELECT $1, competitor_id 
        FROM competitor_competition_details 
        WHERE competition_id = $2 AND is_active = true
        RETURNING competitor_id`,
		groupID, competitionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to insert new competitors: %v", err)})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to scan competitor id: %v", err)})
			return
		}
		competitorIDs = append(competitorIDs, id)
	}

	var result []map[string]interface{}
	if len(competitorIDs) > 0 {
		rows, err = tx.Query(context.Background(),
			`SELECT id, full_name, birth_date, identity, bow, rank, region, federation, club FROM competitors WHERE id = ANY($1)`,
			competitorIDs)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to fetch competitors details: %v", err)})
			return
		}
		defer rows.Close()

		for rows.Next() {
			var c models.Competitor

			if err := rows.Scan(&c.ID, &c.FullName, &c.BirthDate, &c.Identity, &c.Bow, &c.Rank, &c.Region, &c.Federation, &c.Club); err != nil {
				tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to scan competitor details: %v", err)})
				return
			}

			result = append(result, map[string]interface{}{"group_id": groupID, "competitor": c})
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to commit transaction: %v", err)})
		return
	}

	if result == nil {
		tools.WriteJSON(w, http.StatusOK, []interface{}{})
		return
	}

	tools.WriteJSON(w, http.StatusOK, result)
}

func GetQualification(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()

	resp, err := getQualification(conn.Conn(), groupID, r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	tools.WriteJSON(w, http.StatusOK, *resp)
}

func DeleteIndividualGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()

	ctx := r.Context()
	tx, err := conn.Begin(ctx)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, "failed to start transaction")
		return
	}
	defer tx.Rollback(ctx)

	var exists bool
	err = conn.QueryRow(ctx, `SELECT EXISTS(SELECT 1 FROM individual_groups WHERE id = $1)`, groupID).Scan(&exists)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, "failed to check if group exists")
		return
	}
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err := deleteAllGroupData(ctx, tx, groupID); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("failed to delete group data: %v", err))
		return
	}

	if err := tx.Commit(ctx); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, "failed to commit transaction")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func GetFinalGrid(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid group_id"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	exists, err := tools.ExistsInDB(r.Context(), conn,
		"SELECT EXISTS(SELECT 1 FROM individual_groups WHERE id = $1)", groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "database error"})
		return
	}
	if !exists {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "group not found"})
		return
	}

	var finalGrid models.FinalGrid
	finalGrid.GroupID = groupID

	if err := getQuarterfinals(r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
		return
	}
	if err := getSemifinals(r.Context(), groupID, &finalGrid.Semifinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinals"})
		return
	}
	if err := getFinals(r.Context(), groupID, &finalGrid.Final); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get finals"})
		return
	}

	tools.WriteJSON(w, http.StatusOK, finalGrid)
}

func StartQuarterfinal(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	tx, err := conn.Begin(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(r.Context())

	var groupState string
	var bowClass string
	err = tx.QueryRow(r.Context(), `SELECT state, bow FROM individual_groups WHERE id = $1 FOR UPDATE`, groupID).Scan(&groupState, &bowClass)
	if errors.Is(err, pgx.ErrNoRows) {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var finalGrid models.FinalGrid
	finalGrid.GroupID = groupID

	if groupState == "quarterfinal_start" || groupState == "semifinal_start" || groupState == "final_start" || groupState == "completed" {
		if err := getQuarterfinalsTx(tx, r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
			return
		}
		tools.WriteJSON(w, http.StatusOK, finalGrid)
		return
	}

	if groupState != "qualification_end" {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "qualification not completed"})
		return
	}

	qualifiers, err := getQualifiersTx(tx, r.Context(), groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get qualifiers"})
		return
	}

	if len(qualifiers) < 2 {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "not enough qualified competitors"})
		return
	}

	maxSeries := 5
	rangeType := "1-10"
	rangeSize := 3
	if bowClass != "block" {
		maxSeries = 3
		rangeType = "6-10"
	}

	_, err = tx.Exec(r.Context(), `UPDATE individual_groups SET state = 'quarterfinal_start' WHERE id = $1`, groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update group state"})
		return
	}

	if err := createQuarterfinalSparringsTx(tx, r.Context(), groupID, qualifiers, maxSeries, rangeSize, rangeType); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create quarterfinal sparrings"})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
		return
	}

	if err := getQuarterfinals(r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
		return
	}

	tools.WriteJSON(w, http.StatusCreated, finalGrid)
}

func StartSemifinal(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	tx, err := conn.Begin(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(r.Context())

	var groupState string
	var bowClass string
	err = tx.QueryRow(r.Context(), `SELECT state, bow FROM individual_groups WHERE id = $1 FOR UPDATE`, groupID).Scan(&groupState, &bowClass)
	if errors.Is(err, pgx.ErrNoRows) {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var finalGrid models.FinalGrid
	finalGrid.GroupID = groupID

	if groupState == "semifinal_start" || groupState == "final_start" || groupState == "completed" {
		if err := getQuarterfinalsTx(tx, r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
			return
		}

		if err := getSemifinalsTx(tx, r.Context(), groupID, &finalGrid.Semifinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinals"})
			return
		}
		tools.WriteJSON(w, http.StatusOK, finalGrid)
		return
	}

	if groupState != "quarterfinal_start" {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "quarterfinal not started"})
		return
	}

	if err := checkQuarterfinalsCompleted(tx, r.Context(), groupID); err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	winners, err := getQuarterfinalWinners(tx, r.Context(), groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinal winners"})
		return
	}

	maxSeries := 5
	rangeType := "1-10"
	rangeSize := 3
	if bowClass != "block" {
		maxSeries = 3
		rangeType = "6-10"
	}

	_, err = tx.Exec(r.Context(), `UPDATE individual_groups SET state = 'semifinal_start' WHERE id = $1`, groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update group state"})
		return
	}

	if err := createSemifinalSparringsTx(tx, r.Context(), groupID, winners, maxSeries, rangeSize, rangeType); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create semifinal sparrings"})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
		return
	}

	if err := getQuarterfinals(r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
		return
	}

	if err := getSemifinals(r.Context(), groupID, &finalGrid.Semifinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinals"})
		return
	}

	tools.WriteJSON(w, http.StatusCreated, finalGrid)
}

func StartFinal(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	tx, err := conn.Begin(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(r.Context())

	var groupState string
	var bowClass string
	err = tx.QueryRow(r.Context(), `SELECT state, bow FROM individual_groups WHERE id = $1 FOR UPDATE`, groupID).Scan(&groupState, &bowClass)
	if errors.Is(err, pgx.ErrNoRows) {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var finalGrid models.FinalGrid
	finalGrid.GroupID = groupID

	if groupState == "final_start" || groupState == "completed" {
		if err := getQuarterfinalsTx(tx, r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
			return
		}
		if err := getSemifinalsTx(tx, r.Context(), groupID, &finalGrid.Semifinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinals"})
			return
		}
		if err := getFinalsTx(tx, r.Context(), groupID, &finalGrid.Final); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get finals"})
			return
		}
		tools.WriteJSON(w, http.StatusOK, finalGrid)
		return
	}

	if groupState != "semifinal_start" {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "semifinal not started"})
		return
	}

	if err := checkSemifinalsCompleted(tx, r.Context(), groupID); err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	winners, losers, err := getSemifinalWinnersAndLosers(tx, r.Context(), groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinal winners and losers"})
		return
	}

	maxSeries := 5
	rangeType := "1-10"
	rangeSize := 3
	if bowClass != "block" {
		maxSeries = 3
		rangeType = "6-10"
	}

	_, err = tx.Exec(r.Context(), `UPDATE individual_groups SET state = 'final_start' WHERE id = $1`, groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update group state"})
		return
	}

	if err := createFinalSparringsTx(tx, r.Context(), groupID, winners, losers, maxSeries, rangeSize, rangeType); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to create final sparrings"})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
		return
	}

	if err := getQuarterfinals(r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
		return
	}
	if err := getSemifinals(r.Context(), groupID, &finalGrid.Semifinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinals"})
		return
	}
	if err := getFinals(r.Context(), groupID, &finalGrid.Final); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get finals"})
		return
	}

	tools.WriteJSON(w, http.StatusCreated, finalGrid)
}

func EndFinal(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	tx, err := conn.Begin(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(r.Context())

	var groupState string
	err = tx.QueryRow(r.Context(), `SELECT state FROM individual_groups WHERE id = $1 FOR UPDATE`, groupID).Scan(&groupState)
	if errors.Is(err, pgx.ErrNoRows) {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var finalGrid models.FinalGrid
	finalGrid.GroupID = groupID

	if groupState == "completed" {
		if err := getQuarterfinalsTx(tx, r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
			return
		}
		if err := getSemifinalsTx(tx, r.Context(), groupID, &finalGrid.Semifinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinals"})
			return
		}
		if err := getFinalsTx(tx, r.Context(), groupID, &finalGrid.Final); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get finals"})
			return
		}
		tools.WriteJSON(w, http.StatusOK, finalGrid)
		return
	}

	if groupState != "final_start" {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "final not started"})
		return
	}

	if err := checkFinalsCompleted(tx, r.Context(), groupID); err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": err.Error()})
		return
	}

	_, err = tx.Exec(r.Context(), `UPDATE individual_groups SET state = 'completed' WHERE id = $1`, groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update group state"})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
		return
	}

	if err := getQuarterfinals(r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
		return
	}
	if err := getSemifinals(r.Context(), groupID, &finalGrid.Semifinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinals"})
		return
	}
	if err := getFinals(r.Context(), groupID, &finalGrid.Final); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get finals"})
		return
	}

	tools.WriteJSON(w, http.StatusOK, finalGrid)
}

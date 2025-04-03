package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"app-server/pkg/tools"

	"github.com/jackc/pgx/v5"

	"app-server/internal/models"
)

func CreateIndividualGroup(w http.ResponseWriter, r *http.Request) {
	var individualGroup models.IndividualGroup
	err := json.NewDecoder(r.Body).Decode(&individualGroup)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	competitionId, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID ENDPOINT"})
		return
	}

	individualGroup.CompetitionID = competitionId
	_, err = conn.Exec(context.Background(), "INSERT INTO individual_groups (competition_id, bow, identity, state) VALUES ($1, $2, $3, $4)", individualGroup.CompetitionID, individualGroup.Bow, individualGroup.Identity, individualGroup.State)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to insert data: %v", err))
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetIndividualGroups(w http.ResponseWriter, r *http.Request) {
	groupId, err := tools.ParseParamToInt(r, "individual_group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	var individualGroup models.IndividualGroup
	err = conn.QueryRow(context.Background(), `SELECT * FROM individual_groups WHERE id = $1`, groupId).Scan(&individualGroup.ID, &individualGroup.CompetitionID, &individualGroup.Bow, &individualGroup.Identity, &individualGroup.State)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "individual group not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	role := r.Context().Value("role")
	if role != "admin" {
		id, ok := r.Context().Value("user_id").(int)
		if !ok {
			http.Error(w, "invalid user_id", http.StatusInternalServerError)
			return
		}

		var check bool
		err = conn.QueryRow(context.Background(), `SELECT EXISTS(
    SELECT 1 FROM competitor_group_details 
    WHERE group_id = $1 AND competitor_id = $2
)`, groupId, id).Scan(&check)
		if err != nil || !check {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "insufficient access rights"})
			return
		}
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(individualGroup); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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

func deleteShootOuts(ctx context.Context, tx pgx.Tx, groupID int) error {
	query := `
		DELETE FROM shoot_outs 
		WHERE place_id IN (
			SELECT id FROM sparring_places 
			WHERE range_group_id IN (
				SELECT id FROM range_groups 
				WHERE id IN (
					SELECT range_group_id FROM qualification_rounds 
					WHERE section_id IN (
						SELECT id FROM qualification_sections 
						WHERE group_id = $1
					)
				)
			)
		)`
	_, err := tx.Exec(ctx, query, groupID)
	return err
}

func deleteSparringPlaces(ctx context.Context, tx pgx.Tx, groupID int) error {
	query := `
		DELETE FROM sparring_places 
		WHERE range_group_id IN (
			SELECT id FROM range_groups 
			WHERE id IN (
				SELECT range_group_id FROM qualification_rounds 
				WHERE section_id IN (
					SELECT id FROM qualification_sections 
					WHERE group_id = $1
				)
			)
		)`
	_, err := tx.Exec(ctx, query, groupID)
	return err
}

func deleteQualificationRounds(ctx context.Context, tx pgx.Tx, groupID int) error {
	query := `
		DELETE FROM qualification_rounds 
		WHERE section_id IN (
			SELECT id FROM qualification_sections 
			WHERE group_id = $1
		)`
	_, err := tx.Exec(ctx, query, groupID)
	return err
}

func deleteQualificationSections(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `DELETE FROM qualification_sections WHERE group_id = $1`, groupID)
	return err
}

func deleteQualifications(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `DELETE FROM qualifications WHERE group_id = $1`, groupID)
	return err
}

func deleteCompetitorGroupDetails(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `DELETE FROM competitor_group_details WHERE group_id = $1`, groupID)
	return err
}

func deleteIndividualGroups(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `DELETE FROM individual_groups WHERE id = $1`, groupID)
	return err
}

func DeleteIndividualGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		http.Error(w, "invalid group_id", http.StatusBadRequest)
		return
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		http.Error(w, "Failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {

		}
	}(tx, context.Background())

	operations := []func(context.Context, pgx.Tx, int) error{
		deleteShootOuts,
		deleteSparringPlaces,
		deleteQualificationRounds,
		deleteQualificationSections,
		deleteQualifications,
		deleteCompetitorGroupDetails,
		deleteIndividualGroups,
	}

	for _, op := range operations {
		if err := op(r.Context(), tx, groupID); err != nil {
			http.Error(w, fmt.Sprintf("Failed to delete data: %v", err), http.StatusInternalServerError)
			return
		}
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "Failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	_, err = w.Write([]byte("Group and related data deleted successfully"))
	if err != nil {
		http.Error(w, "Failed to write response", http.StatusInternalServerError)
		return
	}
}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {
	groupId, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		http.Error(w, "invalid group_id", http.StatusBadRequest)
	}

	var req struct {
		_          int               `json:"group_id"`
		Competitor models.Competitor `json:"competitor"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "invalid request body", http.StatusBadRequest)
		return
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		http.Error(w, "failed to start transaction", http.StatusInternalServerError)
		return
	}
	defer func(tx pgx.Tx, ctx context.Context) {
		err := tx.Rollback(ctx)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	}(tx, context.Background())

	_, err = tx.Exec(context.Background(),
		`INSERT INTO competitors (id, full_name, birth_date, identity, bow, rank, region, federation, club)
		 VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
		 ON CONFLICT (id) DO UPDATE SET
			full_name = EXCLUDED.full_name,
			birth_date = EXCLUDED.birth_date,
			identity = EXCLUDED.identity,
			bow = EXCLUDED.bow,
			rank = EXCLUDED.rank,
			region = EXCLUDED.region,
			federation = EXCLUDED.federation,
			club = EXCLUDED.club`,
		req.Competitor.ID,
		req.Competitor.FullName,
		req.Competitor.BirthDate,
		req.Competitor.Identity,
		req.Competitor.Bow,
		req.Competitor.Rank,
		req.Competitor.Region,
		req.Competitor.Federation,
		req.Competitor.Club,
	)
	if err != nil {
		http.Error(w, "failed to upsert competitor", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(),
		`INSERT INTO competitor_group_details (group_id, competitor_id)
		 VALUES ($1, $2)
		 ON CONFLICT (group_id, competitor_id) DO NOTHING`,
		groupId,
		req.Competitor.ID,
	)
	if err != nil {
		fmt.Println(err)
		http.Error(w, "failed to add competitor to group", http.StatusInternalServerError)
		return
	}

	var competitionId int
	err = tx.QueryRow(context.Background(),
		`SELECT competition_id FROM individual_groups WHERE id = $1`,
		groupId,
	).Scan(&competitionId)
	if err != nil {
		http.Error(w, "failed to get competition_id", http.StatusInternalServerError)
		return
	}

	_, err = tx.Exec(context.Background(),
		`INSERT INTO competitor_competition_details (competition_id, competitor_id, is_active)
		 VALUES ($1, $2, true)
		 ON CONFLICT (competition_id, competitor_id) DO UPDATE SET
			is_active = EXCLUDED.is_active`,
		competitionId,
		req.Competitor.ID,
	)
	if err != nil {
		http.Error(w, "failed to update competition details", http.StatusInternalServerError)
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		http.Error(w, "failed to commit transaction", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
}

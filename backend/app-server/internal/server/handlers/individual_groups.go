package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"

	"app-server/internal/models"

	"github.com/gorilla/mux"
)

func ParseParam(r *http.Request, str string) (int, error) {
	vars := mux.Vars(r)
	result := vars[str]
	res, err := strconv.Atoi(result)
	return res, err
}

func CreateIndividualGroup(w http.ResponseWriter, r *http.Request) {
	var individualGroup models.IndividualGroup
	err := json.NewDecoder(r.Body).Decode(&individualGroup)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	competitionId, err := ParseParam(r, "competition_id")
	if err != nil {
		http.Error(w, "invalid competition_id", http.StatusBadRequest)
	}

	individualGroup.CompetitionID = competitionId
	_, err = conn.Exec(context.Background(), "INSERT INTO individual_groups (competition_id, bow, identity, state) VALUES ($1, $2, $3, $4)", individualGroup.CompetitionID, individualGroup.Bow, individualGroup.Identity, individualGroup.State)
	if err != nil {
		http.Error(w, "unable to insert data: %v\n", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
}

func GetIndividualGroups(w http.ResponseWriter, r *http.Request) {
	groupId, err := ParseParam(r, "group_id")
	if err != nil {
		http.Error(w, "invalid group_id", http.StatusBadRequest)
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

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(individualGroup); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
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
	groupID, err := ParseParam(r, "group_id")
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

func UpdateGroup(w http.ResponseWriter, r *http.Request) {}

func GetCompetitors(w http.ResponseWriter, r *http.Request) {
	groupId, err := ParseParam(r, "group_id")
	if err != nil {
		http.Error(w, "invalid group_id", http.StatusBadRequest)
	}

	var competitor models.Competitor
	err = conn.QueryRow(context.Background(), `SELECT c.id, c.full_name FROM competitor_group_details cgd JOIN competitors c ON cgd.competitor_id = c.id WHERE cgd.group_id = $1 `, groupId).Scan(&competitor.ID, &competitor.FullName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "competitor not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(competitor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

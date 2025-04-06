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

func SyncIndividualGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to begin transaction: %v", err))
		return
	}
	defer tx.Rollback(context.Background())

	var competitionID int
	err = tx.QueryRow(context.Background(),
		"SELECT competition_id FROM individual_groups WHERE id = $1", groupID).Scan(&competitionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "GROUP NOT FOUND"})
		return
	}

	_, err = tx.Exec(context.Background(),
		"DELETE FROM competitor_group_details WHERE group_id = $1", groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to delete old competitors: %v", err))
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
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to insert new competitors: %v", err))
		return
	}
	defer rows.Close()

	for rows.Next() {
		var id int
		if err := rows.Scan(&id); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to scan competitor id: %v", err))
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
			tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to fetch competitors details: %v", err))
			return
		}
		defer rows.Close()

		for rows.Next() {
			var c models.Competitor

			if err := rows.Scan(&c.ID, &c.FullName, &c.BirthDate, &c.Identity, &c.Bow, &c.Rank, &c.Region, &c.Federation, &c.Club); err != nil {
				tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to scan competitor details: %v", err))
				return
			}

			result = append(result, map[string]interface{}{"group_id": groupID, "competitor": c})
		}
	}

	var hasQualification bool
	err = tx.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM qualifications WHERE group_id = $1)", groupID).Scan(&hasQualification)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to check qualification: %v", err))
		return
	}

	if hasQualification {
		_, err = tx.Exec(context.Background(),
			"DELETE FROM qualification_sections WHERE group_id = $1", groupID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to delete old qualification sections: %v", err))
			return
		}

		_, err = tx.Exec(context.Background(), `
            INSERT INTO qualification_sections (group_id, competitor_id)
            SELECT $1, competitor_id 
            FROM competitor_competition_details 
            WHERE competition_id = $2 AND is_active = true`,
			groupID, competitionID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to insert new qualification sections: %v", err))
			return
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to commit transaction: %v", err))
		return
	}

	tools.WriteJSON(w, http.StatusOK, result)
}

func GetQualifications(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	var resp models.QualificationTable
	err = conn.QueryRow(context.Background(),
		`SELECT group_id, distance, round_count 
         FROM qualifications 
         WHERE group_id = $1`, groupID).Scan(
		&resp.GroupID, &resp.Distance, &resp.RoundCount)
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "QUALIFICATION NOT FOUND"})
		return
	}

	sections, err := getQualificationSections(groupID, r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to fetch sections: %v", err))
		return
	}

	resp.Sections = sections
	tools.WriteJSON(w, http.StatusOK, resp)
}

func getQualificationSections(groupID int, r *http.Request) ([]models.QualificationSectionForTable, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT qs.id, c.id, c.full_name, qs.place
         FROM qualification_sections qs
         JOIN competitors c ON qs.competitor_id = c.id
         WHERE qs.group_id = $1
         ORDER BY qs.id`, groupID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	role := r.Context().Value("role")
	id := r.Context().Value("user_id")
	ok := false
	var sections []models.QualificationSectionForTable
	for rows.Next() {
		var section models.QualificationSectionForTable
		if err := rows.Scan(&section.ID, &section.Competitor.ID, &section.Competitor.FullName, &section.Place); err != nil {
			return nil, err
		}

		if id == section.Competitor.ID {
			ok = true
		}
		rounds, totalScore, tensCount, ninesCount, err := getSectionRoundsStats(section.ID)
		if err != nil {
			return nil, err
		}

		section.Round = rounds
		section.Total = totalScore
		section.CountTen = tensCount
		section.CountNine = ninesCount

		sections = append(sections, section)
	}

	if role != "admin" && !ok {
		return nil, errors.New("no access rights")
	}

	return sections, nil
}

func getSectionRoundsStats(sectionID int) ([]models.RoundShrinked, int, int, int, error) {
	rows, err := conn.Query(context.Background(),
		`SELECT round_ordinal, is_active 
         FROM qualification_rounds 
         WHERE section_id = $1
         ORDER BY round_ordinal`, sectionID)
	if err != nil {
		return nil, 0, 0, 0, err
	}
	defer rows.Close()

	var rounds []models.RoundShrinked
	var totalScore, tensCount, ninesCount int

	for rows.Next() {
		var round models.RoundShrinked
		if err := rows.Scan(&round.RoundOrdinal, &round.IsActive); err != nil {
			return nil, 0, 0, 0, err
		}

		roundScore, tens, nines, err := getRoundStats(sectionID, round.RoundOrdinal)
		if err != nil {
			return nil, 0, 0, 0, err
		}

		round.TotalScore = roundScore
		rounds = append(rounds, round)

		totalScore += roundScore
		tensCount += tens
		ninesCount += nines
	}

	return rounds, totalScore, tensCount, ninesCount, nil
}

func getRoundStats(sectionID int, roundOrdinal int) (int, int, int, error) {
	var score, tens, nines int
	err := conn.QueryRow(context.Background(),
		`SELECT 
            COALESCE(SUM(
                CASE WHEN s.score ~ '^[0-9]+$' THEN s.score::integer ELSE 0 END
            ), 0),
            COALESCE(SUM(CASE WHEN s.score = '10' THEN 1 ELSE 0 END), 0),
            COALESCE(SUM(CASE WHEN s.score = '9' THEN 1 ELSE 0 END), 0)
         FROM shots s
         JOIN ranges r ON s.range_id = r.id
         JOIN qualification_rounds qr ON r.group_id = qr.range_group_id
         WHERE qr.section_id = $1 AND qr.round_ordinal = $2`,
		sectionID, roundOrdinal).Scan(&score, &tens, &nines)

	return score, tens, nines, err
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

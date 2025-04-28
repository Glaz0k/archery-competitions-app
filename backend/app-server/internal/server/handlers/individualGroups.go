package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"app-server/internal/dto"

	"app-server/pkg/tools"

	"github.com/jackc/pgx/v5"

	"app-server/internal/models"
)

func GetIndividualGroups(w http.ResponseWriter, r *http.Request) {
	groupId, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}

	var individualGroup models.IndividualGroup
	err = conn.QueryRow(context.Background(), `SELECT * FROM individual_groups WHERE id = $1`, groupId).Scan(&individualGroup.ID, &individualGroup.CompetitionID, &individualGroup.Bow, &individualGroup.Identity, &individualGroup.State)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "NOT FOUND", http.StatusNotFound)
		} else {
			http.Error(w, "DATABASE ERROR", http.StatusInternalServerError)
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

	var competitors []models.Competitor
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

	var hasQualification bool
	err = tx.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM qualifications WHERE group_id = $1)", groupID).Scan(&hasQualification)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to check qualification: %v", err)})
		return
	}

	if hasQualification {
		_, err = tx.Exec(context.Background(),
			"DELETE FROM qualification_rounds WHERE section_id IN (SELECT id FROM qualification_sections WHERE group_id = $1)",
			groupID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to delete old qualification rounds: %v", err)})
			return
		}

		_, err = tx.Exec(context.Background(),
			"DELETE FROM qualification_sections WHERE group_id = $1", groupID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to delete old qualification sections: %v", err)})
			return
		}

		_, err = tx.Exec(context.Background(), `
            INSERT INTO qualification_sections (group_id, competitor_id)
            SELECT $1, competitor_id 
            FROM competitor_competition_details 
            WHERE competition_id = $2 AND is_active = true`,
			groupID, competitionID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to insert new qualification sections: %v", err)})
			return
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to commit transaction: %v", err)})
		return
	}

	if result == nil {
		tools.WriteJSON(w, http.StatusOK, []interface{}{})
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
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to fetch sections: %v", err)})
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
		return nil, fmt.Errorf("query failed: %w", err)
	}

	role := r.Context().Value("role")
	userID := r.Context().Value("user_id")
	var sections []models.QualificationSectionForTable
	var hasAccess bool

	for rows.Next() {
		var section models.QualificationSectionForTable
		if err := rows.Scan(&section.ID, &section.Competitor.ID, &section.Competitor.FullName, &section.Place); err != nil {
			return nil, fmt.Errorf("scan failed: %w", err)
		}

		if userID == section.Competitor.ID {
			hasAccess = true
		}

		sections = append(sections, section)
	}

	if role != "admin" && !hasAccess {
		return nil, errors.New("no access rights")
	}
	rows.Close()

	for _, section := range sections {
		rounds, totalScore, tensCount, ninesCount, err := getSectionRoundsStats(section.ID)
		if err != nil {
			return nil, err
		}

		section.Round = rounds
		section.Total = totalScore
		section.CountTen = tensCount
		section.CountNine = ninesCount
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
		return nil, 0, 0, 0, fmt.Errorf("query rounds failed: %w", err)
	}

	var rounds []models.RoundShrinked
	var totalScore, tensCount, ninesCount int

	for rows.Next() {
		var round models.RoundShrinked
		if err := rows.Scan(&round.RoundOrdinal, &round.IsActive); err != nil {
			return nil, 0, 0, 0, fmt.Errorf("scan round failed: %w", err)
		}

		rounds = append(rounds, round)
	}
	rows.Close()

	for _, round := range rounds {
		roundScore, tens, nines, err := getRoundStats(sectionID, round.RoundOrdinal)
		if err != nil {
			return nil, 0, 0, 0, fmt.Errorf("get round stats failed: %w", err)
		}

		round.TotalScore = roundScore

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
            COALESCE(SUM(CASE WHEN s.score = 'X' THEN 10 WHEN s.score = 'M' THEN 0 ELSE CAST(s.score AS INTEGER) END), 0),
            COALESCE(SUM(CASE WHEN s.score = '10' THEN 1 WHEN s.score = 'X' THEN 1 ELSE 0 END), 0),
            COALESCE(SUM(CASE WHEN s.score = '9' THEN 1 ELSE 0 END), 0)
         FROM shots s
         JOIN ranges r ON s.range_id = r.id
         JOIN qualification_rounds qr ON r.group_id = qr.range_group_id
         WHERE qr.section_id = $1 AND qr.round_ordinal = $2`,
		sectionID, roundOrdinal).Scan(&score, &tens, &nines)

	if err != nil {
		return 0, 0, 0, fmt.Errorf("query stats failed: %w", err)
	}

	return score, tens, nines, nil
}

func deleteShots(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM shots 
        WHERE range_id IN (
            SELECT r.id 
            FROM ranges r
            JOIN range_groups rg ON r.group_id = rg.id
            JOIN sparring_places sp ON sp.range_group_id = rg.id
            JOIN sparrings s ON sp.id = s.top_place_id OR sp.id = s.bot_place_id
            JOIN (
                SELECT group_id, sparring1_id AS sparring_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring2_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring3_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring4_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring5_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring6_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_gold_id FROM finals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_bronze_id FROM finals WHERE group_id = $1
            ) stages ON s.id = stages.sparring_id
            WHERE stages.group_id = $1
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete shots: %w", err)
	}
	return nil
}

func deleteShootOuts(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM shoot_outs 
        WHERE place_id IN (
            SELECT sp.id 
            FROM sparring_places sp
            JOIN sparrings s ON sp.id = s.top_place_id OR sp.id = s.bot_place_id
            JOIN (
                SELECT group_id, sparring1_id AS sparring_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring2_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring3_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring4_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring5_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring6_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_gold_id FROM finals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_bronze_id FROM finals WHERE group_id = $1
            ) stages ON s.id = stages.sparring_id
            WHERE stages.group_id = $1
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete shoot outs: %w", err)
	}
	return nil
}

func deleteSparrings(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM sparrings 
        WHERE id IN (
            SELECT sparring_id 
            FROM (
                SELECT sparring1_id AS sparring_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT sparring2_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT sparring3_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT sparring4_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT sparring5_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT sparring6_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT sparring_gold_id FROM finals WHERE group_id = $1
                UNION
                SELECT sparring_bronze_id FROM finals WHERE group_id = $1
            ) stages
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete sparrings: %w", err)
	}
	return nil
}

func deleteSparringPlaces(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM sparring_places 
        WHERE id IN (
            SELECT sp.id 
            FROM sparring_places sp
            JOIN sparrings s ON sp.id = s.top_place_id OR sp.id = s.bot_place_id
            JOIN (
                SELECT group_id, sparring1_id AS sparring_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring2_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring3_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring4_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring5_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring6_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_gold_id FROM finals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_bronze_id FROM finals WHERE group_id = $1
            ) stages ON s.id = stages.sparring_id
            WHERE stages.group_id = $1
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete sparring places: %w", err)
	}
	return nil
}

func deleteRanges(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM ranges 
        WHERE group_id IN (
            SELECT rg.id 
            FROM range_groups rg
            JOIN sparring_places sp ON sp.range_group_id = rg.id
            JOIN sparrings s ON sp.id = s.top_place_id OR sp.id = s.bot_place_id
            JOIN (
                SELECT group_id, sparring1_id AS sparring_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring2_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring3_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring4_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring5_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring6_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_gold_id FROM finals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_bronze_id FROM finals WHERE group_id = $1
            ) stages ON s.id = stages.sparring_id
            WHERE stages.group_id = $1
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete ranges: %w", err)
	}
	return nil
}

func deleteRangeGroups(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM range_groups 
        WHERE id IN (
            SELECT rg.id 
            FROM range_groups rg
            JOIN sparring_places sp ON sp.range_group_id = rg.id
            JOIN sparrings s ON sp.id = s.top_place_id OR sp.id = s.bot_place_id
            JOIN (
                SELECT group_id, sparring1_id AS sparring_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring2_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring3_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring4_id FROM quarterfinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring5_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring6_id FROM semifinals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_gold_id FROM finals WHERE group_id = $1
                UNION
                SELECT group_id, sparring_bronze_id FROM finals WHERE group_id = $1
            ) stages ON s.id = stages.sparring_id
            WHERE stages.group_id = $1
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete range groups: %w", err)
	}
	return nil
}

func deleteQuarterfinals(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM quarterfinals
        WHERE group_id = $1`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete quarterfinals: %v", err)
	}
	return nil
}

func deleteSemifinals(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM semifinals
        WHERE group_id = $1`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete semifinals: %v", err)
	}
	return nil
}

func deleteFinals(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM finals
        WHERE group_id = $1`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete finals: %v", err)
	}
	return nil
}

func deleteQualificationRounds(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM qualification_rounds 
        WHERE section_id IN (
            SELECT id FROM qualification_sections 
            WHERE group_id = $1
        )`, groupID)
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

func deleteIndividualGroup(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `DELETE FROM individual_groups WHERE id = $1`, groupID)
	return err
}

func DeleteIndividualGroup(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, "invalid group_id")
		return
	}

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
		tools.WriteJSON(w, http.StatusNotFound, "group not found")
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

func deleteAllGroupData(ctx context.Context, tx pgx.Tx, groupID int) error {
	if err := deleteShootOuts(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete shoot outs: %v", err)
	}

	if err := deleteFinals(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete finals: %v", err)
	}

	if err := deleteSemifinals(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete semifinals: %v", err)
	}

	if err := deleteQuarterfinals(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete quarterfinals: %v", err)
	}

	if err := deleteSparrings(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete sparrings: %v", err)
	}

	if err := deleteSparringPlaces(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete sparring places: %v", err)
	}

	if err := deleteShots(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete shots: %v", err)
	}

	if err := deleteRanges(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete ranges: %v", err)
	}

	if err := deleteQualificationRounds(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete qualification rounds: %v", err)
	}

	if err := deleteRangeGroups(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete range groups: %v", err)
	}

	if err := deleteQualificationSections(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete qualification sections: %v", err)
	}

	if err := deleteQualifications(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete qualifications: %v", err)
	}

	if err := deleteCompetitorGroupDetails(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete competitor details: %v", err)
	}

	if err := deleteIndividualGroup(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete individual group: %v", err)
	}

	return nil
}

func GetFinalGrid(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "invalid group_id"})
		return
	}

	exists, err := tools.ExistsInDB(r.Context(), conn,
		"SELECT 1 FROM individual_groups WHERE id = $1", groupID)
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

func getStageSparrings(ctx context.Context, groupID int, stage string, result interface{}) error {
	var sparringCount int
	var sparringFields []string
	switch stage {
	case "quarterfinal":
		sparringCount = 4
		sparringFields = []string{"sparring1_id", "sparring2_id", "sparring3_id", "sparring4_id"}
	case "semifinal":
		sparringCount = 2
		sparringFields = []string{"sparring5_id", "sparring6_id"}
	case "final":
		sparringCount = 2
		sparringFields = []string{"sparring_gold_id", "sparring_bronze_id"}
	default:
		return fmt.Errorf("unsupported stage: %s", stage)
	}

	type sparringInfo struct {
		ID              int
		State           string
		TopPlaceID      sql.NullInt64
		TopCompID       sql.NullInt64
		TopFullName     sql.NullString
		TopRangeGroupID sql.NullInt64
		BotPlaceID      sql.NullInt64
		BotCompID       sql.NullInt64
		BotFullName     sql.NullString
		BotRangeGroupID sql.NullInt64
		SparringNum     int
	}

	caseClauses := make([]string, sparringCount)
	for i := 1; i <= sparringCount; i++ {
		caseClauses[i-1] = fmt.Sprintf("WHEN s.id = q.%s THEN %d", sparringFields[i-1], i)
	}
	caseStatement := fmt.Sprintf("CASE %s END AS sparring_num", strings.Join(caseClauses, " "))
	joinConditions := make([]string, sparringCount)
	for i := 1; i <= sparringCount; i++ {
		joinConditions[i-1] = fmt.Sprintf("q.%s = s.id", sparringFields[i-1])
	}
	joinCondition := strings.Join(joinConditions, " OR ")

	query := fmt.Sprintf(`
        SELECT s.id, s.state, 
               sp_top.id AS top_place_id, sp_top.competitor_id AS top_competitor_id, c_top.full_name AS top_full_name, sp_top.range_group_id AS top_range_group_id,
               sp_bot.id AS bot_place_id, sp_bot.competitor_id AS bot_competitor_id, c_bot.full_name AS bot_full_name, sp_bot.range_group_id AS bot_range_group_id,
               %s
        FROM %ss q
        LEFT JOIN sparrings s ON %s
        LEFT JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        LEFT JOIN competitors c_top ON sp_top.competitor_id = c_top.id
        LEFT JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        LEFT JOIN competitors c_bot ON sp_bot.competitor_id = c_bot.id
        WHERE q.group_id = $1
        ORDER BY sparring_num`, caseStatement, stage, joinCondition)

	rows, err := conn.Query(ctx, query, groupID)
	if err != nil {
		return fmt.Errorf("query %s: %w", stage, err)
	}
	defer rows.Close()

	var sparringInfos []sparringInfo
	for rows.Next() {
		var info sparringInfo
		err = rows.Scan(
			&info.ID,
			&info.State,
			&info.TopPlaceID,
			&info.TopCompID,
			&info.TopFullName,
			&info.TopRangeGroupID,
			&info.BotPlaceID,
			&info.BotCompID,
			&info.BotFullName,
			&info.BotRangeGroupID,
			&info.SparringNum,
		)
		if err != nil {
			return fmt.Errorf("failed to scan %s sparring: %w", stage, err)
		}
		sparringInfos = append(sparringInfos, info)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating %s sparring rows: %w", stage, err)
	}

	sparrings := make([]*models.Sparring, sparringCount)
	for i := range sparrings {
		sparrings[i] = &models.Sparring{}
	}

	rangeGroups := make(map[int64]*models.RangeGroup)
	for _, info := range sparringInfos {
		if info.TopRangeGroupID.Valid {
			rg := &models.RangeGroup{ID: int(info.TopRangeGroupID.Int64)}
			if err := getRangeGroup(rg); err != nil {
				return fmt.Errorf("failed to get top range group %d: %w", rg.ID, err)
			}
			rangeGroups[int64(rg.ID)] = rg
		}
		if info.BotRangeGroupID.Valid {
			rg := &models.RangeGroup{ID: int(info.BotRangeGroupID.Int64)}
			if err := getRangeGroup(rg); err != nil {
				return fmt.Errorf("failed to get bot range group %d: %w", rg.ID, err)
			}
			rangeGroups[int64(rg.ID)] = rg
		}
	}

	shootOuts := make(map[int64]models.ShootOuts)
	for _, info := range sparringInfos {
		if info.TopPlaceID.Valid {
			var so models.ShootOuts
			if getShotOut(&so, int(info.TopPlaceID.Int64)) {
				shootOuts[info.TopPlaceID.Int64] = so
			}
		}
		if info.BotPlaceID.Valid {
			var so models.ShootOuts
			if getShotOut(&so, int(info.BotPlaceID.Int64)) {
				shootOuts[info.BotPlaceID.Int64] = so
			}
		}
	}

	for _, info := range sparringInfos {
		var sparring models.Sparring
		var topPlace models.SparringPlace
		var botPlace models.SparringPlace
		var topComp models.CompetitorShrinked
		var botComp models.CompetitorShrinked

		sparring.ID = info.ID
		sparring.State = info.State

		if info.TopPlaceID.Valid {
			topPlace.ID = int(info.TopPlaceID.Int64)
			if info.TopCompID.Valid {
				topComp.ID = int(info.TopCompID.Int64)
				topComp.FullName = info.TopFullName.String
			}
			topPlace.Competitor = topComp

			if info.TopRangeGroupID.Valid {
				if rg, exists := rangeGroups[info.TopRangeGroupID.Int64]; exists {
					topPlace.RangeGroup = *rg
				}
			}

			if so, exists := shootOuts[int64(topPlace.ID)]; exists {
				topPlace.ShootOut = &so
			}

			topPlace.SparringScore = 0
		}

		if info.BotPlaceID.Valid {
			botPlace.ID = int(info.BotPlaceID.Int64)
			if info.BotCompID.Valid {
				botComp.ID = int(info.BotCompID.Int64)
				botComp.FullName = info.BotFullName.String
			}
			botPlace.Competitor = botComp

			if info.BotRangeGroupID.Valid {
				if rg, exists := rangeGroups[info.BotRangeGroupID.Int64]; exists {
					botPlace.RangeGroup = *rg
				}
			}

			if so, exists := shootOuts[int64(botPlace.ID)]; exists {
				botPlace.ShootOut = &so
			}

			botPlace.SparringScore = 0
		}

		if sparring.State == "ongoing" {
			topPlace.IsActive = true
			botPlace.IsActive = true
		} else {
			if info.TopPlaceID.Valid {
				topPlace.IsActive = sparring.State == "top_win"
			}
			if info.BotPlaceID.Valid {
				botPlace.IsActive = sparring.State == "bot_win"
			}
		}

		sparring.TopPlace = topPlace
		sparring.BotPlace = botPlace

		if info.SparringNum > 0 {
			sparrings[info.SparringNum-1] = &sparring
		}
	}

	switch stage {
	case "quarterfinal":
		qf := result.(*models.Quarterfinal)
		for i, sparring := range sparrings {
			if i == 0 {
				qf.Sparring1 = *sparring
			} else if i == 1 {
				qf.Sparring2 = *sparring
			} else if i == 2 {
				qf.Sparring3 = *sparring
			} else if i == 3 {
				qf.Sparring4 = *sparring
			}
		}
	case "semifinal":
		sf := result.(*models.Semifinal)
		for i, sparring := range sparrings {
			if i == 0 {
				sf.Sparring5 = *sparring
			} else if i == 1 {
				sf.Sparring6 = *sparring
			}
		}
	case "final":
		f := result.(*models.Final)
		for i, sparring := range sparrings {
			if i == 0 {
				f.SparringGold = *sparring
			} else if i == 1 {
				f.SparringBronze = *sparring
			}
		}
	}

	return nil
}

func getQuarterfinals(ctx context.Context, groupID int, qf *models.Quarterfinal) error {
	return getStageSparrings(ctx, groupID, "quarterfinal", qf)
}

func getSemifinals(ctx context.Context, groupID int, sf *models.Semifinal) error {
	return getStageSparrings(ctx, groupID, "semifinal", sf)
}

func getFinals(ctx context.Context, groupID int, f *models.Final) error {
	return getStageSparrings(ctx, groupID, "final", f)
}

func StartQuarterfinal(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

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

func getQualifiersTx(tx pgx.Tx, ctx context.Context, groupID int) ([]qualifier, error) {
	rows, err := tx.Query(ctx, `
        SELECT qs.competitor_id, qs.place
        FROM qualification_sections qs
        WHERE qs.group_id = $1 AND qs.place IS NOT NULL
        ORDER BY qs.place`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get qualifiers: %w", err)
	}
	defer rows.Close()

	var qualifiers []qualifier
	for rows.Next() {
		var q qualifier
		if err := rows.Scan(&q.CompetitorID, &q.Place); err != nil {
			return nil, fmt.Errorf("failed to scan qualifier: %w", err)
		}
		qualifiers = append(qualifiers, q)
	}
	return qualifiers, nil
}

func createRangeGroupTx(tx pgx.Tx, ctx context.Context, maxSeries, rangeSize int, rangeType string) (int64, error) {
	var rangeGroupID int64
	err := tx.QueryRow(ctx, `INSERT INTO range_groups (ranges_max_count, range_size, type) VALUES ($1, $2, $3) RETURNING id`,
		maxSeries, rangeSize, rangeType).Scan(&rangeGroupID)
	if err != nil {
		return 0, fmt.Errorf("failed to create range group: %w", err)
	}
	return rangeGroupID, nil
}

func createSparringPlaceWithRanges(tx pgx.Tx, ctx context.Context, competitorID, maxSeries, rangeSize int, rangeType string) (int64, error) {
	rangeGroupID, err := createRangeGroupTx(tx, ctx, maxSeries, rangeSize, rangeType)
	if err != nil {
		return 0, fmt.Errorf("failed to create range group: %w", err)
	}

	if rangeGroupID == 0 {
		return 0, fmt.Errorf("invalid range group ID")
	}

	placeID, err := createSparringPlaceTx(tx, ctx, sql.NullInt64{Int64: rangeGroupID, Valid: true}, competitorID)
	if err != nil {
		return 0, fmt.Errorf("failed to create sparring place: %w", err)
	}

	rangeIDs, err := createRangesTx(tx, ctx, rangeGroupID, maxSeries)
	if err != nil {
		return 0, fmt.Errorf("failed to create ranges: %w", err)
	}

	if err := createShotsTx(tx, ctx, rangeIDs, rangeSize); err != nil {
		return 0, fmt.Errorf("failed to create shots: %w", err)
	}

	var savedRangeGroupID sql.NullInt64
	err = tx.QueryRow(ctx, `SELECT range_group_id FROM sparring_places WHERE id = $1`, placeID).Scan(&savedRangeGroupID)
	if err != nil {
		return 0, fmt.Errorf("failed to verify sparring place: %w", err)
	}
	if !savedRangeGroupID.Valid || savedRangeGroupID.Int64 != rangeGroupID {
		return 0, fmt.Errorf("range_group_id mismatch in sparring place")
	}

	return placeID, nil
}

func createQuarterfinalSparringsTx(tx pgx.Tx, ctx context.Context, groupID int, qualifiers []qualifier, maxSeries, rangeSize int, rangeType string) error {
	var quarterfinalID int64
	err := tx.QueryRow(ctx, `INSERT INTO quarterfinals (group_id) VALUES ($1) RETURNING group_id`, groupID).Scan(&quarterfinalID)
	if err != nil {
		return fmt.Errorf("failed to create quarterfinal record: %w", err)
	}

	sparringPairs := [][]int{{1, 8}, {5, 4}, {3, 6}, {7, 2}}

	for i, pair := range sparringPairs {
		topPlace := findQualifier(qualifiers, pair[0])
		botPlace := findQualifier(qualifiers, pair[1])

		var topPlaceID, botPlaceID int64
		var state string

		if topPlace == nil && botPlace == nil {
			continue
		}

		if topPlace != nil {
			if botPlace == nil {
				var nullRangeGroupID sql.NullInt64
				topPlaceID, err = createSparringPlaceTx(tx, ctx, nullRangeGroupID, topPlace.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create top place for pair %v: %w", pair, err)
				}
				state = "top_win"
			} else {
				topPlaceID, err = createSparringPlaceWithRanges(tx, ctx, topPlace.CompetitorID, maxSeries, rangeSize, rangeType)
				if err != nil {
					return fmt.Errorf("failed to create top place with ranges for pair %v: %w", pair, err)
				}
			}
		}

		if botPlace != nil {
			if topPlace == nil {
				var nullRangeGroupID sql.NullInt64
				botPlaceID, err = createSparringPlaceTx(tx, ctx, nullRangeGroupID, botPlace.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create bot place for pair %v: %w", pair, err)
				}
				state = "bot_win"
			} else {
				botPlaceID, err = createSparringPlaceWithRanges(tx, ctx, botPlace.CompetitorID, maxSeries, rangeSize, rangeType)
				if err != nil {
					return fmt.Errorf("failed to create bot place with ranges for pair %v: %w", pair, err)
				}
				state = "ongoing"
			}
		}

		sparringID, err := createSparringTx(tx, ctx, topPlaceID, botPlaceID, state)
		if err != nil {
			return fmt.Errorf("failed to create sparring for pair %v: %w", pair, err)
		}

		if err := linkSparringToQuarterfinalTx(tx, ctx, quarterfinalID, i+1, sparringID); err != nil {
			return fmt.Errorf("failed to link sparring for pair %v: %w", pair, err)
		}
	}
	return nil
}

func createSparringPlaceTx(tx pgx.Tx, ctx context.Context, rangeGroupID sql.NullInt64, competitorID int) (int64, error) {
	var placeID int64
	err := tx.QueryRow(ctx, `INSERT INTO sparring_places (competitor_id, range_group_id) VALUES ($1, $2) RETURNING id`,
		competitorID, rangeGroupID).Scan(&placeID)
	if err != nil {
		return 0, fmt.Errorf("failed to create sparring place: %w", err)
	}
	return placeID, nil
}

func createSparringTx(tx pgx.Tx, ctx context.Context, topPlaceID, botPlaceID int64, state string) (int64, error) {
	var sparringID int64
	err := tx.QueryRow(ctx, `INSERT INTO sparrings (top_place_id, bot_place_id, state) VALUES ($1, $2, $3) RETURNING id`,
		sql.NullInt64{Int64: topPlaceID, Valid: topPlaceID != 0}, sql.NullInt64{Int64: botPlaceID, Valid: botPlaceID != 0}, state).
		Scan(&sparringID)
	if err != nil {
		return 0, fmt.Errorf("failed to create sparring: %w", err)
	}
	return sparringID, nil
}

func linkSparringToQuarterfinalTx(tx pgx.Tx, ctx context.Context, quarterfinalID int64, sparringNum int, sparringID int64) error {
	updateField := fmt.Sprintf("sparring%d_id", sparringNum)
	_, err := tx.Exec(ctx, fmt.Sprintf(`UPDATE quarterfinals SET %s = $1 WHERE group_id = $2`, updateField), sparringID, quarterfinalID)
	if err != nil {
		return fmt.Errorf("failed to link sparring: %w", err)
	}
	return nil
}

func findQualifier(qualifiers []qualifier, place int) *qualifier {
	for _, q := range qualifiers {
		if q.Place == place {
			return &q
		}
	}
	return nil
}

func createRangesTx(tx pgx.Tx, ctx context.Context, rangeGroupID int64, maxSeries int) ([]int64, error) {
	var rangeIDs []int64
	for rangeOrdinal := 1; rangeOrdinal <= maxSeries; rangeOrdinal++ {
		var rangeID int64
		err := tx.QueryRow(ctx, `INSERT INTO ranges (group_id, range_ordinal, is_active) VALUES ($1, $2, $3) RETURNING id`,
			rangeGroupID, rangeOrdinal, false).Scan(&rangeID)
		if err != nil {
			return nil, fmt.Errorf("failed to create range %d for group_id %d: %w", rangeOrdinal, rangeGroupID, err)
		}
		rangeIDs = append(rangeIDs, rangeID)
	}
	return rangeIDs, nil
}

func createShotsTx(tx pgx.Tx, ctx context.Context, rangeIDs []int64, rangeSize int) error {
	for _, rangeID := range rangeIDs {
		for shotOrdinal := 1; shotOrdinal <= rangeSize; shotOrdinal++ {
			_, err := tx.Exec(ctx, `
                INSERT INTO shots (range_id, shot_ordinal, score)
                VALUES ($1, $2, $3)`,
				rangeID, shotOrdinal, "0")
			if err != nil {
				return fmt.Errorf("failed to create shot for range_id %d, shot_ordinal %d: %w", rangeID, shotOrdinal, err)
			}
		}
	}
	return nil
}

func getQuarterfinalsTx(tx pgx.Tx, ctx context.Context, groupID int, qf *models.Quarterfinal) error {
	type sparringInfo struct {
		ID              int
		State           string
		TopPlaceID      sql.NullInt64
		TopCompID       sql.NullInt64
		TopFullName     sql.NullString
		TopRangeGroupID sql.NullInt64
		BotPlaceID      sql.NullInt64
		BotCompID       sql.NullInt64
		BotFullName     sql.NullString
		BotRangeGroupID sql.NullInt64
		SparringNum     int
	}

	rows, err := tx.Query(ctx, `
        SELECT s.id, s.state, 
               sp_top.id AS top_place_id, sp_top.competitor_id AS top_competitor_id, c_top.full_name AS top_full_name, sp_top.range_group_id AS top_range_group_id,
               sp_bot.id AS bot_place_id, sp_bot.competitor_id AS bot_competitor_id, c_bot.full_name AS bot_full_name, sp_bot.range_group_id AS bot_range_group_id,
               CASE 
                   WHEN s.id = q.sparring1_id THEN 1
                   WHEN s.id = q.sparring2_id THEN 2
                   WHEN s.id = q.sparring3_id THEN 3
                   WHEN s.id = q.sparring4_id THEN 4
               END AS sparring_num
        FROM quarterfinals q
        LEFT JOIN sparrings s ON q.sparring1_id = s.id OR q.sparring2_id = s.id OR q.sparring3_id = s.id OR q.sparring4_id = s.id
        LEFT JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        LEFT JOIN competitors c_top ON sp_top.competitor_id = c_top.id
        LEFT JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        LEFT JOIN competitors c_bot ON sp_bot.competitor_id = c_bot.id
        WHERE q.group_id = $1
        ORDER BY sparring_num`, groupID)
	if err != nil {
		return fmt.Errorf("query quarterfinals: %w", err)
	}
	defer rows.Close()

	var sparringInfos []sparringInfo
	for rows.Next() {
		var info sparringInfo
		err = rows.Scan(
			&info.ID,
			&info.State,
			&info.TopPlaceID,
			&info.TopCompID,
			&info.TopFullName,
			&info.TopRangeGroupID,
			&info.BotPlaceID,
			&info.BotCompID,
			&info.BotFullName,
			&info.BotRangeGroupID,
			&info.SparringNum,
		)
		if err != nil {
			return fmt.Errorf("failed to scan sparring: %w", err)
		}
		sparringInfos = append(sparringInfos, info)
	}

	if err = rows.Err(); err != nil {
		return fmt.Errorf("error iterating sparring rows: %w", err)
	}

	sparrings := make([]*models.Sparring, 4)
	for i := range sparrings {
		sparrings[i] = &models.Sparring{}
	}

	rangeGroups := make(map[int64]*models.RangeGroup)
	for _, info := range sparringInfos {
		if info.TopRangeGroupID.Valid {
			rg := &models.RangeGroup{ID: int(info.TopRangeGroupID.Int64)}
			if err := getRangeGroup(rg); err != nil {
				return fmt.Errorf("failed to get top range group %d: %w", rg.ID, err)
			}
			rangeGroups[info.TopRangeGroupID.Int64] = rg
		}
		if info.BotRangeGroupID.Valid {
			rg := &models.RangeGroup{ID: int(info.BotRangeGroupID.Int64)}
			if err := getRangeGroup(rg); err != nil {
				return fmt.Errorf("failed to get bot range group %d: %w", rg.ID, err)
			}
			rangeGroups[info.BotRangeGroupID.Int64] = rg
		}
	}

	shootOuts := make(map[int64]models.ShootOuts)
	for _, info := range sparringInfos {
		if info.TopPlaceID.Valid {
			var so models.ShootOuts
			if getShotOut(&so, int(info.TopPlaceID.Int64)) {
				shootOuts[info.TopPlaceID.Int64] = so
			}
		}
		if info.BotPlaceID.Valid {
			var so models.ShootOuts
			if getShotOut(&so, int(info.BotPlaceID.Int64)) {
				shootOuts[info.BotPlaceID.Int64] = so
			}
		}
	}

	for _, info := range sparringInfos {
		var sparring models.Sparring
		var topPlace models.SparringPlace
		var botPlace models.SparringPlace
		var topComp models.CompetitorShrinked
		var botComp models.CompetitorShrinked

		sparring.ID = info.ID
		sparring.State = info.State

		if info.TopPlaceID.Valid {
			topPlace.ID = int(info.TopPlaceID.Int64)
			if info.TopCompID.Valid {
				topComp.ID = int(info.TopCompID.Int64)
				topComp.FullName = info.TopFullName.String
			}
			topPlace.Competitor = topComp

			if info.TopRangeGroupID.Valid {
				if rg, exists := rangeGroups[info.TopRangeGroupID.Int64]; exists {
					topPlace.RangeGroup = *rg
				}
			} else {
				topPlace.RangeGroup = models.RangeGroup{}
			}

			if so, exists := shootOuts[int64(topPlace.ID)]; exists {
				topPlace.ShootOut = &so
			}

			topPlace.SparringScore = 0
		}

		if info.BotPlaceID.Valid {
			botPlace.ID = int(info.BotPlaceID.Int64)
			if info.BotCompID.Valid {
				botComp.ID = int(info.BotCompID.Int64)
				botComp.FullName = info.BotFullName.String
			}
			botPlace.Competitor = botComp

			if info.BotRangeGroupID.Valid {
				if rg, exists := rangeGroups[info.BotRangeGroupID.Int64]; exists {
					botPlace.RangeGroup = *rg
				}
			} else {
				botPlace.RangeGroup = models.RangeGroup{}
			}

			if so, exists := shootOuts[int64(botPlace.ID)]; exists {
				botPlace.ShootOut = &so
			}

			botPlace.SparringScore = 0
		}

		if sparring.State == "ongoing" {
			topPlace.IsActive = true
			botPlace.IsActive = true
		} else {
			if info.TopPlaceID.Valid {
				topPlace.IsActive = sparring.State == "top_win"
			}
			if info.BotPlaceID.Valid {
				botPlace.IsActive = sparring.State == "bot_win"
			}
		}

		sparring.TopPlace = topPlace
		sparring.BotPlace = botPlace

		if info.SparringNum > 0 {
			sparrings[info.SparringNum-1] = &sparring
		}
	}

	qf.Sparring1 = *sparrings[0]
	qf.Sparring2 = *sparrings[1]
	qf.Sparring3 = *sparrings[2]
	qf.Sparring4 = *sparrings[3]

	return nil
}

type qualifier struct {
	CompetitorID int
	Place        int
}

func StartSemifinal(w http.ResponseWriter, r *http.Request) {
	groupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

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

	if groupState == "semifinal_start" || groupState == "final_start" || groupState == "completed" {
		if err := getSemifinalsTx(tx, r.Context(), groupID, &finalGrid.Semifinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	var bowClass string
	err = tx.QueryRow(r.Context(), `SELECT bow FROM individual_groups WHERE id = $1`, groupID).Scan(&bowClass)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get bow class"})
		return
	}

	rangeGroupID, err := createRangeGroupTx(tx, r.Context(), 1, 1, bowClass)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_, err = tx.Exec(r.Context(), `UPDATE individual_groups SET state = 'semifinal_start' WHERE id = $1`, groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update group state"})
		return
	}

	if err := createSemifinalSparringsTx(tx, r.Context(), groupID, winners, rangeGroupID); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
		return
	}

	tx, err = conn.Begin(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to begin transaction"})
		return
	}
	defer tx.Rollback(r.Context())

	if err := getSemifinalsTx(tx, r.Context(), groupID, &finalGrid.Semifinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
		return
	}

	tools.WriteJSON(w, http.StatusCreated, finalGrid)
}

func checkQuarterfinalsCompleted(tx pgx.Tx, ctx context.Context, groupID int) error {
	var ongoingCount int
	err := tx.QueryRow(ctx, `
       SELECT COUNT(s.id)
       FROM quarterfinals q
       JOIN sparrings s ON q.sparring1_id = s.id OR q.sparring2_id = s.id OR q.sparring3_id = s.id OR q.sparring4_id = s.id
       WHERE q.group_id = $1 AND s.state = 'ongoing'`, groupID).Scan(&ongoingCount)
	if err != nil {
		return fmt.Errorf("failed to check quarterfinals completion: %w", err)
	}
	if ongoingCount > 0 {
		return fmt.Errorf("not all quarterfinal sparrings are completed")
	}
	return nil
}

func getQuarterfinalWinners(tx pgx.Tx, ctx context.Context, groupID int) ([]qualifier, error) {
	rows, err := tx.Query(ctx, `
       SELECT
           CASE
               WHEN s.state = 'top_win' THEN sp_top.competitor_id
               WHEN s.state = 'bot_win' THEN sp_bot.competitor_id
           END AS winner_id,
           CASE
               WHEN s.id = q.sparring1_id THEN 1
               WHEN s.id = q.sparring2_id THEN 2
               WHEN s.id = q.sparring3_id THEN 3
               WHEN s.id = q.sparring4_id THEN 4
           END AS sparring_num
       FROM quarterfinals q
       JOIN sparrings s ON q.sparring1_id = s.id OR q.sparring2_id = s.id OR q.sparring3_id = s.id OR q.sparring4_id = s.id
       LEFT JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
       LEFT JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
       WHERE q.group_id = $1 AND s.state IN ('top_win', 'bot_win')
       ORDER BY sparring_num`, groupID)
	if err != nil {
		return nil, fmt.Errorf("failed to get quarterfinal winners: %w", err)
	}
	defer rows.Close()

	var winners []qualifier
	for rows.Next() {
		var q qualifier
		var winnerID sql.NullInt64
		if err := rows.Scan(&winnerID, &q.Place); err != nil {
			return nil, fmt.Errorf("failed to scan winner: %w", err)
		}
		if winnerID.Valid {
			q.CompetitorID = int(winnerID.Int64)
			winners = append(winners, q)
		}
	}
	return winners, rows.Err()
}

func createSemifinalSparringsTx(tx pgx.Tx, ctx context.Context, groupID int, winners []qualifier, rangeGroupID int64) error {
	var semifinalID int64
	err := tx.QueryRow(ctx, `INSERT INTO semifinals (group_id) VALUES ($1) RETURNING group_id`, groupID).Scan(&semifinalID)
	if err != nil {
		return fmt.Errorf("failed to create semifinal record: %w", err)
	}

	sparringPairs := [][]int{{1, 2}, {3, 4}}

	for i, pair := range sparringPairs {
		topWinner := findQualifier(winners, pair[0])
		botWinner := findQualifier(winners, pair[1])

		var topPlaceID, botPlaceID int64
		var topWinnerFlag, botWinnerFlag bool

		if topWinner != nil {
			if botWinner == nil {
				topPlaceID, err = createSparringPlaceTx(tx, ctx, sql.NullInt64{Int64: rangeGroupID}, topWinner.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create top place: %w", err)
				}
				topWinnerFlag = true
			} else {
				topPlaceID, err = createSparringPlaceTx(tx, ctx, sql.NullInt64{Int64: rangeGroupID}, topWinner.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create top place: %w", err)
				}
			}
		}

		if botWinner != nil {
			if topWinner == nil {
				botPlaceID, err = createSparringPlaceTx(tx, ctx, sql.NullInt64{Int64: rangeGroupID}, botWinner.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create bot place: %w", err)
				}
				botWinnerFlag = true
			} else {
				botPlaceID, err = createSparringPlaceTx(tx, ctx, sql.NullInt64{Int64: rangeGroupID}, botWinner.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create bot place: %w", err)
				}
			}
		}

		if topWinner != nil || botWinner != nil {
			state := "ongoing"
			if topWinnerFlag {
				state = "top_win"
			} else if botWinnerFlag {
				state = "bot_win"
			}

			sparringID, err := createSparringTx(tx, ctx, topPlaceID, botPlaceID, state)
			if err != nil {
				return fmt.Errorf("failed to create sparring: %w", err)
			}

			if err := linkSparringToSemifinalTx(tx, ctx, semifinalID, i+1, sparringID); err != nil {
				return fmt.Errorf("failed to link sparring: %w", err)
			}
		}
	}
	return nil
}

func linkSparringToSemifinalTx(tx pgx.Tx, ctx context.Context, semifinalID int64, sparringNum int, sparringID int64) error {
	updateField := fmt.Sprintf("sparring%d_id", sparringNum)
	_, err := tx.Exec(ctx, fmt.Sprintf(`UPDATE semifinals SET %s = $1 WHERE group_id = $2`, updateField), sparringID, semifinalID)
	if err != nil {
		return fmt.Errorf("failed to link sparring: %w", err)
	}
	return nil
}

func getSemifinalsTx(tx pgx.Tx, ctx context.Context, groupID int, sf *models.Semifinal) error {
	rows, err := tx.Query(ctx, `SELECT s.id, s.state, sp_top.id, sp_top.competitor_id, c_top.full_name, sp_bot.id, sp_bot.competitor_id, c_bot.full_name
      FROM semifinals sf
      JOIN sparrings s ON sf.sparring1_id = s.id OR sf.sparring2_id = s.id
      LEFT JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
      LEFT JOIN competitors c_top ON sp_top.competitor_id = c_top.id
      LEFT JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
      LEFT JOIN competitors c_bot ON sp_bot.competitor_id = c_bot.id
      WHERE sf.group_id = $1
      ORDER BY
          CASE
              WHEN s.id = sf.sparring1_id THEN 1
              WHEN s.id = sf.sparring2_id THEN 2
          END`, groupID)
	if err != nil {
		return fmt.Errorf("query semifinals: %w", err)
	}
	defer rows.Close()

	//sparrings, err := getSparringFromRows(rows)
	if err != nil {
		return fmt.Errorf("get sparrings: %w", err)
	}

	//sf.Sparring5 = *sparrings[0]
	//sf.Sparring6 = *sparrings[1]

	return nil
}

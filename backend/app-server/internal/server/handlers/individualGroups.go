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
            SELECT id FROM ranges
            WHERE group_id IN (
                SELECT id FROM range_groups 
                WHERE id IN (
                    SELECT range_group_id FROM qualification_rounds
                    WHERE section_id IN (
                        SELECT id FROM qualification_sections 
                        WHERE group_id = $1
                    )
                )
            )
        )`, groupID)
	return err
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

func deleteShootOuts(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM shoot_outs 
        WHERE place_id IN (
            SELECT sp.id FROM sparring_places sp
            JOIN range_groups rg ON sp.range_group_id = rg.id
            JOIN qualification_rounds qr ON rg.id = qr.range_group_id
            JOIN qualification_sections qs ON qr.section_id = qs.id
            WHERE qs.group_id = $1
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete shoot outs: %v", err)
	}
	return nil
}

func deleteSparrings(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM sparrings 
        WHERE top_place_id IN (
            SELECT id FROM sparring_places 
            WHERE range_group_id IN (
                SELECT rg.id FROM range_groups rg
                JOIN qualification_rounds qr ON rg.id = qr.range_group_id
                JOIN qualification_sections qs ON qr.section_id = qs.id
                WHERE qs.group_id = $1
            )
        ) OR bot_place_id IN (
            SELECT id FROM sparring_places 
            WHERE range_group_id IN (
                SELECT rg.id FROM range_groups rg
                JOIN qualification_rounds qr ON rg.id = qr.range_group_id
                JOIN qualification_sections qs ON qr.section_id = qs.id
                WHERE qs.group_id = $1
            )
        )
        OR EXISTS (
            SELECT 1 FROM quarterfinals qf
            WHERE qf.group_id = $1
            AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
        )
        OR EXISTS (
            SELECT 1 FROM semifinals sf
            WHERE sf.group_id = $1
            AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
        )
        OR EXISTS (
            SELECT 1 FROM finals f
            WHERE f.group_id = $1
            AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete sparrings: %v", err)
	}
	return nil
}

func deleteSparringPlaces(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM sparring_places 
        WHERE range_group_id IN (
            SELECT rg.id FROM range_groups rg
            JOIN qualification_rounds qr ON rg.id = qr.range_group_id
            JOIN qualification_sections qs ON qr.section_id = qs.id
            WHERE qs.group_id = $1
        )
        OR id IN (
            SELECT top_place_id FROM sparrings
            WHERE EXISTS (
                SELECT 1 FROM quarterfinals qf
                WHERE qf.group_id = $1
                AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
            )
            OR EXISTS (
                SELECT 1 FROM semifinals sf
                WHERE sf.group_id = $1
                AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
            )
            OR EXISTS (
                SELECT 1 FROM finals f
                WHERE f.group_id = $1
                AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
            )
        )
        OR id IN (
            SELECT bot_place_id FROM sparrings
            WHERE EXISTS (
                SELECT 1 FROM quarterfinals qf
                WHERE qf.group_id = $1
                AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
            )
            OR EXISTS (
                SELECT 1 FROM semifinals sf
                WHERE sf.group_id = $1
                AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
            )
            OR EXISTS (
                SELECT 1 FROM finals f
                WHERE f.group_id = $1
                AND (sparrings.top_place_id = id OR sparrings.bot_place_id = id)
            )
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete sparring places: %v", err)
	}
	return nil
}

func deleteRanges(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM ranges 
        WHERE group_id IN (
            SELECT id FROM range_groups 
            WHERE id IN (
                SELECT range_group_id FROM qualification_rounds
                WHERE section_id IN (
                    SELECT id FROM qualification_sections 
                    WHERE group_id = $1
                )
            )
        )`, groupID)
	return err
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

func deleteRangeGroups(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM range_groups 
        WHERE id IN (
            SELECT range_group_id 
            FROM qualification_rounds 
            WHERE section_id IN (
                SELECT id 
                FROM qualification_sections 
                WHERE group_id = $1
            )
        )`, groupID)
	if err != nil {
		return fmt.Errorf("failed to delete range groups: %v", err)
	}
	return nil
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

	err = getQuarterfinals(r.Context(), groupID, &finalGrid.Quarterfinal)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get quarterfinals"})
		return
	}

	err = getSemifinals(r.Context(), groupID, &finalGrid.Semifinal)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get semifinals"})
		return
	}

	err = getFinals(r.Context(), groupID, &finalGrid.Final)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get finals"})
		return
	}

	tools.WriteJSON(w, http.StatusOK, finalGrid)
}

func getQuarterfinals(ctx context.Context, groupID int, qf *models.Quarterfinal) error {
	rows, err := conn.Query(ctx, `
        SELECT s.id, s.state, sp_top.id, sp_top.competitor_id, c.full_name, sp_bot.id, sp_bot.competitor_id, c2.full_name
        FROM quarterfinals q
        JOIN sparrings s ON q.sparring1_id = s.id OR q.sparring2_id = s.id OR q.sparring3_id = s.id OR q.sparring4_id = s.id
        JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        JOIN competitors c ON sp_top.competitor_id = c.id
        JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        JOIN competitors c2 ON sp_bot.competitor_id = c2.id
        WHERE q.group_id = $1`, groupID)
	if err != nil {
		return fmt.Errorf("failed to get quarterfinals: %v", err)
	}
	defer rows.Close()

	tmp, err := getSparringFromRows(rows)
	if err != nil {
		return fmt.Errorf("failed to get sparring from rows: %v", err)
	}
	qf.Sparring1 = *tmp[0]
	qf.Sparring2 = *tmp[1]
	qf.Sparring3 = *tmp[2]
	qf.Sparring4 = *tmp[3]

	return nil
}

func getSemifinals(ctx context.Context, groupID int, sf *models.Semifinal) error {
	rows, err := conn.Query(ctx, `
        SELECT s.id, s.state, sp_top.id, sp_top.competitor_id, c.full_name, sp_bot.id, sp_bot.competitor_id, c2.full_name
        FROM semifinals sm
        JOIN sparrings s ON sm.sparring5_id = s.id OR sm.sparring6_id = s.id
        JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        JOIN competitors c ON sp_top.competitor_id = c.id
        JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        JOIN competitors c2 ON sp_bot.competitor_id = c2.id
        WHERE sm.group_id = $1`, groupID)
	if err != nil {
		return fmt.Errorf("query semifinals: %w", err)
	}
	defer rows.Close()

	tmp, err := getSparringFromRows(rows)
	if err != nil {
		return fmt.Errorf("rows error: %w", err)
	}
	sf.Sparring5 = *tmp[0]
	sf.Sparring6 = *tmp[1]

	return nil
}

func getFinals(ctx context.Context, groupID int, f *models.Final) error {
	rows, err := conn.Query(ctx, `
        SELECT s.id, s.state, sp_top.id, sp_top.competitor_id, c.full_name, sp_bot.id, sp_bot.competitor_id, c2.full_name
        FROM finals fl
        JOIN sparrings s ON fl.sparring_gold_id = s.id OR fl.sparring_bronze_id = s.id
        JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        JOIN competitors c ON sp_top.competitor_id = c.id
        JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        JOIN competitors c2 ON sp_bot.competitor_id = c2.id
        WHERE fl.group_id = $1`, groupID)
	if err != nil {
		return fmt.Errorf("query finals: %w", err)
	}
	defer rows.Close()

	tmp, err := getSparringFromRows(rows)
	if err != nil {
		return fmt.Errorf("rows error: %w", err)
	}
	f.SparringGold = *tmp[0]
	f.SparringBronze = *tmp[1]

	return nil
}

type qualifier struct {
	CompetitorID int
	Place        int
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

	if groupState == "quarterfinal_start" || groupState == "semifinal_start" || groupState == "final_start" || groupState == "completed" {
		if err := getQuarterfinalsTx(tx, r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if len(qualifiers) < 2 {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "not enough qualified competitors"})
		return
	}

	var bowClass string
	err = tx.QueryRow(r.Context(), `SELECT bow FROM individual_groups WHERE id = $1`, groupID).Scan(&bowClass)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to get bow class"})
		return
	}

	rangeGroupID, err := createRangeGroupTx(tx, r.Context(), bowClass)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_, err = tx.Exec(r.Context(), `UPDATE individual_groups SET state = 'quarterfinal_start' WHERE id = $1`, groupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to update group state"})
		return
	}

	if err := createQuarterfinalSparringsTx(tx, r.Context(), groupID, qualifiers, rangeGroupID, bowClass); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	if err := tx.Commit(r.Context()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "failed to commit transaction"})
		return
	}

	if err := getQuarterfinals(r.Context(), groupID, &finalGrid.Quarterfinal); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
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

func createRangeGroupTx(tx pgx.Tx, ctx context.Context, bowClass string) (int64, error) {
	maxSeries := 5
	rangeType := "1-10"
	if bowClass != "block" {
		maxSeries = 3
		rangeType = "6-10"
	}
	rangeSize := 3

	var rangeGroupID int64
	err := tx.QueryRow(ctx, `INSERT INTO range_groups (ranges_max_count, range_size, type) VALUES ($1, $2, $3) RETURNING id`,
		maxSeries, rangeSize, rangeType).Scan(&rangeGroupID)
	if err != nil {
		return 0, fmt.Errorf("failed to create range group: %w", err)
	}
	return rangeGroupID, nil
}

func createQuarterfinalSparringsTx(tx pgx.Tx, ctx context.Context, groupID int, qualifiers []qualifier, rangeGroupID int64, bowClass string) error {
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
		var topWinner, botWinner bool

		if topPlace != nil {
			if botPlace == nil {
				topPlaceID, err = createSparringPlaceTx(tx, ctx, rangeGroupID, topPlace.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create top place: %w", err)
				}
				topWinner = true
			} else {
				topPlaceID, err = createSparringPlaceTx(tx, ctx, rangeGroupID, topPlace.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create top place: %w", err)
				}
			}
		}

		if botPlace != nil {
			if topPlace == nil {
				botPlaceID, err = createSparringPlaceTx(tx, ctx, rangeGroupID, botPlace.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create bot place: %w", err)
				}
				botWinner = true
			} else {
				botPlaceID, err = createSparringPlaceTx(tx, ctx, rangeGroupID, botPlace.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create bot place: %w", err)
				}
			}
		}

		if topPlace != nil || botPlace != nil {
			state := "ongoing"
			if topWinner {
				state = "top_win"
			} else if botWinner {
				state = "bot_win"
			}

			sparringID, err := createSparringTx(tx, ctx, topPlaceID, botPlaceID, state)
			if err != nil {
				return fmt.Errorf("failed to create sparring: %w", err)
			}

			if err := linkSparringToQuarterfinalTx(tx, ctx, quarterfinalID, i+1, sparringID); err != nil {
				return fmt.Errorf("failed to link sparring: %w", err)
			}
		}
	}
	return nil
}

func createSparringPlaceTx(tx pgx.Tx, ctx context.Context, rangeGroupID int64, competitorID int) (int64, error) {
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

func getQuarterfinalsTx(tx pgx.Tx, ctx context.Context, groupID int, qf *models.Quarterfinal) error {
	rows, err := tx.Query(ctx, `
        SELECT s.id, s.state, sp_top.id, sp_top.competitor_id, c_top.full_name, sp_bot.id, sp_bot.competitor_id, c_bot.full_name
        FROM quarterfinals q
        JOIN sparrings s ON q.sparring1_id = s.id OR q.sparring2_id = s.id OR q.sparring3_id = s.id OR q.sparring4_id = s.id
        LEFT JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        LEFT JOIN competitors c_top ON sp_top.competitor_id = c_top.id
        LEFT JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        LEFT JOIN competitors c_bot ON sp_bot.competitor_id = c_bot.id
        WHERE q.group_id = $1
        ORDER BY 
            CASE 
                WHEN s.id = q.sparring1_id THEN 1
                WHEN s.id = q.sparring2_id THEN 2
                WHEN s.id = q.sparring3_id THEN 3
                WHEN s.id = q.sparring4_id THEN 4
            END`, groupID)
	if err != nil {
		return fmt.Errorf("query quarterfinals: %w", err)
	}
	defer rows.Close()

	sparrings, err := getSparringFromRows(rows)
	if err != nil {
		return fmt.Errorf("get sparrings: %w", err)
	}

	if len(sparrings) > 0 {
		qf.Sparring1 = *sparrings[0]
	}
	if len(sparrings) > 1 {
		qf.Sparring2 = *sparrings[1]
	}
	if len(sparrings) > 2 {
		qf.Sparring3 = *sparrings[2]
	}
	if len(sparrings) > 3 {
		qf.Sparring4 = *sparrings[3]
	}

	return nil
}

func getSparringFromRows(rows pgx.Rows) ([]*models.Sparring, error) {
	var sparrings []*models.Sparring
	for rows.Next() {
		var sparring models.Sparring
		var topPlace models.SparringPlace
		var botPlace models.SparringPlace
		var topComp models.CompetitorShrinked
		var botComp models.CompetitorShrinked

		var topPlaceID, botPlaceID sql.NullInt64
		var topCompID, botCompID sql.NullInt64
		var topFullName, botFullName sql.NullString

		err := rows.Scan(
			&sparring.ID,
			&sparring.State,
			&topPlaceID,
			&topCompID,
			&topFullName,
			&botPlaceID,
			&botCompID,
			&botFullName,
		)
		if err != nil {
			return nil, fmt.Errorf("failed to scan sparring: %w", err)
		}

		if topPlaceID.Valid {
			topPlace.ID = int(topPlaceID.Int64)
		}
		if topCompID.Valid {
			topComp.ID = int(topCompID.Int64)
		}
		topComp.FullName = topFullName.String
		topPlace.Competitor = topComp

		if botPlaceID.Valid {
			botPlace.ID = int(botPlaceID.Int64)
		}
		if botCompID.Valid {
			botComp.ID = int(botCompID.Int64)
		}
		botComp.FullName = botFullName.String
		botPlace.Competitor = botComp

		sparring.TopPlace = topPlace
		sparring.BotPlace = botPlace
		sparrings = append(sparrings, &sparring)
	}
	return sparrings, rows.Err()
}

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
            SELECT r.id FROM ranges r
            JOIN range_groups rg ON r.group_id = rg.id
            JOIN qualification_rounds qr ON rg.id = qr.range_group_id
            JOIN qualification_sections qs ON qr.section_id = qs.id
            WHERE qs.group_id = $1
        )`, groupID)
	return err
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
	return err
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
        )`, groupID)
	return err
}

func deleteSparringPlaces(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM sparring_places 
        WHERE range_group_id IN (
            SELECT rg.id FROM range_groups rg
            JOIN qualification_rounds qr ON rg.id = qr.range_group_id
            JOIN qualification_sections qs ON qr.section_id = qs.id
            WHERE qs.group_id = $1
        )`, groupID)
	return err
}

func deleteRanges(ctx context.Context, tx pgx.Tx, groupID int) error {
	_, err := tx.Exec(ctx, `
        DELETE FROM ranges 
        WHERE group_id IN (
            SELECT rg.id FROM range_groups rg
            JOIN qualification_rounds qr ON rg.id = qr.range_group_id
            JOIN qualification_sections qs ON qr.section_id = qs.id
            WHERE qs.group_id = $1
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
            SELECT range_group_id FROM qualification_rounds 
            WHERE section_id IN (
                SELECT id FROM qualification_sections 
                WHERE group_id = $1
            )
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

	exists, err := tools.ExistsInDB(context.Background(), conn, `SELECT * FROM individual_groups WHERE group_id = $1`, groupID)
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

	if err := deleteSparringPlaces(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete sparring places: %v", err)
	}

	if err := deleteSparrings(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete sparrings: %v", err)
	}

	if err := deleteShots(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete shots: %v", err)
	}

	if err := deleteRanges(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete ranges: %v", err)
	}

	if err := deleteRangeGroups(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete range groups: %v", err)
	}

	if err := deleteQualificationRounds(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete qualification rounds: %v", err)
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

func getSparringFromRows(ctx context.Context, rows pgx.Rows) ([]models.Sparring, error) {
	var tmp []models.Sparring
	for rows.Next() {
		var sparring models.Sparring
		var topPlace, botPlace models.SparringPlace
		var topComp, botComp models.CompetitorShrinked

		err := rows.Scan(
			&sparring.ID, &sparring.State,
			&topPlace.ID, &topComp.ID, &topComp.FullName,
			&botPlace.ID, &botComp.ID, &botComp.FullName,
		)
		if err != nil {
			return nil, err
		}

		topPlace.Competitor = topComp
		botPlace.Competitor = botComp

		err = loadRangeGroupAndShots(ctx, &topPlace)
		if err != nil {
			return nil, err
		}

		err = loadRangeGroupAndShots(ctx, &botPlace)
		if err != nil {
			return nil, err
		}

		sparring.TopPlace = topPlace
		sparring.BotPlace = botPlace

		tmp = append(tmp, sparring)
	}
	return tmp, rows.Err()
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

	tmp, err := getSparringFromRows(ctx, rows)
	if err != nil {
		return fmt.Errorf("failed to get sparring from rows: %v", err)
	}
	qf.Sparring1 = tmp[0]
	qf.Sparring2 = tmp[1]
	qf.Sparring3 = tmp[2]
	qf.Sparring4 = tmp[3]

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

	tmp, err := getSparringFromRows(ctx, rows)
	if err != nil {
		return fmt.Errorf("rows error: %w", err)
	}
	sf.Sparring5 = tmp[0]
	sf.Sparring6 = tmp[1]

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

	tmp, err := getSparringFromRows(ctx, rows)
	if err != nil {
		return fmt.Errorf("rows error: %w", err)
	}
	f.SparringGold = tmp[0]
	f.SparringBronze = tmp[1]

	return nil
}

func loadRangeGroupAndShots(ctx context.Context, place *models.SparringPlace) error {
	var rg models.RangeGroup
	err := conn.QueryRow(ctx, `
        SELECT rg.id, rg.ranges_max_count, rg.range_size
        FROM range_groups rg
        JOIN sparring_places sp ON sp.range_group_id = rg.id
        WHERE sp.id = $1`, place.ID).Scan(&rg.ID, &rg.RangesMaxCount, &rg.RangeSize)
	if err != nil {
		return fmt.Errorf("query range group: %w", err)
	}

	rows, err := conn.Query(ctx, `
        SELECT r.id, r.range_ordinal, r.is_active
        FROM ranges r
        WHERE r.group_id = $1
        ORDER BY r.range_ordinal`, rg.ID)
	if err != nil {
		return fmt.Errorf("query ranges: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var r models.Range
		err := rows.Scan(&r.ID, &r.RangeNumber, &r.IsOngoing)
		if err != nil {
			return fmt.Errorf("scan range: %w", err)
		}

		if err := loadShotsForRange(ctx, &r); err != nil {
			return fmt.Errorf("load shots: %w", err)
		}

		rg.Ranges = append(rg.Ranges, r)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("rows error: %w", err)
	}

	place.RangeGroup = rg
	return nil
}

func loadShotsForRange(ctx context.Context, r *models.Range) error {
	rows, err := conn.Query(ctx, `
        SELECT s.shot_ordinal, s.score
        FROM shots s
        WHERE s.range_id = $1
        ORDER BY s.shot_ordinal`, r.ID)
	if err != nil {
		return fmt.Errorf("query shots: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var s models.Shot
		err := rows.Scan(&s.ShotNumber, &s.Score)
		if err != nil {
			return fmt.Errorf("scan shot: %w", err)
		}
		r.Shots = append(r.Shots, s)
	}

	return rows.Err()
}

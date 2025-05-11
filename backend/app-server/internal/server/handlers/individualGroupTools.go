package handlers

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strings"

	"app-server/internal/models"

	"github.com/jackc/pgx/v5"
)

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

func createSemifinalSparringsTx(tx pgx.Tx, ctx context.Context, groupID int, winners []qualifier, maxSeries, rangeSize int, rangeType string) error {
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
		var state string

		if topWinner == nil && botWinner == nil {
			continue
		}

		if topWinner != nil {
			if botWinner == nil {
				var nullRangeGroupID sql.NullInt64
				topPlaceID, err = createSparringPlaceTx(tx, ctx, nullRangeGroupID, topWinner.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create top place for pair %v: %w", pair, err)
				}
				state = "top_win"
			} else {
				topPlaceID, err = createSparringPlaceWithRanges(tx, ctx, topWinner.CompetitorID, maxSeries, rangeSize, rangeType)
				if err != nil {
					return fmt.Errorf("failed to create top place with ranges for pair %v: %w", pair, err)
				}
			}
		}

		if botWinner != nil {
			if topWinner == nil {
				var nullRangeGroupID sql.NullInt64
				botPlaceID, err = createSparringPlaceTx(tx, ctx, nullRangeGroupID, botWinner.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create bot place for pair %v: %w", pair, err)
				}
				state = "bot_win"
			} else {
				botPlaceID, err = createSparringPlaceWithRanges(tx, ctx, botWinner.CompetitorID, maxSeries, rangeSize, rangeType)
				if err != nil {
					return fmt.Errorf("failed to create bot place with ranges for pair %v: %w", pair, err)
				}
				state = "ongoing"
			}
		}

		if topPlaceID != 0 || botPlaceID != 0 {
			sparringID, err := createSparringTx(tx, ctx, topPlaceID, botPlaceID, state)
			if err != nil {
				return fmt.Errorf("failed to create sparring for pair %v: %w", pair, err)
			}

			if err := linkSparringToSemifinalTx(tx, ctx, semifinalID, i+5, sparringID); err != nil {
				return fmt.Errorf("failed to link sparring for pair %v: %w", pair, err)
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
	rows, err := tx.Query(ctx, `
        SELECT s.id, s.state, 
               sp_top.id, sp_top.competitor_id, sp_top.range_group_id, c_top.full_name,
               sp_bot.id, sp_bot.competitor_id, sp_bot.range_group_id, c_bot.full_name
        FROM semifinals sf
        JOIN sparrings s ON sf.sparring5_id = s.id OR sf.sparring6_id = s.id
        LEFT JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        LEFT JOIN competitors c_top ON sp_top.competitor_id = c_top.id
        LEFT JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        LEFT JOIN competitors c_bot ON sp_bot.competitor_id = c_bot.id
        WHERE sf.group_id = $1
        ORDER BY
            CASE
                WHEN s.id = sf.sparring5_id THEN 1
                WHEN s.id = sf.sparring6_id THEN 2
            END`, groupID)
	if err != nil {
		return fmt.Errorf("failed to query semifinals: %w", err)
	}
	defer rows.Close()

	var sparrings []*models.Sparring
	for rows.Next() {
		var s models.Sparring
		var topPlaceID, topCompetitorID, botPlaceID, botCompetitorID sql.NullInt64
		var topRangeGroupID, botRangeGroupID sql.NullInt64
		var topFullName, botFullName sql.NullString

		if err := rows.Scan(&s.ID, &s.State,
			&topPlaceID, &topCompetitorID, &topRangeGroupID, &topFullName,
			&botPlaceID, &botCompetitorID, &botRangeGroupID, &botFullName); err != nil {
			return fmt.Errorf("failed to scan sparring: %w", err)
		}

		s.TopPlace = &models.SparringPlace{ID: int(topPlaceID.Int64), IsActive: true}
		if topCompetitorID.Valid {
			s.TopPlace.Competitor = models.CompetitorShrinked{
				ID:       int(topCompetitorID.Int64),
				FullName: topFullName.String,
			}
		}
		if topRangeGroupID.Valid {
			s.TopPlace.RangeGroup.ID = int(topRangeGroupID.Int64)
			if err := getRangeGroup(&s.TopPlace.RangeGroup); err != nil {
				return fmt.Errorf("failed to get top range group: %w", err)
			}
		}

		s.BotPlace = &models.SparringPlace{ID: int(botPlaceID.Int64)}
		if botCompetitorID.Valid {
			s.BotPlace.Competitor = models.CompetitorShrinked{
				ID:       int(botCompetitorID.Int64),
				FullName: botFullName.String,
			}
			s.BotPlace.IsActive = true
		}
		if botRangeGroupID.Valid {
			s.BotPlace.RangeGroup.ID = int(botRangeGroupID.Int64)
			if err := getRangeGroup(&s.BotPlace.RangeGroup); err != nil {
				return fmt.Errorf("failed to get bot range group: %w", err)
			}
		}

		sparrings = append(sparrings, &s)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating sparring rows: %w", err)
	}

	if sf.Sparring6.BotPlace.Competitor.FullName == "" {
		sf = nil
		return nil
	}

	if len(sparrings) > 0 {
		sf.Sparring5 = *sparrings[0]
	}
	if len(sparrings) > 1 {
		sf.Sparring6 = *sparrings[1]
	}

	return nil
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
				rangeID, shotOrdinal, nil)
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

		sparring.TopPlace = &topPlace
		sparring.BotPlace = &botPlace

		if topPlace.Competitor.FullName == "" {
			sparring.TopPlace = nil
		}
		if botPlace.Competitor.FullName == "" {
			sparring.BotPlace = nil
		}

		if sparring.BotPlace == nil || sparring.TopPlace == nil {
			if sparring.TopPlace == nil {
				sparring.BotPlace.IsActive = false
				sparring.BotPlace.RangeGroup.Type = "1-10"
				rgs := make([]models.Range, 0)
				sparring.BotPlace.RangeGroup.Ranges = rgs
			} else {
				sparring.TopPlace.IsActive = false
				sparring.TopPlace.RangeGroup.Type = "1-10"
				rgs := make([]models.Range, 0)
				sparring.TopPlace.RangeGroup.Ranges = rgs
			}
		}
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
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed Acquire")
	}
	defer conn.Release()

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

		sparring.TopPlace = &topPlace
		sparring.BotPlace = &botPlace

		if topPlace.Competitor.FullName == "" {
			sparring.TopPlace = nil
		}
		if botPlace.Competitor.FullName == "" {
			sparring.BotPlace = nil
		}

		if info.SparringNum > 0 {
			sparrings[info.SparringNum-1] = &sparring
		}
	}

	switch stage {
	case "quarterfinal":
		qf := result.(*models.Quarterfinal)
		for i, sparring := range sparrings {
			if sparring.BotPlace == nil || sparring.TopPlace == nil {
				if sparring.TopPlace == nil {
					sparring.BotPlace.IsActive = false
					sparring.BotPlace.RangeGroup.Type = "1-10"
					rgs := make([]models.Range, 0)
					sparring.BotPlace.RangeGroup.Ranges = rgs
				} else {
					sparring.TopPlace.IsActive = false
					sparring.TopPlace.RangeGroup.Type = "1-10"
					rgs := make([]models.Range, 0)
					sparring.TopPlace.RangeGroup.Ranges = rgs
				}
			}
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

func getQualificationSections(groupID int, r *http.Request) ([]models.QualificationSectionForTable, error) {
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		return nil, fmt.Errorf("acquire failed: %w", err)
	}
	defer conn.Release()
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

	for i, section := range sections {
		rounds, totalScore, tensCount, ninesCount, err := getSectionRoundsStats(section.ID)
		if err != nil {
			return nil, err
		}

		section.Round = rounds
		section.Total = totalScore
		section.CountTen = tensCount
		section.CountNine = ninesCount
		sections[i] = section
	}

	return sections, nil
}

func getSectionRoundsStats(sectionID int) ([]models.RoundShrinked, int, int, int, error) {
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return nil, 0, 0, 0, fmt.Errorf("acquire failed: %w", err)
	}
	defer conn.Release()
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

	for i, round := range rounds {
		roundScore, tens, nines, err := getRoundStats(sectionID, round.RoundOrdinal)
		if err != nil {
			return nil, 0, 0, 0, fmt.Errorf("get round stats failed: %w", err)
		}

		round.TotalScore = roundScore

		totalScore += roundScore
		tensCount += tens
		ninesCount += nines
		rounds[i] = round
	}

	return rounds, totalScore, tensCount, ninesCount, nil
}

func getRoundStats(sectionID int, roundOrdinal int) (int, int, int, error) {
	var score, tens, nines int
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return 0, 0, 0, fmt.Errorf("acquire failed: %w", err)
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(),
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

	if err := deleteShots(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete shots: %v", err)
	}

	if err := deleteRanges(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete ranges: %v", err)
	}

	if err := deleteSparrings(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete sparrings: %v", err)
	}

	if err := deleteSparringPlaces(ctx, tx, groupID); err != nil {
		return fmt.Errorf("failed to delete sparring places: %v", err)
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

func checkSemifinalsCompleted(tx pgx.Tx, ctx context.Context, groupID int) error {
	var ongoingCount int
	err := tx.QueryRow(ctx, `
        SELECT COUNT(s.id)
        FROM semifinals sf
        JOIN sparrings s ON sf.sparring5_id = s.id OR sf.sparring6_id = s.id
        WHERE sf.group_id = $1 AND s.state = 'ongoing'`, groupID).Scan(&ongoingCount)
	if err != nil {
		return fmt.Errorf("failed to check semifinals completion: %w", err)
	}
	if ongoingCount > 0 {
		return fmt.Errorf("not all semifinal sparrings are completed")
	}
	return nil
}

func getSemifinalWinnersAndLosers(tx pgx.Tx, ctx context.Context, groupID int) ([]qualifier, []qualifier, error) {
	rows, err := tx.Query(ctx, `
        SELECT
            CASE
                WHEN s.state = 'top_win' THEN sp_top.competitor_id
                WHEN s.state = 'bot_win' THEN sp_bot.competitor_id
            END AS winner_id,
            CASE
                WHEN s.state = 'top_win' THEN sp_bot.competitor_id
                WHEN s.state = 'bot_win' THEN sp_top.competitor_id
            END AS loser_id,
            CASE
                WHEN s.id = sf.sparring5_id THEN 5
                WHEN s.id = sf.sparring6_id THEN 6
            END AS sparring_num
        FROM semifinals sf
        JOIN sparrings s ON sf.sparring5_id = s.id OR sf.sparring6_id = s.id
        LEFT JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        LEFT JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        WHERE sf.group_id = $1 AND s.state IN ('top_win', 'bot_win')
        ORDER BY sparring_num`, groupID)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to get semifinal winners and losers: %w", err)
	}
	defer rows.Close()

	var winners, losers []qualifier
	for rows.Next() {
		var w, l qualifier
		var winnerID, loserID sql.NullInt64
		if err := rows.Scan(&winnerID, &loserID, &w.Place); err != nil {
			return nil, nil, fmt.Errorf("failed to scan winner and loser: %w", err)
		}
		l.Place = w.Place
		if winnerID.Valid {
			w.CompetitorID = int(winnerID.Int64)
			winners = append(winners, w)
		}
		if loserID.Valid {
			l.CompetitorID = int(loserID.Int64)
			losers = append(losers, l)
		}
	}
	return winners, losers, rows.Err()
}

func createFinalSparringsTx(tx pgx.Tx, ctx context.Context, groupID int, winners, losers []qualifier, maxSeries, rangeSize int, rangeType string) error {
	var finalID int64
	err := tx.QueryRow(ctx, `INSERT INTO finals (group_id) VALUES ($1) RETURNING group_id`, groupID).Scan(&finalID)
	if err != nil {
		return fmt.Errorf("failed to create final record: %w", err)
	}

	sparringPairs := []struct {
		topPlace, botPlace int
		isGold             bool
	}{
		{topPlace: 5, botPlace: 6, isGold: true},
		{topPlace: 5, botPlace: 6, isGold: false},
	}

	for _, pair := range sparringPairs {
		var topWinner, botWinner *qualifier
		if pair.isGold {
			topWinner = findQualifier(winners, pair.topPlace)
			botWinner = findQualifier(winners, pair.botPlace)
		} else {
			topWinner = findQualifier(losers, pair.topPlace)
			botWinner = findQualifier(losers, pair.botPlace)
		}

		var topPlaceID, botPlaceID int64
		var state string

		if topWinner == nil && botWinner == nil {
			continue
		}

		if topWinner != nil {
			if botWinner == nil {
				var nullRangeGroupID sql.NullInt64
				topPlaceID, err = createSparringPlaceTx(tx, ctx, nullRangeGroupID, topWinner.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create top place for pair %v: %w", pair, err)
				}
				state = "top_win"
			} else {
				topPlaceID, err = createSparringPlaceWithRanges(tx, ctx, topWinner.CompetitorID, maxSeries, rangeSize, rangeType)
				if err != nil {
					return fmt.Errorf("failed to create top place with ranges for pair %v: %w", pair, err)
				}
			}
		}

		if botWinner != nil {
			if topWinner == nil {
				var nullRangeGroupID sql.NullInt64
				botPlaceID, err = createSparringPlaceTx(tx, ctx, nullRangeGroupID, botWinner.CompetitorID)
				if err != nil {
					return fmt.Errorf("failed to create bot place for pair %v: %w", pair, err)
				}
				state = "bot_win"
			} else {
				botPlaceID, err = createSparringPlaceWithRanges(tx, ctx, botWinner.CompetitorID, maxSeries, rangeSize, rangeType)
				if err != nil {
					return fmt.Errorf("failed to create bot place with ranges for pair %v: %w", pair, err)
				}
				state = "ongoing"
			}
		}

		if topPlaceID != 0 || botPlaceID != 0 {
			sparringID, err := createSparringTx(tx, ctx, topPlaceID, botPlaceID, state)
			if err != nil {
				return fmt.Errorf("failed to create sparring for pair %v: %w", pair, err)
			}

			updateField := "sparring_bronze_id"
			if pair.isGold {
				updateField = "sparring_gold_id"
			}
			_, err = tx.Exec(ctx, fmt.Sprintf(`UPDATE finals SET %s = $1 WHERE group_id = $2`, updateField), sparringID, finalID)
			if err != nil {
				return fmt.Errorf("failed to link sparring for pair %v: %w", pair, err)
			}
		}
	}
	return nil
}

func getFinalsTx(tx pgx.Tx, ctx context.Context, groupID int, f *models.Final) error {
	rows, err := tx.Query(ctx, `
        SELECT s.id, s.state, 
               sp_top.id, sp_top.competitor_id, sp_top.range_group_id, c_top.full_name,
               sp_bot.id, sp_bot.competitor_id, sp_bot.range_group_id, c_bot.full_name
        FROM finals f
        JOIN sparrings s ON f.sparring_gold_id = s.id OR f.sparring_bronze_id = s.id
        LEFT JOIN sparring_places sp_top ON s.top_place_id = sp_top.id
        LEFT JOIN competitors c_top ON sp_top.competitor_id = c_top.id
        LEFT JOIN sparring_places sp_bot ON s.bot_place_id = sp_bot.id
        LEFT JOIN competitors c_bot ON sp_bot.competitor_id = c_bot.id
        WHERE f.group_id = $1
        ORDER BY
            CASE
                WHEN s.id = f.sparring_gold_id THEN 1
                WHEN s.id = f.sparring_bronze_id THEN 2
            END`, groupID)
	if err != nil {
		return fmt.Errorf("failed to query finals: %w", err)
	}
	defer rows.Close()

	var sparrings []*models.Sparring
	for rows.Next() {
		var s models.Sparring
		var topPlaceID, topCompetitorID, botPlaceID, botCompetitorID sql.NullInt64
		var topRangeGroupID, botRangeGroupID sql.NullInt64
		var topFullName, botFullName sql.NullString

		if err := rows.Scan(&s.ID, &s.State,
			&topPlaceID, &topCompetitorID, &topRangeGroupID, &topFullName,
			&botPlaceID, &botCompetitorID, &botRangeGroupID, &botFullName); err != nil {
			return fmt.Errorf("failed to scan sparring: %w", err)
		}

		s.TopPlace = &models.SparringPlace{ID: int(topPlaceID.Int64), IsActive: true}
		if topCompetitorID.Valid {
			s.TopPlace.Competitor = models.CompetitorShrinked{
				ID:       int(topCompetitorID.Int64),
				FullName: topFullName.String,
			}
		}
		if topRangeGroupID.Valid {
			s.TopPlace.RangeGroup.ID = int(topRangeGroupID.Int64)
			if err := getRangeGroup(&s.TopPlace.RangeGroup); err != nil {
				return fmt.Errorf("failed to get top range group: %w", err)
			}
		}

		s.BotPlace = &models.SparringPlace{ID: int(botPlaceID.Int64)}
		if botCompetitorID.Valid {
			s.BotPlace.Competitor = models.CompetitorShrinked{
				ID:       int(botCompetitorID.Int64),
				FullName: botFullName.String,
			}
			s.BotPlace.IsActive = true
		}
		if botRangeGroupID.Valid {
			s.BotPlace.RangeGroup.ID = int(botRangeGroupID.Int64)
			if err := getRangeGroup(&s.BotPlace.RangeGroup); err != nil {
				return fmt.Errorf("failed to get bot range group: %w", err)
			}
		}

		sparrings = append(sparrings, &s)
	}
	if err := rows.Err(); err != nil {
		return fmt.Errorf("error iterating sparring rows: %w", err)
	}

	if f.SparringGold.BotPlace.Competitor.FullName == "" {
		f = nil
		return nil
	}

	if len(sparrings) > 0 {
		f.SparringGold = *sparrings[0]
	}
	if len(sparrings) > 1 {
		f.SparringBronze = *sparrings[1]
	}

	return nil
}

func checkFinalsCompleted(tx pgx.Tx, ctx context.Context, groupID int) error {
	var ongoingCount int
	err := tx.QueryRow(ctx, `
        SELECT COUNT(s.id)
        FROM finals f
        JOIN sparrings s ON f.sparring_gold_id = s.id OR f.sparring_bronze_id = s.id
        WHERE f.group_id = $1 AND s.state = 'ongoing'`, groupID).Scan(&ongoingCount)
	if err != nil {
		return fmt.Errorf("failed to check finals completion: %w", err)
	}
	if ongoingCount > 0 {
		return fmt.Errorf("not all final sparrings are completed")
	}
	return nil
}

func getQualification(conn *pgx.Conn, groupID int, r *http.Request) (*models.QualificationTable, error) {
	var resp models.QualificationTable
	err := conn.QueryRow(context.Background(),
		`SELECT group_id, distance, round_count 
         FROM qualifications 
         WHERE group_id = $1`, groupID).Scan(
		&resp.GroupID, &resp.Distance, &resp.RoundCount)
	if err != nil {
		return nil, fmt.Errorf("failed to query qualifications: %w", err)
	}

	resp.Sections, err = getQualificationSections(groupID, r)
	if err != nil {
		return nil, fmt.Errorf("failed to get qualification sections: %w", err)
	}

	return &resp, nil
}

func getCompetitorGroup(rows pgx.Rows, groupID int) ([]map[string]interface{}, error) {
	var result []map[string]interface{}
	for rows.Next() {
		var c models.Competitor

		if err := rows.Scan(&c.ID, &c.FullName, &c.BirthDate, &c.Identity, &c.Bow, &c.Rank, &c.Region, &c.Federation, &c.Club); err != nil {
			return nil, err
		}

		result = append(result, map[string]interface{}{"group_id": groupID, "competitor": c})
	}
	return result, nil
}

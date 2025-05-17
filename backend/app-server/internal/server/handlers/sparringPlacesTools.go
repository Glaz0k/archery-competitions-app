package handlers

import (
	"app-server/internal/dto"
	"app-server/internal/models"
	"app-server/pkg/tools"
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

func getRangeGroup(rg *models.RangeGroup) error {
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed Acquire")
	}
	defer conn.Release()
	query := `SELECT
    r.id AS range_id,
    r.range_ordinal,
    r.is_active AS range_is_active,
    s.shot_ordinal,
    s.score
	FROM range_groups rg
	LEFT JOIN ranges r ON rg.id = r.group_id
	LEFT JOIN shots s ON r.id = s.range_id
	WHERE rg.id = $1
	ORDER BY r.range_ordinal ASC, s.shot_ordinal ASC`

	rows, err := conn.Query(context.Background(), query, rg.ID)
	defer rows.Close()

	if err != nil {
		return err
	}
	rangesMap := make(map[int]*models.Range)

	for rows.Next() {
		var (
			rangeID      sql.NullInt64
			rangeOrdinal sql.NullInt64
			isActive     sql.NullBool
			shotOrdinal  sql.NullInt64
			score        sql.NullString
		)

		err = rows.Scan(&rangeID, &rangeOrdinal, &isActive, &shotOrdinal, &score)
		if err != nil {
			return err
		}

		if !rangeID.Valid {
			continue
		}

		if _, exists := rangesMap[int(rangeID.Int64)]; !exists {
			rangesMap[int(rangeID.Int64)] = &models.Range{
				ID:           int(rangeID.Int64),
				RangeOrdinal: int(rangeOrdinal.Int64),
				IsActive:     isActive.Bool,
			}
		}

		if int(shotOrdinal.Int64) != 0 {
			shot := models.Shot{
				ShotOrdinal: int(shotOrdinal.Int64),
				Score:       &score.String,
			}
			if !score.Valid {
				shot.Score = nil
			}
			rangesMap[int(rangeID.Int64)].Shots = append(rangesMap[int(rangeID.Int64)].Shots, shot)
		}
	}
	if err = rows.Err(); err != nil {
		return err
	}
	for _, r := range rangesMap {
		r.RangeScore = r.CalculateScore()
		rg.Ranges = append(rg.Ranges, *r)
	}
	rg.TotalScore = rg.CalculateTotalScore()

	getRangeGroupQuery := `SELECT ranges_max_count, range_size, type FROM range_groups WHERE range_groups.id = $1`
	err = conn.QueryRow(context.Background(), getRangeGroupQuery, rg.ID).Scan(
		&rg.RangesMaxCount,
		&rg.RangeSize,
		&rg.Type)
	if err != nil {
		return err
	}
	return nil
}

func getRange(r *models.Range, sparringPlaceId, rangeOrdinal int) error {
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed Acquire")
	}
	defer conn.Release()
	query := `
        SELECT 
            r.id,
            r.range_ordinal,
            r.is_active
        FROM ranges r
        JOIN range_groups rg ON r.group_id = rg.id
        JOIN sparring_places sp ON rg.id = sp.range_group_id
        WHERE sp.id = $1 AND r.range_ordinal = $2
    `

	err = conn.QueryRow(context.Background(), query, sparringPlaceId, rangeOrdinal).Scan(
		&r.ID,
		&r.RangeOrdinal,
		&r.IsActive,
	)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return fmt.Errorf("range not found for sparring place %d and ordinal %d",
				sparringPlaceId, rangeOrdinal)
		}
		return fmt.Errorf("failed to get range: %w", err)
	}
	shotsQuery := `
        SELECT shot_ordinal, score 
        FROM shots 
        WHERE range_id = $1 
        ORDER BY shot_ordinal
    `

	rows, err := conn.Query(context.Background(), shotsQuery, r.ID)
	if err != nil {
		return fmt.Errorf("failed to get shots: %w", err)
	}
	defer rows.Close()

	r.Shots = make([]models.Shot, 0)
	for rows.Next() {
		var shot models.Shot
		if err := rows.Scan(&shot.ShotOrdinal, &shot.Score); err != nil {
			return fmt.Errorf("failed to scan shot: %w", err)
		}
		r.Shots = append(r.Shots, shot)
	}

	if err := rows.Err(); err != nil {
		return fmt.Errorf("shots rows error: %w", err)
	}
	r.RangeScore = r.CalculateScore()
	return nil
}

func checkRangeExist(spID, rangeOrdinal int) (bool, error) {
	var isRangeExist bool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed Acquire")
	}
	defer conn.Release()
	queryCheck := `SELECT EXISTS (SELECT 1 
		FROM ranges r 
		JOIN sparring_places sp ON r.group_id = sp.range_group_id
		WHERE sp.id = $1
		AND r.range_ordinal = $2)`
	err = conn.QueryRow(context.Background(), queryCheck, spID, rangeOrdinal).Scan(&isRangeExist)
	if err != nil {
		return false, err
	}
	return isRangeExist, nil
}

func checkRangeActive(spID, rangeOrdinal int) (bool, error) {
	var isActive bool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return false, fmt.Errorf("failed Acquire")
	}
	defer conn.Release()
	queryCheck := `SELECT is_active FROM ranges r 
		JOIN sparring_places sp ON r.group_id = sp.range_group_id
		WHERE sp.id = $1
		AND r.range_ordinal = $2`
	err = conn.QueryRow(context.Background(), queryCheck, spID, rangeOrdinal).Scan(&isActive)
	if err != nil {
		return false, err
	}
	return isActive, nil
}

func isValidScore(score string, rangeType string) bool {
	if score == "X" || score == "M" {
		return true
	}
	val, err := strconv.Atoi(score)
	if err != nil {
		return false
	}
	return val >= 1 && val <= 10 && rangeType == "1-10" || val >= 6 && val <= 10 && rangeType == "6-10"
}

func checkAccess(r *http.Request, spID int) (string, int, bool, error) {
	var competitorID, rangeGroupID int
	rangeGroupIDQuery := `SELECT competitor_id, range_group_id FROM sparring_places WHERE id = $1`
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return "", 0, false, fmt.Errorf("failed Acquire")
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), rangeGroupIDQuery,
		spID).Scan(&competitorID, &rangeGroupID)
	if err != nil {
		return "", 0, false, err
	}

	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		return "", 0, false, err
	}
	if role != "admin" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil {
			return "", 0, false, err
		}
		if competitorID != userID {
			return "", 0, false, err
		}
	}
	return role, rangeGroupID, true, nil
}

// TODO: Add Acquire error handle
func getShotOut(shot *models.ShootOuts, placeID int) bool {
	shootOutsQuery := `SELECT score, priority FROM shoot_outs WHERE place_id = $1`
	var score sql.NullString
	var priority sql.NullBool
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return false
	}
	defer conn.Release()
	err = conn.QueryRow(context.Background(), shootOutsQuery, placeID).Scan(&score, &priority)
	if err == nil {
		shot.ID = placeID
		shot.Score = score.String
		shot.Priority = priority.Bool
		return true
	}
	return false
}

func getOpponentPlaceID(spID int) (int, error) {
	var opSpID int
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return 0, fmt.Errorf("failed Acquire")
	}
	defer conn.Release()

	opponentQuery := `SELECT 
    CASE 
        WHEN top_place_id = $1 THEN bot_place_id 
        WHEN bot_place_id = $1 THEN top_place_id 
    END AS opponent_id
	FROM sparrings 
	WHERE top_place_id = $1 OR bot_place_id = $1`
	err = conn.QueryRow(context.Background(), opponentQuery, spID).Scan(&opSpID)
	if err != nil {
		return 0, err
	}
	return opSpID, nil
}

func calculateSparringPlaceScore(sparringPlace, opponentSparringPlace *models.SparringPlace, bowType string) {
	endedRanges := make(map[int]*models.RangeScorePair)
	for i := 0; i < len(sparringPlace.RangeGroup.Ranges); i++ {
		if !sparringPlace.RangeGroup.Ranges[i].IsActive {
			endedRanges[sparringPlace.RangeGroup.Ranges[i].RangeOrdinal] = &models.RangeScorePair{}
			endedRanges[sparringPlace.RangeGroup.Ranges[i].RangeOrdinal].
				CompScore = sparringPlace.RangeGroup.Ranges[i].RangeScore
		}
	}

	for i := 0; i < len(opponentSparringPlace.RangeGroup.Ranges); i++ {
		if !opponentSparringPlace.RangeGroup.Ranges[i].IsActive {
			_, exists := endedRanges[opponentSparringPlace.RangeGroup.Ranges[i].RangeOrdinal]
			if exists {
				endedRanges[opponentSparringPlace.RangeGroup.Ranges[i].RangeOrdinal].
					OppScore = opponentSparringPlace.RangeGroup.Ranges[i].RangeScore
			}
		}
	}
	if len(endedRanges) != 0 {
		compScore, oppScore := tools.CalculatePoints(endedRanges, bowType)
		sparringPlace.SparringScore = compScore
		opponentSparringPlace.SparringScore = oppScore
	}
}

func editRange(tx pgx.Tx, changeRange dto.ChangeRange, sparringPlaceId int) error {
	query := `UPDATE shots s
			SET score = $1 
			FROM ranges r 
			JOIN sparring_places sp ON r.group_id = sp.range_group_id
			WHERE s.range_id = r.id
			AND s.shot_ordinal = $2
			AND r.range_ordinal = $3
			AND sp.id = $4`
	for _, s := range changeRange.Shots {
		_, err := tx.Exec(context.Background(), query, s.Score, s.ShotOrdinal,
			changeRange.RangeOrdinal, sparringPlaceId)
		if err != nil {
			return err
		}
	}
	return nil
}

func markRangeAsCompleted(ctx context.Context, tx pgx.Tx, groupID, ordinal int) error {
	_, err := tx.Exec(ctx, `
		UPDATE ranges
		SET is_active = false
		WHERE group_id = $1 AND range_ordinal = $2
	`, groupID, ordinal)
	return err
}

func activateNextRange(ctx context.Context, tx pgx.Tx, groupID, currentOrdinal int) error {
	_, err := tx.Exec(ctx, `
		UPDATE ranges
		SET is_active = true
		WHERE group_id = $1 AND range_ordinal = $2
	`, groupID, currentOrdinal+1)
	return err
}

func checkAllShotsNotNull(ctx context.Context, conn *pgxpool.Conn, rangeOrdinal int, rangeGroupID int) (bool, error) {
	query := `
        SELECT NOT EXISTS (
            SELECT 1 
            FROM shots s
            JOIN ranges r ON s.range_id = r.id
            WHERE r.range_ordinal = $1 
              AND r.group_id = $2 
              AND s.score IS NULL
        ) AS all_shots_not_null
    `

	var allShotsNotNull bool
	err := conn.QueryRow(ctx, query, rangeOrdinal, rangeGroupID).Scan(&allShotsNotNull)
	if err != nil {
		return false, fmt.Errorf("error while checking all shots not null: %v", err)
	}

	return allShotsNotNull, nil
}

func endRange(ctx context.Context, tx pgx.Tx,
	sparringPlaceID, groupID, ordinal int, bowType string) (*models.Range, error) {
	var r *models.Range
	var err error
	switch bowType {
	case "block":
		r, err = endRangeCompound(ctx, tx, sparringPlaceID, groupID, ordinal, bowType)
	default:
		r, err = endRangeWinPoints(ctx, tx, sparringPlaceID, groupID, ordinal, bowType)
	}
	err = markRangeAsCompleted(ctx, tx, groupID, ordinal)
	if err != nil {
		return nil, err
	}
	r.IsActive = false
	return r, err
}

func endRangeWinPoints(ctx context.Context, tx pgx.Tx, sparringPlaceID, rgID,
	rangeOrdinal int, bowType string) (*models.Range, error) {
	getRangesMaxCountQuery := `SELECT ranges_max_count FROM range_groups WHERE id = $1`
	var rangeSize int
	err := tx.QueryRow(ctx, getRangesMaxCountQuery, rgID).Scan(&rangeSize)
	if err != nil {
		return nil, err
	}
	endedRange := models.Range{}

	if rangeOrdinal < 3 {
		err = activateNextRange(ctx, tx, rgID, rangeOrdinal)
		if err != nil {
			return nil, err
		}
		err = getRange(&endedRange, sparringPlaceID, rangeOrdinal)
		if err != nil {
			return nil, err
		}
		return &endedRange, nil
	}

	opPlaceID, err := getOpponentPlaceID(sparringPlaceID)
	if err != nil {
		return nil, err
	}
	opRG := models.RangeGroup{}
	opRG.Ranges = make([]models.Range, rangeSize)
	for i := 1; i <= rangeOrdinal; i++ {
		err = getRange(&opRG.Ranges[i], opPlaceID, i)
		if err != nil {
			return nil, err
		}
		if opRG.Ranges[i].IsActive {
			err = getRange(&endedRange, sparringPlaceID, rangeOrdinal)
			if err != nil {
				return nil, err
			}
			return &endedRange, nil
		}
	}
	sp := models.SparringPlace{}
	sp.RangeGroup.ID = rgID
	spOpponent := models.SparringPlace{}
	getRGIDopponent := `SELECT range_group_id FROM sparring_places WHERE id = $1`
	err = tx.QueryRow(context.Background(), getRGIDopponent, opPlaceID).
		Scan(&spOpponent.RangeGroup.ID)
	if err != nil {
		return nil, err
	}
	err = getRangeGroup(&sp.RangeGroup)
	if err != nil {
		return nil, err
	}
	err = getRangeGroup(&spOpponent.RangeGroup)
	if err != nil {
		return nil, err
	}
	calculateSparringPlaceScore(&sp, &spOpponent, bowType)
	if sp.SparringScore > spOpponent.SparringScore {
		err = endSparring(ctx, tx, sparringPlaceID, opPlaceID, rgID)
		if err != nil {
			return nil, err
		}
	} else if sp.SparringScore < spOpponent.SparringScore {
		err = endSparring(ctx, tx, opPlaceID, sparringPlaceID, rgID)
		if err != nil {
			return nil, err
		}
	} else {
		if rangeOrdinal < 5 {
			err = activateNextRange(ctx, tx, rgID, rangeOrdinal+1)
			if err != nil {
				return nil, err
			}
			err = activateNextRange(ctx, tx, opRG.ID, rangeOrdinal+1)
			if err != nil {
				return nil, err
			}
		} else {
			sparringID, err := getSparringIDForPlaces(ctx, tx, sparringPlaceID, opPlaceID)
			if err != nil {
				return nil, err
			}
			err = createShootOuts(ctx, tx, sparringID, sparringPlaceID, opPlaceID)
			if err != nil {
				return nil, err
			}
		}
	}
	err = getRange(&endedRange, sparringPlaceID, rangeOrdinal)
	if err != nil {
		return nil, err
	}
	return &endedRange, nil
}

func endRangeCompound(ctx context.Context, tx pgx.Tx, sparringPlaceID, rgID, rangeOrdinal int, bowType string) (*models.Range, error) {
	getRangesMaxCountQuery := `SELECT ranges_max_count FROM range_groups WHERE id = $1`
	var rangeSize int
	err := tx.QueryRow(ctx, getRangesMaxCountQuery, rgID).Scan(&rangeSize)
	if err != nil {
		return nil, err
	}
	endedRange := models.Range{}

	// Если серия не последняя
	if rangeOrdinal < rangeSize {
		err = activateNextRange(ctx, tx, rgID, rangeOrdinal)
		if err != nil {
			return nil, err
		}
		err = getRange(&endedRange, sparringPlaceID, rangeOrdinal)
		if err != nil {
			return nil, err
		}
		return &endedRange, nil
	}

	opPlaceID, err := getOpponentPlaceID(sparringPlaceID)
	if err != nil {
		return nil, err
	}
	opRG := models.RangeGroup{}
	opRG.Ranges = make([]models.Range, rangeSize)
	for i := 1; i <= rangeSize; i++ {
		err = getRange(&opRG.Ranges[i], opPlaceID, i)
		if err != nil {
			return nil, err
		}
		if opRG.Ranges[i].IsActive {
			err = getRange(&endedRange, sparringPlaceID, rangeOrdinal)
			if err != nil {
				return nil, err
			}
			return &endedRange, nil
		}
	}

	sp := models.SparringPlace{}
	sp.RangeGroup.ID = rgID
	spOpponent := models.SparringPlace{}
	getRGIDopponent := `SELECT range_group_id FROM sparring_places WHERE id = $1`
	err = tx.QueryRow(context.Background(), getRGIDopponent, opPlaceID).Scan(&spOpponent.RangeGroup.ID)
	if err != nil {
		return nil, err
	}
	err = getRangeGroup(&sp.RangeGroup)
	if err != nil {
		return nil, err
	}
	err = getRangeGroup(&spOpponent.RangeGroup)
	if err != nil {
		return nil, err
	}
	calculateSparringPlaceScore(&sp, &spOpponent, bowType)
	if sp.SparringScore > spOpponent.SparringScore {
		err = endSparring(ctx, tx, sparringPlaceID, opPlaceID, rgID)
		if err != nil {
			return nil, err
		}
	} else if sp.SparringScore < spOpponent.SparringScore {
		err = endSparring(ctx, tx, opPlaceID, sparringPlaceID, rgID)
		if err != nil {
			return nil, err
		}
	} else {
		sparringID, err := getSparringIDForPlaces(ctx, tx, sparringPlaceID, opPlaceID)
		if err != nil {
			return nil, err
		}
		err = createShootOuts(ctx, tx, sparringID, sparringPlaceID, opPlaceID)
		if err != nil {
			return nil, err
		}
	}
	err = getRange(&endedRange, sparringPlaceID, rangeOrdinal)
	if err != nil {
		return nil, err
	}
	return &endedRange, nil
}

func getSparringIDForPlaces(ctx context.Context, tx pgx.Tx, placeID1, placeID2 int) (int, error) {
	var sparringID int
	err := tx.QueryRow(ctx, `
		SELECT id FROM sparrings 
		WHERE (top_place_id = $1 AND bot_place_id = $2)
		   OR (top_place_id = $2 AND bot_place_id = $1)
	`, placeID1, placeID2).Scan(&sparringID)
	if err != nil {
		return 0, fmt.Errorf("places don't belong to the same sparring: %w", err)
	}
	return sparringID, nil
}

func endSparring(ctx context.Context, tx pgx.Tx, winnerPlaceID, loserPlaceID, rangeGroupID int) error {
	sparringID, err := getSparringIDForPlaces(ctx, tx, winnerPlaceID, loserPlaceID)
	if err != nil {
		return fmt.Errorf("failed to get sparring ID: %w", err)
	}

	var currentState string
	err = tx.QueryRow(ctx, `
		SELECT state FROM sparrings WHERE id = $1
	`, sparringID).Scan(&currentState)
	if err != nil {
		return fmt.Errorf("failed to get current sparring state: %w", err)
	}

	if currentState != "ongoing" {
		return errors.New("sparring is already completed")
	}

	var topPlaceID, botPlaceID int
	err = tx.QueryRow(ctx, `
		SELECT top_place_id, bot_place_id FROM sparrings WHERE id = $1
	`, sparringID).Scan(&topPlaceID, &botPlaceID)
	if err != nil {
		return fmt.Errorf("failed to get sparring places: %w", err)
	}

	var newState string
	if winnerPlaceID == topPlaceID && loserPlaceID == botPlaceID {
		newState = "top_win"
	} else if winnerPlaceID == botPlaceID && loserPlaceID == topPlaceID {
		newState = "bot_win"
	} else {
		return errors.New("invalid winner/loser combination for this sparring")
	}

	_, err = tx.Exec(ctx, `
		UPDATE sparrings SET state = $1 WHERE id = $2
	`, newState, sparringID)
	if err != nil {
		return fmt.Errorf("failed to update sparring state: %w", err)
	}

	_, err = tx.Exec(ctx, `
		UPDATE ranges SET is_active = false 
		WHERE group_id = $1 AND is_active = true
	`, rangeGroupID)
	if err != nil {
		return fmt.Errorf("failed to deactivate ranges: %w", err)
	}
	return nil
}

func createShootOuts(ctx context.Context, tx pgx.Tx, sparringID int, topPlaceID int, botPlaceID int) error {
	var exists bool
	err := tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM sparrings WHERE id = $1)`,
		sparringID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check sparring existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("sparring with ID %d not found", sparringID)
	}

	err = tx.QueryRow(ctx,
		`SELECT EXISTS(SELECT 1 FROM sparring_places WHERE id = $1 OR id = $2)`,
		topPlaceID, botPlaceID,
	).Scan(&exists)
	if err != nil {
		return fmt.Errorf("failed to check places existence: %w", err)
	}
	if !exists {
		return fmt.Errorf("places %d, %d not found", topPlaceID, botPlaceID)
	}

	_, err = tx.Exec(ctx,
		`INSERT INTO shoot_outs (place_id, score, priority) 
         VALUES ($1, NULL, NULL), ($2, NULL, NULL) 
         ON CONFLICT (place_id) DO NOTHING`,
		topPlaceID, botPlaceID,
	)
	if err != nil {
		return fmt.Errorf("failed to create shoot-outs: %w", err)
	}
	_, err = tx.Exec(ctx,
		`UPDATE sparrings SET state = 'ongoing' 
         WHERE id = $1 AND state != 'ongoing'`,
		sparringID,
	)
	if err != nil {
		return fmt.Errorf("failed to update sparring state: %w", err)
	}
	return nil
}

func determineWinner(topScore string, topPriority bool, botScore string, botPriority bool) int {
	topVal, err1 := strconv.Atoi(topScore)
	botVal, err2 := strconv.Atoi(botScore)
	if err1 != nil || err2 != nil {
		return 0
	}
	if topVal > botVal {
		return 1
	}
	if botVal > topVal {
		return 2
	}
	if topPriority && !botPriority {
		return 1
	}
	if botPriority && !topPriority {
		return 2
	}
	return 0
}

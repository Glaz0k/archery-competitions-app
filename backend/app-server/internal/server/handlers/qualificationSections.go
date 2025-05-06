package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"sort"
	"strconv"

	"app-server/internal/dto"

	"app-server/internal/models"
	"app-server/pkg/tools"

	"github.com/jackc/pgx/v5"
)

func StartQualification(w http.ResponseWriter, r *http.Request) {
	type request struct {
		Distance    string `json:"distance"`
		RoundCount  int    `json:"round_count"`
		RangesCount int    `json:"ranges_count"`
		RangeSize   int    `json:"range_size"`
	}
	var req request
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	individualGroupID, err := tools.ParseParamToInt(r, "group_id")
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
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to begin transaction: %v", err)})
		return
	}
	defer tx.Rollback(context.Background())

	var groupExists bool
	err = tx.QueryRow(context.Background(), `SELECT EXISTS(SELECT 1 FROM individual_groups WHERE id = $1)`, individualGroupID).Scan(&groupExists)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to check group existence: %v", err)})
		return
	}
	if !groupExists {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}

	var hasQualification bool
	err = tx.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM qualifications WHERE group_id = $1)", individualGroupID).Scan(&hasQualification)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to check qualification: %v", err)})
		return
	}

	if hasQualification {
		res, err := getQualification(conn.Conn(), individualGroupID, r)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to get qualification: %v", err)})
			return
		}
		tools.WriteJSON(w, http.StatusOK, res)
		return
	}

	_, err = tx.Exec(context.Background(), `INSERT INTO qualifications (group_id, distance, round_count) VALUES ($1, $2, $3)`,
		individualGroupID, req.Distance, req.RoundCount)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create qualification: %v", err)})
		return
	}

	var bowType string
	err = tx.QueryRow(context.Background(), `SELECT bow FROM individual_groups WHERE id = $1`, individualGroupID).Scan(&bowType)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to get group bow type: %v", err)})
		return
	}

	rangeType := "1-10"
	if bowType == "classic" || bowType == "block" {
		rangeType = "6-10"
	}

	var competitorIDs []int
	rows, err := tx.Query(context.Background(), `SELECT competitor_id FROM competitor_group_details
    	WHERE group_id = $1`, individualGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to get competitors: %v", err)})
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

	if len(competitorIDs) < 5 {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "few participants"})
		return
	}

	i := 1
	for _, competitorID := range competitorIDs {
		var sectionID int
		err = tx.QueryRow(context.Background(), `INSERT INTO qualification_sections (group_id, competitor_id, place) VALUES ($1, $2, $3)
             RETURNING id`, individualGroupID, competitorID, i).Scan(&sectionID)
		i++
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create qualification section: %v", err)})
			return
		}

		first := true
		for roundOrdinal := 1; roundOrdinal <= req.RoundCount; roundOrdinal++ {
			var rangeGroupID int
			err = tx.QueryRow(context.Background(), `INSERT INTO range_groups (ranges_max_count, range_size, type) VALUES ($1, $2, $3) RETURNING id`,
				req.RangesCount, req.RangeSize, rangeType).Scan(&rangeGroupID)
			if err != nil {
				tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create range group: %v", err)})
				return
			}

			for rangeOrdinal := 1; rangeOrdinal <= req.RangesCount; rangeOrdinal++ {
				var rangeID int
				err = tx.QueryRow(context.Background(), `INSERT INTO ranges (group_id, range_ordinal, is_active) VALUES ($1, $2, $3) RETURNING id`,
					rangeGroupID, rangeOrdinal, true).Scan(&rangeID)
				if err != nil {
					tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create range: %v", err)})
					return
				}

				for shotOrdinal := 1; shotOrdinal <= req.RangeSize; shotOrdinal++ {
					_, err = tx.Exec(context.Background(), `INSERT INTO shots (range_id, shot_ordinal, score) VALUES ($1, $2, $3)`,
						rangeID, shotOrdinal, nil)
					if err != nil {
						tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create shot: %v", err)})
						return
					}
				}
			}

			_, err = tx.Exec(context.Background(), `INSERT INTO qualification_rounds (section_id, round_ordinal, is_active, range_group_id)
                 VALUES ($1, $2, $3, $4)`, sectionID, roundOrdinal, first, rangeGroupID)
			if err != nil {
				tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create qualification round: %v", err)})
				return
			}
			if first {
				first = false
			}
		}
	}

	_, err = tx.Exec(context.Background(), `UPDATE individual_groups SET state = 'qualification_start' WHERE id = $1`, individualGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to update group state: %v", err)})
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to commit transaction: %v", err)})
		return
	}

	res, err := getQualification(conn.Conn(), individualGroupID, r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to get qualification: %v", err)})
		return
	}

	tools.WriteJSON(w, http.StatusCreated, res)
}

func EndQualification(w http.ResponseWriter, r *http.Request) {
	individualGroupID, err := tools.ParseParamToInt(r, "group_id")
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
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to begin transaction: %v", err)})
		return
	}

	defer tx.Rollback(context.Background())

	var incompleteRounds int
	err = tx.QueryRow(context.Background(), `SELECT COUNT(*) FROM qualification_rounds qr JOIN qualification_sections qs 
    	ON qr.section_id = qs.id WHERE qs.group_id = $1 AND qr.is_active = true`, individualGroupID).Scan(&incompleteRounds)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to check incomplete rounds: %v", err)})
		return
	}
	if incompleteRounds > 0 {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "not all rounds are completed"})
		return
	}

	err = editCompetitorsPlaces(tx, individualGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	_, err = tx.Exec(context.Background(), `UPDATE individual_groups SET state = 'qualification_end' WHERE id = $1`,
		individualGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to update group state: %v", err)})
		return
	}

	var distance string
	var roundCount int
	err = tx.QueryRow(context.Background(), `SELECT distance, round_count FROM qualifications WHERE group_id = $1`,
		individualGroupID).Scan(&distance, &roundCount)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to get qualification data: %v", err)})
		return
	}

	var resp models.QualificationTable
	resp.RoundCount = roundCount
	resp.Distance = distance
	resp.GroupID = individualGroupID
	resp.Sections, err = getQualificationSections(individualGroupID, r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to get qualification sections: %v", err)})
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to commit transaction: %v", err)})
		return
	}

	tools.WriteJSON(w, http.StatusOK, resp)
}

func GetQualificationSection(w http.ResponseWriter, r *http.Request) {
	sectionID, err := tools.ParseParamToInt(r, "id")
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
	var section models.QualificationSectionResponse
	var competitorID int
	var place sql.NullInt64
	var rankGained sql.NullString
	err = conn.QueryRow(r.Context(), `
        SELECT qs.id, qs.competitor_id, c.full_name, qs.place
        FROM qualification_sections qs
        JOIN competitors c ON qs.competitor_id = c.id
        WHERE qs.id = $1`, sectionID).Scan(&section.ID, &competitorID, &section.Competitor.FullName, &place)
	if errors.Is(err, pgx.ErrNoRows) {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("DATABASE ERROR", err.Error())})
		return
	}
	section.Competitor.ID = competitorID
	if place.Valid {
		section.Place = int(place.Int64)
	}
	if rankGained.Valid {
		section.RankGained = rankGained.String
	}

	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "UNAUTHORIZED"})
		return
	}
	if role != "admin" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil || userID != competitorID {
			tools.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "FORBIDDEN"})
			return
		}
	}

	rounds, totalScore, tensCount, ninesCount, err := getSectionRoundsStats(sectionID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	section.Rounds = make([]models.RoundResponse, len(rounds))
	for i, r := range rounds {
		section.Rounds[i] = models.RoundResponse{
			RoundOrdinal: r.RoundOrdinal,
			IsOngoing:    r.IsActive,
			Total:        r.TotalScore,
		}
	}
	section.Total = totalScore
	section.CountTen = tensCount
	section.CountNine = ninesCount

	tools.WriteJSON(w, http.StatusOK, section)
}

func GetQualificationRound(w http.ResponseWriter, r *http.Request) {
	sectionID, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}
	roundOrdinal, err := tools.ParseParamToInt(r, "round_ordinal")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID ROUND ORDINAL"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	var competitorID, rangeGroupID int
	var isActive bool
	err = conn.QueryRow(r.Context(), `
        SELECT qs.competitor_id, qr.range_group_id, qr.is_active
        FROM qualification_sections qs
        JOIN qualification_rounds qr ON qs.id = qr.section_id
        WHERE qs.id = $1 AND qr.round_ordinal = $2`, sectionID, roundOrdinal).Scan(&competitorID, &rangeGroupID, &isActive)
	if errors.Is(err, pgx.ErrNoRows) {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "UNAUTHORIZED"})
		return
	}
	if role != "admin" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil || userID != competitorID {
			tools.WriteJSON(w, http.StatusForbidden, map[string]string{"error": "FORBIDDEN"})
			return
		}
	}

	var round models.QualificationRoundResponse
	round.SectionID = sectionID
	round.RoundOrdinal = roundOrdinal
	round.IsActive = isActive
	round.RangeGroup.ID = rangeGroupID

	if err := getRangeGroup(&round.RangeGroup); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "FAILED TO GET RANGE GROUP"})
		return
	}

	roundScore, _, _, err := getRoundStats(sectionID, roundOrdinal)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "FAILED TO GET ROUND STATS"})
		return
	}
	round.RangeGroup.TotalScore = roundScore

	for i := range round.RangeGroup.Ranges {
		round.RangeGroup.Ranges[i].RangeScore = round.RangeGroup.Ranges[i].CalculateScore()
		if round.RangeGroup.Ranges[i].Shots == nil {
			round.RangeGroup.Ranges[i].Shots = []models.Shot{}
		}
		if round.RangeGroup.Ranges[i].RangeScore == 0 && len(round.RangeGroup.Ranges[i].Shots) == 0 {
			round.RangeGroup.Ranges[i].RangeScore = 0
		}
	}

	tools.WriteJSON(w, http.StatusOK, round)
}

func GetQualificationSectionRanges(w http.ResponseWriter, r *http.Request) {
	sectionId, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	round, err := tools.ParseParamToInt(r, "round_ordinal")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	if role == "user" {
		userID, err := tools.GetUserIDFromContext(r)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("%v", err)})
			return
		}

		var exist bool
		queryCheck := `SELECT  EXISTS (SELECT 1 FROM qualification_sections WHERE id = $1 AND competitor_id = $2)`
		err = conn.QueryRow(context.Background(), queryCheck, sectionId, userID).Scan(&exist)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !exist {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	var rangeGroup models.RangeGroup
	query := `SELECT rg.id, rg.ranges_max_count, rg.range_size
		FROM range_groups rg 
		JOIN qualification_rounds qr ON rg.id = qr.range_group_id
		WHERE qr.section_id = $1
		AND qr.round_ordinal = $2`

	err = conn.QueryRow(context.Background(), query, sectionId, round).Scan(&rangeGroup.ID, &rangeGroup.RangesMaxCount, &rangeGroup.RangeSize)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var rg models.Range
	var ranges []models.Range
	var totalScore int

	query = `SELECT id, range_ordinal, is_active FROM ranges WHERE group_id = $1`
	rows, err := conn.Query(context.Background(), query, rangeGroup.ID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()

	for rows.Next() {
		err = rows.Scan(&rg.ID, &rg.RangeOrdinal, &rg.IsActive)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		ranges = append(ranges, rg)
	}

	for i := range ranges {
		err = getShotsFromRange(&ranges[i])
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		rangeScore, err := getRangeScore(ranges[i].Shots)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
		ranges[i].RangeScore = rangeScore
		totalScore += rangeScore
	}

	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	rangeGroup.Ranges = ranges
	rangeGroup.TotalScore = totalScore
	tools.WriteJSON(w, http.StatusOK, rangeGroup)
}

func getShotsFromRange(r *models.Range) error {
	conn, err := dbPool.Acquire(context.Background())
	if err != nil {
		return fmt.Errorf("failed Acquire")
	}
	defer conn.Release()
	query := `SELECT shot_ordinal, score FROM shots WHERE range_id = $1`
	rows, err := conn.Query(context.Background(), query, r.ID)
	if err != nil {
		return err
	}
	defer rows.Close()

	var s models.Shot
	var shots []models.Shot
	var score sql.NullString
	for rows.Next() {
		err = rows.Scan(&s.ShotOrdinal, &score)
		if err != nil {
			return err
		}
		s.Score = score.String
		shots = append(shots, s)
	}
	r.Shots = shots
	return nil
}

func getRangeScore(shots []models.Shot) (int, error) {
	var rangeScore int
	pattern := regexp.MustCompile(`^(10|[1-9])$`)
	for _, shot := range shots {
		score := shot.Score
		if pattern.MatchString(score) {
			s, err := strconv.Atoi(score)
			if err != nil {
				return 0, fmt.Errorf("invalid score")
			}
			rangeScore += s
		} else if score == "X" {
			rangeScore += 10
		} else if !(score == "M" || score == "") {
			return 0, fmt.Errorf("invalid score")
		}
	}
	return rangeScore, nil
}

// TODO валидация значений score
func EditQualificationSectionRanges(w http.ResponseWriter, r *http.Request) {
	sectionId, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	round, err := tools.ParseParamToInt(r, "round_ordinal")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	var changeRange dto.ChangeRange
	err = json.NewDecoder(r.Body).Decode(&changeRange)
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
	var isRangeExist bool
	queryCheck := `SELECT EXISTS (SELECT 1 
		FROM ranges r 
		JOIN qualification_rounds qr ON r.group_id = qr.range_group_id
		WHERE qr.section_id = $1
		AND qr.round_ordinal = $2
		AND r.range_ordinal = $3)`
	err = conn.QueryRow(context.Background(), queryCheck, sectionId, round, changeRange.RangeOrdinal).Scan(&isRangeExist)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !isRangeExist {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	var isRangeActive bool
	var typeRange string
	var rangeGroupId int
	queryCheck = `SELECT r.is_active, rg.type, rg.id
		FROM ranges r 
		JOIN  range_groups rg ON rg.id = r.group_id
		JOIN qualification_rounds qr ON rg.id = qr.range_group_id
		WHERE qr.section_id = $1
		AND qr.round_ordinal = $2
		AND r.range_ordinal = $3`
	err = conn.QueryRow(context.Background(), queryCheck, sectionId, round, changeRange.RangeOrdinal).Scan(&isRangeActive, &typeRange, &rangeGroupId)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
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

		var exist bool
		queryCheck = `SELECT EXISTS (SELECT 1 FROM qualification_sections WHERE id = $1 AND competitor_id = $2)`
		err = conn.QueryRow(context.Background(), queryCheck, sectionId, userID).Scan(&exist)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !exist {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}

		if !isRangeActive {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
		if typeRange == "6-10" {
			for _, c := range changeRange.Shots {
				if c.Score != "X" && c.Score != "M" && c.Score < "6" {
					e := dto.ErrorInvalidType{
						Error: "INVALID SCORE",
						Details: dto.DetailsInvalidType{
							ShotOrdinal: c.ShotOrdinal,
							Type:        typeRange,
						},
					}
					tools.WriteJSON(w, http.StatusBadRequest, e)
					return
				}
			}
		}
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer tx.Rollback(context.Background())

	err = editRanges(tx, changeRange, sectionId, round)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var rg models.Range
	rg.RangeOrdinal = changeRange.RangeOrdinal
	query := `SELECT id, is_active FROM ranges WHERE group_id = $1 AND range_ordinal = $2`
	err = tx.QueryRow(context.Background(), query, rangeGroupId, changeRange.RangeOrdinal).Scan(&rg.ID, &rg.IsActive)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	err = getShotsFromRange(&rg)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	rg.RangeScore, err = getRangeScore(rg.Shots)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	if !isRangeActive {
		var individualGroupId int
		query = `SELECT group_id FROM qualification_sections WHERE id = $1`
		err = tx.QueryRow(context.Background(), query, sectionId).Scan(&individualGroupId)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}

		err = editCompetitorsPlaces(tx, individualGroupId)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
			return
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	tools.WriteJSON(w, http.StatusOK, rg)
}

func editRanges(tx pgx.Tx, changeRange dto.ChangeRange, sectionId int, round int) error {
	query := `UPDATE shots s
			SET score = $1 
			FROM ranges r 
			JOIN qualification_rounds qr ON r.group_id = qr.range_group_id
			WHERE s.range_id = r.id
			AND s.shot_ordinal = $2
			AND r.range_ordinal = $3
			AND qr.section_id = $4
			AND qr.round_ordinal = $5`
	for _, s := range changeRange.Shots {
		_, err := tx.Exec(context.Background(), query, s.Score, s.ShotOrdinal, changeRange.RangeOrdinal, sectionId, round)
		if err != nil {
			return err
		}
	}
	return nil
}

func editCompetitorsPlaces(tx pgx.Tx, individualGroupId int) error {
	type CompetitorResult struct {
		ID    int
		Total int
		Tens  int
		Nines int
	}

	var competitors []CompetitorResult
	rows, err := tx.Query(context.Background(), `
        SELECT 
            qs.competitor_id,
            SUM(CASE WHEN s.score = 'X' THEN 10 WHEN s.score = 'M' THEN 0 ELSE CAST(s.score AS INTEGER) END) as total,
            SUM(CASE WHEN s.score = 'X' THEN 1 WHEN s.score = '10' THEN 1 ELSE 0 END) as tens,
            SUM(CASE WHEN s.score = '9' THEN 1 ELSE 0 END) as nines
        FROM qualification_sections qs
        JOIN qualification_rounds qr ON qr.section_id = qs.id
        JOIN range_groups rg ON qr.range_group_id = rg.id
        JOIN ranges r ON r.group_id = rg.id
        JOIN shots s ON s.range_id = r.id
        WHERE qs.group_id = $1
        GROUP BY qs.competitor_id`, individualGroupId)
	if err != nil {
		return fmt.Errorf("unable to get competitor results: %v", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cr CompetitorResult
		if err := rows.Scan(&cr.ID, &cr.Total, &cr.Tens, &cr.Nines); err != nil {
			return fmt.Errorf("unable to scan competitor result: %v", err)
		}
		competitors = append(competitors, cr)
	}

	sort.Slice(competitors, func(i, j int) bool {
		if competitors[i].Total != competitors[j].Total {
			return competitors[i].Total > competitors[j].Total
		}
		if competitors[i].Tens != competitors[j].Tens {
			return competitors[i].Tens > competitors[j].Tens
		}
		return competitors[i].Nines > competitors[j].Nines
	})

	for place, competitor := range competitors {
		_, err = tx.Exec(context.Background(), `UPDATE qualification_sections SET place = $1 WHERE group_id = $2 AND competitor_id = $3`,
			place+1, individualGroupId, competitor.ID)
		if err != nil {
			return fmt.Errorf("unable to update competitor place: %v", err)
		}
	}
	return nil
}

func EndRange(w http.ResponseWriter, r *http.Request) {
	sectionId, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	round, err := tools.ParseParamToInt(r, "round_ordinal")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	rangeOrdinal, err := tools.ParseParamToInt(r, "range_ordinal")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	var isRangeExist bool
	queryCheck := `SELECT EXISTS (SELECT 1 
		FROM ranges r 
		JOIN qualification_rounds qr ON r.group_id = qr.range_group_id
		WHERE qr.section_id = $1
		AND qr.round_ordinal = $2
		AND r.range_ordinal = $3)`
	err = conn.QueryRow(context.Background(), queryCheck, sectionId, round, rangeOrdinal).Scan(&isRangeExist)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !isRangeExist {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	var isRangeActive bool
	queryCheck = `SELECT r.is_active
		FROM ranges r 
		JOIN qualification_rounds qr ON r.group_id = qr.range_group_id
		WHERE qr.section_id = $1
		AND qr.round_ordinal = $2
		AND r.range_ordinal = $3`
	err = conn.QueryRow(context.Background(), queryCheck, sectionId, round, rangeOrdinal).Scan(&isRangeActive)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !isRangeActive {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
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

		var exist bool
		queryCheck = `SELECT EXISTS (SELECT 1 FROM qualification_sections WHERE id = $1 AND competitor_id = $2)`
		err = conn.QueryRow(context.Background(), queryCheck, sectionId, userID).Scan(&exist)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !exist {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	queryCheck = `SELECT s.score
			FROM shots s
			JOIN ranges r ON s.range_id = r.id
			JOIN qualification_rounds qr ON r.group_id = qr.range_group_id
			WHERE qr.section_id = $1
			AND qr.round_ordinal = $2
			AND r.range_ordinal = $3`
	rows, err := conn.Query(context.Background(), queryCheck, sectionId, round, rangeOrdinal)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer rows.Close()

	var s sql.NullString
	for rows.Next() {
		err = rows.Scan(&s)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !s.Valid {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer tx.Rollback(context.Background())

	var rg models.Range
	query := `UPDATE ranges r
			SET is_active = $1 
			FROM qualification_rounds qr
			WHERE r.group_id = qr.range_group_id
			AND qr.section_id = $2
			AND qr.round_ordinal = $3
			AND r.range_ordinal = $4 
			RETURNING r.id`
	err = tx.QueryRow(context.Background(), query, false, sectionId, round, rangeOrdinal).Scan(&rg.ID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var existNextRange bool
	queryCheck = `SELECT EXISTS( SELECT 1
			FROM ranges r
			JOIN qualification_rounds qr ON r.group_id = qr.range_group_id
			WHERE qr.section_id = $1
			AND qr.round_ordinal = $2
			AND r.range_ordinal = $3)`
	err = tx.QueryRow(context.Background(), queryCheck, sectionId, round, rangeOrdinal+1).Scan(&existNextRange)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	if existNextRange {
		query = `UPDATE ranges r
			SET is_active = $1 
			FROM qualification_rounds qr
			WHERE r.group_id = qr.range_group_id
			AND qr.section_id = $2
			AND qr.round_ordinal = $3
			AND r.range_ordinal = $4`
		_, err = tx.Exec(context.Background(), query, true, sectionId, round, rangeOrdinal+1)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
	} else {
		query = `UPDATE qualification_rounds qr
			SET is_active = $1 
			WHERE qr.section_id = $2
			AND qr.round_ordinal = $3`
		_, err = tx.Exec(context.Background(), query, false, sectionId, round)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}

		var existNextRound bool
		queryCheck = `SELECT EXISTS (SELECT 1 FROM qualification_rounds WHERE section_id = $1 and round_ordinal = $2)`
		err = tx.QueryRow(context.Background(), queryCheck, sectionId, round+1).Scan(&existNextRound)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}

		if existNextRound {
			var groupId int
			query = `UPDATE qualification_rounds SET is_active = $1 WHERE section_id = $2 AND round_ordinal = $3 RETURNING range_group_id`
			err = tx.QueryRow(context.Background(), query, true, sectionId, round+1).Scan(&groupId)
			if err != nil {
				tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
				return
			}
			query = `UPDATE ranges SET is_active = $1 WHERE group_id = $2 AND range_ordinal = $3`
			_, err = tx.Exec(context.Background(), query, true, groupId, 1)
			if err != nil {
				tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
				return
			}
		}
	}

	rg.RangeOrdinal = rangeOrdinal
	rg.IsActive = false
	err = getShotsFromRange(&rg)
	rg.RangeScore, err = getRangeScore(rg.Shots)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": err.Error()})
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	tools.WriteJSON(w, http.StatusOK, rg)
}

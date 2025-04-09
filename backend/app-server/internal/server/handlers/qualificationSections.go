package handlers

import (
	"app-server/internal/dto"
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"sort"

	"app-server/internal/models"
	"app-server/pkg/tools"
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

	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to begin transaction: %v", err)})
		return
	}
	defer tx.Rollback(context.Background())

	var hasQualification bool
	err = tx.QueryRow(context.Background(),
		"SELECT EXISTS(SELECT 1 FROM qualifications WHERE group_id = $1)", individualGroupID).Scan(&hasQualification)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to check qualification: %v", err)})
		return
	}

	if hasQualification {
		err = deleteQualifications(context.Background(), tx, individualGroupID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to delete qualification: %v", err)})
			return
		}
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

	var rangeGroupID int
	err = tx.QueryRow(context.Background(), `INSERT INTO range_groups (ranges_max_count, range_size, type) VALUES ($1, $2, $3) RETURNING id`,
		req.RangesCount, req.RangeSize, rangeType).Scan(&rangeGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create range group: %v", err)})
		return
	}

	for i := 1; i <= req.RangesCount; i++ {
		_, err = tx.Exec(context.Background(), `INSERT INTO ranges (group_id, range_ordinal, is_active) VALUES ($1, $2, $3)`, rangeGroupID, i, false)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to create range: %v", err)})
			return
		}
	}

	var competitorIDs []int
	rows, err := tx.Query(context.Background(), `SELECT competitor_id  FROM competitor_group_details 
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

	var resp models.QualificationTable
	resp.RoundCount = req.RoundCount
	resp.Distance = req.Distance
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

func EndQualification(w http.ResponseWriter, r *http.Request) {
	individualGroupID, err := tools.ParseParamToInt(r, "group_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

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
	fmt.Println(incompleteRounds)
	if incompleteRounds > 0 {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "not all rounds are completed"})
		return
	}

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
        GROUP BY qs.competitor_id`, individualGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to get competitor results: %v", err)})
		return
	}
	defer rows.Close()

	for rows.Next() {
		var cr CompetitorResult
		if err := rows.Scan(&cr.ID, &cr.Total, &cr.Tens, &cr.Nines); err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to scan competitor result: %v", err)})
			return
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
			place+1, individualGroupID, competitor.ID)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": fmt.Sprintf("unable to update competitor place: %v", err)})
			return
		}
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

func CreateQualificationRound(w http.ResponseWriter, r *http.Request) {
	var qualificationRound models.QualificationRound
	err := json.NewDecoder(r.Body).Decode(&qualificationRound)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO qualification_rounds (section_id, round_ordinal, range_group_id) VALUES ($1, $2, $3)", qualificationRound.SectionID, qualificationRound.RoundNumber, qualificationRound.RangeGroupId)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateQualificationSection(w http.ResponseWriter, r *http.Request) {
	var qualificationSection models.QualificationSection
	err := json.NewDecoder(r.Body).Decode(&qualificationSection)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO qualification_sections (group_id, competitor_id, place) VALUES ($1, $2, $3)", qualificationSection.IndividualGroupsID, qualificationSection.CompetitorID, qualificationSection.Place)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetQualificationSection(w http.ResponseWriter, r *http.Request) {
	Qid, err := tools.ParseParamToInt(r, "competition_id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID ENDPOINT"})
		return
	}

	userID, ok := r.Context().Value("user_id").(int)
	if !ok {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "UserID not found"})
		return
	}
	role := r.Context().Value("role")
	if userID != Qid && role != "admin" {
		tools.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "You are not authorized to access this resource"})
		return
	}

	var qualificationSection models.QualificationSection
	err = conn.QueryRow(context.Background(), `SELECT * FROM qualification_sections WHERE id = $1`, Qid).Scan(&qualificationSection.ID, &qualificationSection.IndividualGroupsID, &qualificationSection.CompetitorID, &qualificationSection.Place)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "qualification section not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(qualificationSection); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func GetQualificationSectionRanges(w http.ResponseWriter, r *http.Request) {
	sectionId, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}

	round, err := tools.ParseParamToInt(r, "round_ordinal")
	if err != nil {
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

		var exist bool
		queryCheck := `SELECT  EXISTS (SELECT 1 FROM qualification_sections WHERE id = $1 AND competitor_id = $2)`
		err = conn.QueryRow(context.Background(), queryCheck, sectionId, userID).Scan(&exist)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !exist {
			tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "BAD ACTION"})
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
		err = rows.Scan(&rg.ID, &rg.RangeNumber, &rg.IsOngoing)
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
		err = rows.Scan(&s.ShotNumber, &score)
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
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}

	round, err := tools.ParseParamToInt(r, "round_ordinal")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}

	var changeRange dto.ChangeRange
	err = json.NewDecoder(r.Body).Decode(&changeRange)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	var isRangeExist bool
	queryCheck := `SELECT EXISTS (SELECT 1 
		FROM ranges r 
		JOIN  range_groups rg ON rg.id = r.group_id
		JOIN qualification_rounds qr ON rg.id = qr.range_group_id
		WHERE qr.section_id = $1
		AND qr.round_ordinal = $2
		AND r.range_ordinal = $3)`
	err = conn.QueryRow(context.Background(), queryCheck, sectionId, round, changeRange.RangeOrdinal).Scan(&isRangeExist)
	if err != nil {
		fmt.Println(err) //
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

	role, err := tools.GetRoleFromContext(r)
	if err != nil {
		fmt.Println(err) //
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
			fmt.Println(err) //
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !exist {
			tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "BAD ACTION"})
			return
		}

		if err != nil {
			fmt.Println(err) //
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		if !isRangeActive {
			tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "BAD ACTION"})
			return
		}
		if typeRange == "6-10" {
			for _, c := range changeRange.Shots {
				if c.Score != "X" && c.Score != "M" && c.Score < "6" {
					e := dto.ErrorInvalidType{
						Error: "INVALID SCORE",
						Details: dto.DetailsInvalidType{
							ShotOrdinal: c.ShotNumber,
							Type:        typeRange,
						},
					}
					tools.WriteJSON(w, http.StatusBadRequest, e)
					return
				}
			}
		}
	}
	fmt.Println(changeRange)
	err = editRanges(changeRange, sectionId, round)
	if err != nil {
		fmt.Println(err) //
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	var rg models.Range
	rg.RangeNumber = changeRange.RangeOrdinal
	query := `SELECT id, is_active FROM ranges WHERE group_id = $1 AND range_ordinal = $2`
	err = conn.QueryRow(context.Background(), query, rangeGroupId, changeRange.RangeOrdinal).Scan(&rg.ID, &rg.IsOngoing)
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

	//TODO for admin к завершенной пересчитать места

	tools.WriteJSON(w, http.StatusOK, rg)
}

func editRanges(changeRange dto.ChangeRange, sectionId int, round int) error {
	tx, err := conn.Begin(context.Background())
	if err != nil {
		return err
	}
	defer tx.Rollback(context.Background())

	query := `UPDATE shots s
			SET score = $1 
			FROM ranges r 
			JOIN range_groups rg ON rg.id = r.group_id
			JOIN qualification_rounds qr ON rg.id = qr.range_group_id
			WHERE s.range_id = r.id
			AND s.shot_ordinal = $2
			AND r.range_ordinal = $3
			AND qr.section_id = $4
			AND qr.round_ordinal = $5`
	for _, s := range changeRange.Shots {
		_, err = conn.Exec(context.Background(), query, s.Score, s.ShotNumber, changeRange.RangeOrdinal, sectionId, round)
		if err != nil {
			return err
		}
	}

	err = tx.Commit(context.Background())
	if err != nil {
		return err
	}
	return nil
}

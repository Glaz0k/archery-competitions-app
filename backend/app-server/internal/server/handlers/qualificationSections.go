package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"net/http"

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
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	individualGroupID, err := tools.ParseParamToInt(r, "individual_group_id")
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to begin transaction: %v", err))
		return
	}
	defer tx.Rollback(context.Background())

	_, err = tx.Exec(context.Background(), `INSERT INTO qualifications (group_id, distance, round_count) VALUES ($1, $2, $3)`,
		individualGroupID, req.Distance, req.RoundCount)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to create qualification: %v", err))
		return
	}

	var bowType string
	err = tx.QueryRow(context.Background(), `SELECT bow FROM individual_groups WHERE id = $1`, individualGroupID).Scan(&bowType)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to get group bow type: %v", err))
		return
	}

	rangeType := "6-10"
	if bowType == "classic" || bowType == "block" {
		rangeType = "1-10"
	}

	var rangeGroupID int
	err = tx.QueryRow(context.Background(), `INSERT INTO range_groups (ranges_max_count, range_size, type) VALUES ($1, $2, $3) RETURNING id`,
		req.RangesCount, req.RangeSize, rangeType).Scan(&rangeGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to create range group: %v", err))
		return
	}

	for i := 1; i <= req.RangesCount; i++ {
		_, err = tx.Exec(context.Background(), `INSERT INTO ranges (group_id, range_ordinal, is_active) VALUES ($1, $2, $3)`, rangeGroupID, i, false)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to create range: %v", err))
			return
		}
	}

	var competitorIDs []int
	rows, err := tx.Query(context.Background(), `SELECT competitor_id  FROM competitor_group_details 
         WHERE group_id = $1`, individualGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to get competitors: %v", err))
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

	i := 1
	for _, competitorID := range competitorIDs {
		var sectionID int
		err = tx.QueryRow(context.Background(), `INSERT INTO qualification_sections (group_id, competitor_id, place) VALUES ($1, $2, $3)
             RETURNING id`, individualGroupID, competitorID, i).Scan(&sectionID)
		i++
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to create qualification section: %v", err))
			return
		}

		first := true
		for roundOrdinal := 1; roundOrdinal <= req.RoundCount; roundOrdinal++ {
			_, err = tx.Exec(context.Background(), `INSERT INTO qualification_rounds (section_id, round_ordinal, is_active, range_group_id)
                 VALUES ($1, $2, $3, $4)`, sectionID, roundOrdinal, first, rangeGroupID)
			if err != nil {
				tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to create qualification round: %v", err))
				return
			}
			if first {
				first = false
			}
		}
	}

	_, err = tx.Exec(context.Background(), `UPDATE individual_groups SET state = 'qualification_start' WHERE id = $1`, individualGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to update group state: %v", err))
		return
	}

	var resp models.QualificationTable
	resp.RoundCount = req.RoundCount
	resp.Distance = req.Distance
	resp.GroupID = individualGroupID
	resp.Sections, err = getQualificationSections(individualGroupID, r)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to get qualification sections: %v", err))
		return
	}

	if err := tx.Commit(context.Background()); err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, fmt.Sprintf("unable to commit transaction: %v", err))
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

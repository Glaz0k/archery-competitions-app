package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"

	"app-server/internal/models"
	"app-server/pkg/tools"
)

func StartQualification(w http.ResponseWriter, r *http.Request) {
	var qualification models.Qualification
	err := json.NewDecoder(r.Body).Decode(&qualification)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	qualification.IndividualGroupID, err = tools.ParseParamToInt(r, "group_id")
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO qualifications (group_id, distance, round_count) VALUES ($1, $2, $3)", qualification.IndividualGroupID, qualification.Distance, qualification.RoundCount)
	if err != nil {
		http.Error(w, "unable to insert data", http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
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

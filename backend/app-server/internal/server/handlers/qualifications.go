package handlers

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

func StartQualification(w http.ResponseWriter, r *http.Request) {
	var qualification models.Qualification
	err := json.NewDecoder(r.Body).Decode(&qualification)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO qualifications (group_id, distance, round_count) VALUES ($1, $2, $3)", qualification.IndividualGroupID, qualification.Distance, qualification.RoundCount)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
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
	_, err = conn.Exec(context.Background(), "INSERT INTO qualification_rounds (section_id, round_number, range_group_id) VALUES ($1, $2, $3)", qualificationRound.SectionID, qualificationRound.RoundNumber, qualificationRound.RangeGroupId)
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

package handlers

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"

	"github.com/jackc/pgx/v5"
)

var conn *pgx.Conn

func InitDB(dbConn *pgx.Conn) {
	conn = dbConn
}

func CreateCup(w http.ResponseWriter, r *http.Request) {
	var cup models.Cup
	err := json.NewDecoder(r.Body).Decode(&cup)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO cups (title, address, season) VALUES ($1, $2, $3)", cup.Title, cup.Address, cup.Season)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateCompetition(w http.ResponseWriter, r *http.Request) {
	var competition models.Competition
	err := json.NewDecoder(r.Body).Decode(&competition)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO competitions (cup_id, stage, start_date, end_date, is_ended) VALUES ($1, $2, $3, $4, $5)", competition.CupID, competition.Stage, competition.StartDate, competition.EndDate, competition.IsEnded)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateIndividualGroup(w http.ResponseWriter, r *http.Request) {
	var individualGroup models.IndividualGroup
	err := json.NewDecoder(r.Body).Decode(&individualGroup)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO individual_groups (competition_id, bow, identity, state) VALUES ($1, $2, $3, $4)", individualGroup.CompetitionID, individualGroup.Bow, individualGroup.Identity, individualGroup.State)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateRangeGroup(w http.ResponseWriter, r *http.Request) {
	var rangeGroup models.RangeGroup
	err := json.NewDecoder(r.Body).Decode(&rangeGroup)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO range_groups (ranges_count, range_size) VALUES ($1, $2)", rangeGroup.RangesCount, rangeGroup.RangeSize)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateQualification(w http.ResponseWriter, r *http.Request) {
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

func CreateShot(w http.ResponseWriter, r *http.Request) {
	var shot models.Shot
	err := json.NewDecoder(r.Body).Decode(&shot)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO shots (range_id, shot_number, score) VALUES ($1, $2, $3)", shot.RangeId, shot.ShotNumber, shot.Score)
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
	_, err = conn.Exec(context.Background(), "INSERT INTO qualification_sections (group_id, competitor_id, place) VALUES ($1, $2, $3)", qualificationSection.IndividualGroupsID, qualificationSection.CompetitorsID, qualificationSection.Place)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func CreateRange(w http.ResponseWriter, r *http.Request) {
	var rangeGroup models.Range
	err := json.NewDecoder(r.Body).Decode(&rangeGroup)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO ranges (group_id, range_number, is_completed) VALUES ($1, $2, $3)", rangeGroup.GroupId, rangeGroup.RangeNumber, rangeGroup.IsCompleted)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

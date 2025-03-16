package handlers

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"app-server/internal/models"

	"github.com/gorilla/mux"
)

func CreateIndividualGroup(w http.ResponseWriter, r *http.Request) {
	var individualGroup models.IndividualGroup
	err := json.NewDecoder(r.Body).Decode(&individualGroup)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	vars := mux.Vars(r)
	competitionID := vars["competition_id"]
	competitionId, err := strconv.Atoi(competitionID)
	if err != nil {
		http.Error(w, "invalid competition_id", http.StatusBadRequest)
	}
	individualGroup.CompetitionID = competitionId
	_, err = conn.Exec(context.Background(), "INSERT INTO individual_groups (competition_id, bow, identity, state) VALUES ($1, $2, $3, $4)", individualGroup.CompetitionID, individualGroup.Bow, individualGroup.Identity, individualGroup.State)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func GetIndividualGroups(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["group_id"]
	groupId, err := strconv.Atoi(groupID)
	if err != nil {
		http.Error(w, "invalid group_id", http.StatusBadRequest)
		return
	}

	var individualGroup models.IndividualGroup
	err = conn.QueryRow(context.Background(), `SELECT * FROM individual_groups WHERE id = $1`, groupId).Scan(&individualGroup.ID, &individualGroup.CompetitionID, &individualGroup.Bow, &individualGroup.Identity, &individualGroup.State)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "individual group not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(individualGroup); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func DeleteGroup(w http.ResponseWriter, r *http.Request) {}

func UpdateGroup(w http.ResponseWriter, r *http.Request) {}

func GetCompetitors(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	groupID := vars["group_id"]
	groupId, err := strconv.Atoi(groupID)
	if err != nil {
		http.Error(w, "invalid group_id", http.StatusBadRequest)
		return
	}

	var competitor models.Competitor
	err = conn.QueryRow(context.Background(), `SELECT competitor_id FROM competitor_group_details WHERE group_id = $1`, groupId).Scan(&competitor.ID)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "competitor not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	err = conn.QueryRow(context.Background(), `SELECT full_name FROM competitors WHERE id = $1`, competitor.ID).Scan(&competitor.FullName)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			http.Error(w, "competitor not found", http.StatusNotFound)
		} else {
			http.Error(w, "database error", http.StatusInternalServerError)
		}
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(competitor); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

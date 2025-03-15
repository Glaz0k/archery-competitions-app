package handlers

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
	"strconv"

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

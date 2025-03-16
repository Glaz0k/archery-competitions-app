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

func CreateShot(w http.ResponseWriter, r *http.Request) {
	var shot models.Shot
	rangeId := mux.Vars(r)["range_id"]
	rangeID, err := strconv.Atoi(rangeId)
	if err != nil {
		http.Error(w, "invalid range_id", http.StatusBadRequest)
	}
	shot.RangeID = rangeID
	err = json.NewDecoder(r.Body).Decode(&shot)
	if err != nil {
		http.Error(w, "invalid request payload", http.StatusBadRequest)
		return
	}
	_, err = conn.Exec(context.Background(), "INSERT INTO shots (range_id, shot_number, score) VALUES ($1, $2, $3)", shot.RangeID, shot.ShotNumber, shot.Score)
	if err != nil {
		log.Fatalf("unable to insert data: %v\n", err)
		return
	}
	w.WriteHeader(http.StatusOK)
}

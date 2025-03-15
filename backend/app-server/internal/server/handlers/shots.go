package handlers

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

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

package handlers

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"log"
	"net/http"
)

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

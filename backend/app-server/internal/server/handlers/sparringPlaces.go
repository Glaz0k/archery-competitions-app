package handlers

import (
	"app-server/internal/models"
	"app-server/pkg/tools"
	"context"
	"net/http"
)

func GetSparringPlace(w http.ResponseWriter, r *http.Request, isAdmin bool) {
	id, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}

	var sparringPlace models.SparringPlace

	checkCompetitorQuery := `SELECT competitor_id FROM sparring_places WHERE id = $1`
	err = conn.QueryRow(context.Background(), checkCompetitorQuery, id).Scan(&sparringPlace.Competitor.ID)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	if !isAdmin {
		if sparringPlace.Competitor.ID != id {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}
	sparringPlace.ID = id
	getName := `SELECT full_name FROM competitors WHERE id = $1`
	err = conn.QueryRow(context.Background(), getName, sparringPlace.Competitor.ID).Scan(&sparringPlace.Competitor.FullName)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}
	// TODO: посчитать очки завершенных серий
	tools.WriteJSON(w, http.StatusOK, sparringPlace)
}

func GetSparringPlaceAdmin(w http.ResponseWriter, r *http.Request) {
	GetSparringPlace(w, r, true)
}

func GetSparringPlaceUser(w http.ResponseWriter, r *http.Request) {
	GetSparringPlace(w, r, false)
}

func GetRanges(w http.ResponseWriter, r *http.Request, isAdmin bool) {

}

func GetRangesAdmin(w http.ResponseWriter, r *http.Request) {
	GetRanges(w, r, true)
}

func GetRangesUser(w http.ResponseWriter, r *http.Request) {
	GetRanges(w, r, false)
}

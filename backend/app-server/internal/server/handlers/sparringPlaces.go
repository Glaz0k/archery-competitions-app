package handlers

import (
	"app-server/internal/dto"
	"app-server/internal/models"
	"app-server/pkg/tools"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	"github.com/jackc/pgx/v5"
)

func GetSparringPlace(w http.ResponseWriter, r *http.Request) {
	id, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		fmt.Printf("err: %v\n", err)
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}
	var sparringPlace models.SparringPlace

	var bowType string
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	competitorSparringQuery := `SELECT c.id, c.full_name, c.bow 
	FROM sparring_places s
    JOIN competitors c ON s.competitor_id = c.id WHERE s.id = $1`
	err = conn.QueryRow(context.Background(), competitorSparringQuery,
		id).Scan(&sparringPlace.Competitor.ID, &sparringPlace.Competitor.FullName, &bowType)

	if err != nil {
		fmt.Printf("err1: %v\n", err)
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	_, rgID, acs, err := checkAccess(r, id)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !acs {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	sparringPlace.ID = id
	sparringPlace.RangeGroup.ID = rgID

	err = getRangeGroup(&sparringPlace.RangeGroup)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	sparringPlace.ShootOut = &models.ShootOuts{}
	shootOutExists := getShotOut(sparringPlace.ShootOut, id)
	if !shootOutExists {
		sparringPlace.ShootOut = nil
	}
	opponentSparringPlaceID, err := getOpponentPlaceID(id)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	opponentRangeGroupQuery := `SELECT range_group_id FROM sparring_places WHERE id = $1`
	var opponentRangeGroupID int
	err = conn.QueryRow(context.Background(), opponentRangeGroupQuery, opponentSparringPlaceID).Scan(&opponentRangeGroupID)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	opponentSparringPlace := models.SparringPlace{}
	opponentSparringPlace.RangeGroup.ID = opponentRangeGroupID
	err = getRangeGroup(&opponentSparringPlace.RangeGroup)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	calculateSparringPlaceScore(&sparringPlace, &opponentSparringPlace, bowType)
	tools.WriteJSON(w, http.StatusOK, sparringPlace)
}

func GetRanges(w http.ResponseWriter, r *http.Request) {
	id, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "NOT FOUND"})
		return
	}
	var rangeGroup models.RangeGroup

	_, rgID, acs, err := checkAccess(r, id)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !acs {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	rangeGroup.ID = rgID

	err = getRangeGroup(&rangeGroup)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, rangeGroup)
}

func EditSparringPlaceRange(w http.ResponseWriter, r *http.Request) {
	spID, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	role, rangeGroupID, acs, err := checkAccess(r, spID)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	if !acs {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	var changeRange dto.ChangeRange
	err = json.NewDecoder(r.Body).Decode(&changeRange)
	if err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID PARAMETERS"})
		return
	}

	isRangeExist, err := checkRangeExist(spID, changeRange.RangeOrdinal)
	if err != nil {
		fmt.Printf("err1: %v\n", err)
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !isRangeExist {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	var isRangeActive bool
	var typeRange string
	queryRange := `SELECT is_active, rg.type FROM ranges 
    			   JOIN range_groups rg ON ranges.group_id = rg.id 
                   WHERE group_id = $1 AND range_ordinal = $2`
	err = conn.QueryRow(context.Background(), queryRange, rangeGroupID, changeRange.RangeOrdinal).Scan(&isRangeActive, &typeRange)
	if err != nil {
		fmt.Printf("err2: %v\n", err)
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if role == "user" {
		if !isRangeActive {
			tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
			return
		}
	}
	for _, c := range changeRange.Shots {
		if c.Score == nil || !isValidScore(*c.Score, typeRange) {
			e := dto.ErrorInvalidType{
				Error: "INVALID SCORE",
				Details: dto.DetailsInvalidType{
					ShotOrdinal: c.ShotOrdinal,
					Type:        typeRange,
				},
			}
			tools.WriteJSON(w, http.StatusBadRequest, e)
			return
		}
	}
	tx, err := conn.Begin(context.Background())
	if err != nil {
		fmt.Printf("err2: %v\n", err)
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer tx.Rollback(context.Background())

	err = editRange(tx, changeRange, spID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	err = tx.Commit(context.Background())
	if err != nil {
		fmt.Printf("err4: %v\n", err)
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	spRange := models.Range{}
	err = getRange(&spRange, spID, changeRange.RangeOrdinal)
	if err != nil {
		fmt.Printf("err5: %v\n", err)
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	tools.WriteJSON(w, http.StatusOK, spRange)
}

func EndSparringPlaceRange(w http.ResponseWriter, r *http.Request) {
	sparringPlaceID, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	_, rgID, acs, err := checkAccess(r, sparringPlaceID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !acs {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}

	rangeOrdinal, err := tools.ParseParamToInt(r, "range_ordinal")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	isRangeExist, err := checkRangeExist(sparringPlaceID, rangeOrdinal)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !isRangeExist {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	isRangeActive, err := checkRangeActive(sparringPlaceID, rangeOrdinal)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	curRange := models.Range{}
	if !isRangeActive {
		err := getRange(&curRange, sparringPlaceID, rangeOrdinal)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
			return
		}
		tools.WriteJSON(w, http.StatusOK, curRange)
		return
	}

	isAllShootsNotNull, err := checkAllShotsNotNull(context.Background(), conn, rangeOrdinal, rgID)
	if !isAllShootsNotNull {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	tx, err := conn.Begin(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer tx.Rollback(context.Background())

	bowTypeQuery := `SELECT bow FROM competitors WHERE id = $1`
	var bowType string
	err = tx.QueryRow(context.Background(), bowTypeQuery, sparringPlaceID).Scan(&bowType)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}

	endedRange, err := endRange(context.Background(), tx, sparringPlaceID, rgID, rangeOrdinal, bowType)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	err = tx.Commit(context.Background())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	tools.WriteJSON(w, http.StatusOK, endedRange)
}

func EditShootOut(w http.ResponseWriter, r *http.Request) {
	sparringPlaceID, err := tools.ParseParamToInt(r, "id")
	if err != nil {
		tools.WriteJSON(w, http.StatusNotFound, map[string]string{"error": "NOT FOUND"})
		return
	}

	_, _, acs, err := checkAccess(r, sparringPlaceID)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	if !acs {
		fmt.Printf("err3: %v\n", err)
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "BAD ACTION"})
		return
	}
	conn, err := dbPool.Acquire(r.Context())
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "DATABASE ERROR"})
		return
	}
	defer conn.Release()
	var so models.ShootOuts
	if err := json.NewDecoder(r.Body).Decode(&so); err != nil {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID REQUEST"})
		return
	}

	if so.Score != "" && !isValidScore(so.Score, "1-10") {
		tools.WriteJSON(w, http.StatusBadRequest, map[string]string{"error": "INVALID SCORE"})
		return
	}

	_, err = conn.Exec(r.Context(),
		`UPDATE shoot_outs 
         SET score = $1, priority = $2 
         WHERE place_id = $3`,
		so.Score, so.Priority, sparringPlaceID,
	)
	so.ID = sparringPlaceID
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "UPDATE FAILED"})
		return
	}
	var sparringID, otherPlaceID int64
	var currentPlaceIsTop bool
	err = conn.QueryRow(r.Context(),
		`SELECT 
            s.id, 
            CASE WHEN s.top_place_id = $1 THEN s.bot_place_id ELSE s.top_place_id END,
            s.top_place_id = $1
         FROM sparrings s
         WHERE s.top_place_id = $1 OR s.bot_place_id = $1`,
		sparringPlaceID,
	).Scan(&sparringID, &otherPlaceID, &currentPlaceIsTop)
	if err != nil {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "SPARRING NOT FOUND"})
		return
	}

	var otherScore string
	err = conn.QueryRow(r.Context(),
		`SELECT score FROM shoot_outs WHERE place_id = $1`,
		otherPlaceID,
	).Scan(&otherScore)
	if err != nil && !errors.Is(err, pgx.ErrNoRows) {
		tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "CHECK FAILED"})
		return
	}

	if otherScore != "" && so.Score != "" {
		var topScore, botScore string
		var topPriority, botPriority bool

		if currentPlaceIsTop {
			topScore = so.Score
			topPriority = so.Priority
			err = conn.QueryRow(r.Context(),
				`SELECT score, priority FROM shoot_outs WHERE place_id = $1`,
				otherPlaceID,
			).Scan(&botScore, &botPriority)
		} else {
			botScore = so.Score
			botPriority = so.Priority
			err = conn.QueryRow(r.Context(),
				`SELECT score, priority FROM shoot_outs WHERE place_id = $1`,
				otherPlaceID,
			).Scan(&topScore, &topPriority)
		}

		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "RESULTS FETCH FAILED"})
			return
		}

		winner := determineWinner(topScore, topPriority, botScore, botPriority)
		newState := "ongoing"
		if winner == 1 {
			newState = "top_win"
		} else if winner == 2 {
			newState = "bot_win"
		}

		_, err = conn.Exec(r.Context(),
			`UPDATE sparrings SET state = $1 WHERE id = $2`,
			newState, sparringID,
		)
		if err != nil {
			tools.WriteJSON(w, http.StatusInternalServerError, map[string]string{"error": "SPARRING UPDATE FAILED"})
			return
		}
	}

	tools.WriteJSON(w, http.StatusOK, so)

}

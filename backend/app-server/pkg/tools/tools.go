package tools

import (
	"app-server/internal/models"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5/pgxpool"
)

func ParseParamToInt(r *http.Request, str string) (int, error) {
	vars := mux.Vars(r)
	result := vars[str]
	res, err := strconv.Atoi(result)
	return res, err
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	_ = json.NewEncoder(w).Encode(data)
}

func ExistsInDB(ctx context.Context, conn *pgxpool.Conn, query string, args ...interface{}) (bool, error) {
	var exists bool
	err := conn.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		return false, fmt.Errorf("database error: %w", err)
	}
	return exists, nil
}

func GetUserIDFromContext(r *http.Request) (int, error) {
	userID := r.Context().Value("user_id")
	if userID == nil {
		return -1, fmt.Errorf("user_id not found in context")
	}

	userId, ok := userID.(int)
	if !ok {
		return -1, fmt.Errorf("user_id has invalid type")
	}
	return userId, nil
}

func GetRoleFromContext(r *http.Request) (string, error) {
	role := r.Context().Value("role")
	if role == nil {
		return "", fmt.Errorf("role not found in context")
	}
	return role.(string), nil
}

func CalculatePoints(sp map[int]*models.RangeScorePair, bowType string) (int, int) {
	totalTop := 0
	totalBot := 0

	switch bowType {
	case "block":
		for _, v := range sp {
			totalTop += v.CompScore
			totalBot += v.OppScore
		}
	default:
		for _, v := range sp {
			if v.CompScore > v.OppScore {
				totalTop += 2
			} else if v.CompScore < v.OppScore {
				totalBot += 2
			} else {
				totalTop++
				totalBot++
			}
		}
	}

	return totalTop, totalBot
}

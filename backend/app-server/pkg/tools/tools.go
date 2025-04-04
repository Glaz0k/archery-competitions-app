package tools

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	"github.com/jackc/pgx/v5"
)

func ParseParamToInt(r *http.Request, str string) (int, error) {
	vars := mux.Vars(r)
	result := vars[str]
	res, err := strconv.Atoi(result)
	return res, err
}

func WriteJSON(w http.ResponseWriter, statusCode int, data interface{}) error {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	return json.NewEncoder(w).Encode(data)
}

func ExistsInDB(ctx context.Context, conn *pgx.Conn, query string, args ...interface{}) (bool, error) {
	var exists int

	err := conn.QueryRow(ctx, query, args...).Scan(&exists)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("database error: %w", err)
	}

	return true, nil
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

package handlers_test

import (
	"app-server/internal/server/handlers"
	"context"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/jackc/pgx/v5"
	"github.com/stretchr/testify/assert"
)

func setupTestDB(t *testing.T) *pgx.Conn {
	config, err := pgx.ParseConfig("postgres://root:root_password@localhost:5432/BowCompetitions?sslmode=disable")
	if err != nil {
		t.Fatalf("Unable to parse config: %v\n", err)
	}

	conn, err := pgx.ConnectConfig(context.Background(), config)
	if err != nil {
		t.Fatalf("Unable to connect to database: %v\n", err)
	}

	return conn
}

func teardownTestDB(t *testing.T, conn *pgx.Conn) {
	conn.Close(context.Background())
}

func TestCreateCup(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	handlers.InitDB(conn)

	reqBody := `{"title": "Test Cup", "address": "Test Address", "season": "2023"}`
	req, err := http.NewRequest("POST", "/cups", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handlers.CreateCup(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var title, address, season string
	err = conn.QueryRow(context.Background(), "SELECT title, address, season FROM cups WHERE title = $1", "Test Cup").Scan(&title, &address, &season)
	assert.NoError(t, err)
	assert.Equal(t, "Test Cup", title)
	assert.Equal(t, "Test Address", address)
	assert.Equal(t, "2023     ", season)
}

func TestCreateCompetition(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	handlers.InitDB(conn)

	reqBody := `{"cup_id": 1, "stage": "II", "start_date": "2023-01-01", "end_date": "2023-01-02", "is_ended": false}`
	req, err := http.NewRequest("POST", "/competitions", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handlers.CreateCompetition(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var cupID int
	var stage string
	var isEnded bool
	err = conn.QueryRow(context.Background(), "SELECT cup_id, stage, is_ended FROM competitions WHERE stage = $1", "II").Scan(&cupID, &stage, &isEnded)
	assert.NoError(t, err)
	assert.Equal(t, 1, cupID)
	assert.Equal(t, "II", stage)
	assert.Equal(t, false, isEnded)
}

func TestCreateIndividualGroup(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	handlers.InitDB(conn)

	reqBody := `{"competition_id": 1, "bow": "block", "identity": "male", "state": "created"}`
	req, err := http.NewRequest("POST", "/individual-groups", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handlers.CreateIndividualGroup(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var competitionID int
	var bow, identity, state string
	err = conn.QueryRow(context.Background(), "SELECT competition_id, bow, identity, state FROM individual_groups WHERE bow = $1", "block").Scan(&competitionID, &bow, &identity, &state)
	assert.NoError(t, err)
	assert.Equal(t, 1, competitionID)
	assert.Equal(t, "block", bow)
	assert.Equal(t, "male", identity)
	assert.Equal(t, "created", state)
}

func TestCreateQualification(t *testing.T) {
	conn := setupTestDB(t)
	defer teardownTestDB(t, conn)

	handlers.InitDB(conn)

	_, err := conn.Exec(context.Background(), "INSERT INTO individual_groups (id, competition_id, bow, identity, state) VALUES ($1, $2, $3, $4, $5)",
		2, 1, "block", "male", "created")
	if err != nil {
		t.Fatalf("Unable to insert individual group: %v\n", err)
	}

	reqBody := `{"group_id": 1, "distance": "12", "round_count": 5}`
	req, err := http.NewRequest("POST", "/qualifications", strings.NewReader(reqBody))
	if err != nil {
		t.Fatalf("Could not create request: %v\n", err)
	}
	req.Header.Set("Content-Type", "application/json")

	rr := httptest.NewRecorder()

	handlers.CreateQualification(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)

	var groupID, roundCount int
	var distance string
	err = conn.QueryRow(context.Background(), "SELECT group_id, distance, round_count FROM qualifications WHERE group_id = $1", 1).Scan(&groupID, &distance, &roundCount)
	assert.NoError(t, err)
	assert.Equal(t, 1, groupID)
	assert.Equal(t, "12", distance)
	assert.Equal(t, 5, roundCount)
}

//TODO: дописать тесты когда-нибудь

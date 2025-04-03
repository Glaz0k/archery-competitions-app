package dto

import (
	"app-server/internal/models"
	"time"
)

type Comprtitor struct {
	CompetitorID int `json:"competitor_id"`
}

type CompetitorCompetitionDetails struct {
	Competition_ID int                 `json:"competition_id"`
	Competitors    []models.Competitor `json:"competitor"`
	Is_active      bool                `json:"is_active"`
	Created_at     time.Time           `json:"created_at"`
}

package dto

import (
	"app-server/internal/models"
	"time"
)

type CompetitorCompetitionDetails struct {
	Competiton_ID int                 `json:"competiton_id"`
	Competitors   []models.Competitor `json:`
	Is_active     bool                `json:"is_active"`
	Created_at    time.Time           `json:"created_at"`
}

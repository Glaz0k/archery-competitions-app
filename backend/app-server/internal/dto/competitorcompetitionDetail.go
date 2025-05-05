package dto

import (
	"app-server/internal/models"
)

type Competitor struct {
	CompetitorID int `json:"competitor_id"`
}

type CompetitorCompetitionDetails struct {
	CompetitionID int               `json:"competition_id"`
	Competitor    models.Competitor `json:"competitor"`
	IsActive      bool              `json:"is_active"`
	CreatedAt     string            `json:"created_at"`
}

package dto

import "app-server/internal/models"

type CompetitorGroupDetail struct {
	GroupID     int                 `json:"group_id"`
	Competitors []models.Competitor `json:"competitor"`
}

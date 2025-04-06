package models

import "time"

type CompetitorCompetitionDetails struct {
	CompetitionID int       `json:"competiton_id"`
	CompetitorID  int       `json:"competitor_id"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

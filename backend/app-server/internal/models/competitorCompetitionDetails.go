package models

import "time"

type CompetitorCompetitionDetails struct {
	Competiton_ID int       `json:"competiton_id"`
	Competitor_ID int       `json:"competitor_id"`
	Is_active     bool      `json:"is_active"`
	Created_at    time.Time `json:"created_at"`
}

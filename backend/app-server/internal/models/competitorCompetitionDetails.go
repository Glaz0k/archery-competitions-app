package models

type CompetitorCompetitionDetails struct {
	CompetitionID int   `json:"competition_id"`
	CompetitorID  int   `json:"competitor_id"`
	IsActive      bool  `json:"is_active"`
	CreatedAt     *Date `json:"created_at"`
}

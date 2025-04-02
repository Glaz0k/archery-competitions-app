package models

import "time"

type Competitor struct {
	ID         int       `json:"id"`
	FullName   string    `json:"full_name"`
	BirthDate  time.Time `json:"birth_date"`
	Identity   string    `json:"identity"`
	Bow        string    `json:"bow"`
	Rank       string    `json:"rank"`
	Region     string    `json:"region"`
	Federation string    `json:"federation"`
	Club       string    `json:"club"`
}

type CompetitorCompetitionDetails struct {
	CompetitionID int       `json:"competition_id"`
	CompetitorID  int       `json:"competitor_id"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
}

type CompetitorGroupDetails struct {
	GroupID      int `json:"group_id"`
	CompetitorID int `json:"competitor_id"`
}

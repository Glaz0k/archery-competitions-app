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

package models

import "time"

type Competition struct {
	ID        int       `json:"id"`
	CupID     int       `json:"cup_id"`
	Stage     string    `json:"stage"` // TODO: enum
	StartDate time.Time `json:"start_date"`
	EndDate   time.Time `json:"end_date"`
	IsEnded   bool      `json:"is_ended"`
}

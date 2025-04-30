package models

type Competition struct {
	ID        int     `json:"id"`
	CupID     int     `json:"cup_id"`
	Stage     *string `json:"stage"`
	StartDate *Date   `json:"start_date"`
	EndDate   *Date   `json:"end_date"`
	IsEnded   bool    `json:"is_ended"`
}

var CompetitionUpdateData struct {
	StartDate *Date `json:"start_date"`
	EndDate   *Date `json:"end_date"`
}

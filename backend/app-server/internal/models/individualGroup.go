package models

type IndividualGroup struct {
	ID            int     `json:"id"`
	CompetitionID int     `json:"competition_id"`
	Bow           *string `json:"bow"`
	Identity      string  `json:"identity"`
	State         string  `json:"state"`
}

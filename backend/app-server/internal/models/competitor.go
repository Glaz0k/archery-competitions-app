package models

type Competitor struct {
	ID         int     `json:"id"`
	FullName   string  `json:"full_name"`
	BirthDate  *Date   `json:"birth_date"`
	Identity   string  `json:"identity"`
	Bow        *string `json:"bow"`
	Rank       *string `json:"rank"`
	Region     *string `json:"region"`
	Federation *string `json:"federation"`
	Club       *string `json:"club"`
}

type CompetitorGroupDetails struct {
	GroupID      int `json:"group_id"`
	CompetitorID int `json:"competitor_id"`
}

type CompetitorShrinked struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}

package models

type Competitor struct {
	ID         int    `json:"id"`
	FullName   string `json:"full_name"`
	BirthDate  string `json:"birth_date"`
	Identity   string `json:"identity"`
	Bow        string `json:"bow"`
	Rank       int    `json:"rank"`
	Region     string `json:"region"`
	Federation string `json:"federation"`
	Club       string `json:"club"`
}

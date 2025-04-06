package models

type SparringPlace struct {
	ID            int                `json:"id"`
	Competitor    CompetitorShrinked `json:"competitor"`
	RangeGroup    RangeGroup         `json:"range_group"`
	IsActive      bool               `json:"is_active"`
	ShootOut      bool               `json:"shoot_out"`
	SparringScore int                `json:"sparring_score"`
}

type Shot struct {
	ShotNumber int    `json:"shot_original"`
	Score      string `json:"score"`
}

type ShotOut struct {
	ID       int    `json:"id"`
	Score    string `json:"score"`
	Priority bool   `json:"priority"`
}

type Sparring struct {
	ID       int           `json:"id"`
	TopPlace SparringPlace `json:"top_place"`
	BotPlace SparringPlace `json:"bot_place"`
	State    string        `json:"state"`
}

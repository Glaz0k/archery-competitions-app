package models

type SparringPlace struct {
	ID            int                `json:"id"`
	Competitor    CompetitorShrinked `json:"competitor"`
	RangeGroup    RangeGroup         `json:"range_group"`
	IsActive      bool               `json:"is_active"`
	ShootOut      *ShootOuts         `json:"shoot_out"`
	SparringScore int                `json:"sparring_score"`
}

type RangeScorePair struct {
	CompScore   int
	OppScore    int
	IsEndedComp bool
	IsEndedOpp  bool
}

type Shot struct {
	ShotOrdinal int     `json:"shot_ordinal"`
	Score       *string `json:"score"`
}

type ShootOuts struct {
	ID       int    `json:"sparring_place_id"`
	Score    string `json:"score"`
	Priority bool   `json:"priority"`
}

type Sparring struct {
	ID       int            `json:"id"`
	TopPlace *SparringPlace `json:"top_place"`
	BotPlace *SparringPlace `json:"bot_place"`
	State    string         `json:"state"`
}

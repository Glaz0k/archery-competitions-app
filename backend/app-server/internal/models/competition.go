package models

type Shot struct {
	RangeId    int    `json:"range_id"`
	ShotNumber int    `json:"shot_number"`
	Score      string `json:"score"`
}

type Range struct {
	ID          int  `json:"id"`
	GroupId     int  `json:"group_id"`
	RangeNumber int  `json:"range_number"`
	IsCompleted bool `json:"is_completed"`
}

type RangeGroup struct {
	ID          int `json:"id"`
	RangesCount int `json:"ranges_count"`
	RangeSize   int `json:"range_size"`
}

type QualificationRound struct {
	SectionID    int `json:"section_id"`
	RoundNumber  int `json:"round_number"`
	RangeGroupId int `json:"range_group_id"`
}

type QualificationSection struct {
	ID                 int   `json:"id"`
	IndividualGroupsID []int `json:"groups_id"`
	CompetitorsID      []int `json:"competitors_id"`
	Place              int   `json:"place"`
}

type Qualification struct {
	IndividualGroupID int    `json:"group_id"`
	Distance          string `json:"distance"` // String?
	RoundCount        int    `json:"round_count"`
}

type IndividualGroup struct {
	ID            int    `json:"id"`
	CompetitionID int    `json:"competition_id"`
	Bow           string `json:"bow"`
	Identity      string `json:"identity"`
	State         string `json:"state"`
}

type Competition struct {
	ID        int    `json:"id"`
	CupID     int    `json:"cup_id"`
	Stage     string `json:"stage"` // TODO: enum
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"` // TODO: date
	IsEnded   bool   `json:"is_ended"`
}

type Cup struct {
	ID           int           `json:"id"`
	Title        string        `json:"title"`
	Address      string        `json:"address"`
	Season       string        `json:"season"`
	Competitions []Competition `json:"competitions"`
}

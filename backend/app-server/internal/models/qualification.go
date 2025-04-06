package models

type Shot struct {
	ShotNumber int    `json:"shot_number"`
	Score      string `json:"score"`
	RangeID    int    `json:"range_id"`
}
type Range struct {
	ID          int    `json:"id"`
	RangeNumber int    `json:"range_number"`
	IsOngoing   bool   `json:"is_ongoing"`
	Shots       []Shot `json:"shots"`
	RangeScore  int    `json:"range_score"`
}

type RangeGroup struct {
	ID             int     `json:"id"`
	RangesMaxCount int     `json:"ranges_max_count"`
	RangeSize      int     `json:"range_size"`
	Ranges         []Range `json:"ranges"`
	TotalScore     int     `json:"total_score"`
}

type QualificationRound struct {
	SectionID    int  `json:"section_id"`
	RoundNumber  int  `json:"round_number"`
	RangeGroupId int  `json:"range_group"`
	IsOngoing    bool `json:"is_ongoing"`
}

type CompetitorShrinked struct {
	ID       int    `json:"id"`
	FullName string `json:"full_name"`
}

type RoundShrinked struct {
	RoundOrdinal int  `json:"round_ordinal"`
	IsActive     bool `json:"is_active"`
	TotalScore   int  `json:"total_score"`
}

type QualificationSectionForTable struct {
	ID         int                `json:"id"`
	Competitor CompetitorShrinked `json:"competitor"`
	Place      int                `json:"place"`
	Round      []RoundShrinked    `json:"rounds"`
	Total      int                `json:"total"`
	CountTen   int                `json:"10_s"`
	CountNine  int                `json:"9_s"`
	RankGained string             `json:"rank_gained"`
}

type QualificationSection struct {
	ID                 int   `json:"id"`
	IndividualGroupsID []int `json:"groups_id"`
	CompetitorID       []int `json:"competitor_id"`
	Place              int   `json:"place"`
}

type Qualification struct {
	IndividualGroupID int    `json:"group_id"`
	Distance          string `json:"distance"`
	RoundCount        int    `json:"round_count"`
}

type QualificationTable struct {
	GroupID    int                            `json:"group_id"`
	Distance   string                         `json:"distance"`
	RoundCount int                            `json:"round_count"`
	Sections   []QualificationSectionForTable `json:"sections"`
}

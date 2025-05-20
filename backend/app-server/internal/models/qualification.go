package models

import "strconv"

type Range struct {
	ID           int    `json:"id"`
	RangeOrdinal int    `json:"range_ordinal"`
	IsActive     bool   `json:"is_active"`
	Shots        []Shot `json:"shots"`
	RangeScore   int    `json:"range_score"`
}

func (r Range) CalculateScore() int {
	score := 0
	for _, shot := range r.Shots {
		if shot.Score == nil {
			continue
		}
		if *shot.Score == "X" {
			score += 10
		} else if *shot.Score == "M" {
			score += 0
		} else if *shot.Score != "" {
			val, _ := strconv.Atoi(*shot.Score)
			score += val
		}
	}
	return score
}

type RangeGroup struct {
	ID             int     `json:"id"`
	RangesMaxCount int     `json:"ranges_max_count"`
	RangeSize      int     `json:"range_size"`
	Type           string  `json:"type"`
	Ranges         []Range `json:"ranges"`
	TotalScore     int     `json:"total_score"`
}

func (r RangeGroup) CalculateTotalScore() int {
	score := 0
	for _, r := range r.Ranges {
		score += r.CalculateScore()
	}
	return score
}

type QualificationRound struct {
	SectionID    int  `json:"section_id"`
	RoundNumber  int  `json:"round_number"`
	RangeGroupId int  `json:"range_group"`
	IsOngoing    bool `json:"is_ongoing"`
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
	RankGained *string            `json:"rank_gained"`
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

type QualificationSectionResponse struct {
	ID         int                `json:"id"`
	Competitor CompetitorShrinked `json:"competitor"`
	Place      int                `json:"place"`
	Rounds     []RoundResponse    `json:"rounds"`
	Total      int                `json:"total"`
	CountTen   int                `json:"10_s"`
	CountNine  int                `json:"9_s"`
	RankGained *string            `json:"rank_gained"`
}

type RoundResponse struct {
	RoundOrdinal int  `json:"round_ordinal"`
	IsOngoing    bool `json:"is_active"`
	Total        int  `json:"total_score"`
}

type QualificationRoundResponse struct {
	SectionID    int        `json:"section_id"`
	RoundOrdinal int        `json:"round_ordinal"`
	IsActive     bool       `json:"is_active"`
	RangeGroup   RangeGroup `json:"range_group"`
}

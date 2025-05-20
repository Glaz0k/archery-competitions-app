package dto

import "app-server/internal/models"

type EditSparringPlaceRangeResponse struct {
	RangeOrdinal int           `json:"range_ordinal"`
	Shots        []models.Shot `json:"shots"`
	Error        error         `json:"error"`
}

type Err struct {
	Error   string  `json:"error"`
	Details Details `json:"details"`
}

type Details struct {
	ShotOrdinal int    `json:"shot_ordinal"`
	Type        string `json:"type"`
}

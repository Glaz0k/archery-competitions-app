package dto

import "app-server/internal/models"

type ChangeRange struct {
	RangeOrdinal int           `json:"range_ordinal"`
	Shots        []models.Shot `json:"shots"`
}

type ErrorInvalidType struct {
	Error   string             `json:"error"`
	Details DetailsInvalidType `json:"details"`
}

type DetailsInvalidType struct {
	ShotOrdinal int    `json:"shot_ordinal"`
	Type        string `json:"type"`
}

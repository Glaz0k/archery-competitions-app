package models

type Shot struct {
	Range_ID    int    `json:"range_id"`
	Shot_Number int    `json:"shot_number"`
	Score       string `json:"score"`
}

package models

type Cup struct {
	ID      int    `json:"id"`
	Title   string `json:"title"`
	Address string `json:"address"`
	Season  string `json:"season"`
}

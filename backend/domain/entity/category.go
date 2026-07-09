package entity

import "time"

type Category struct {
	ID           string    `json:"id"`
	Name         string    `json:"name"`
	Abbreviation string    `json:"abbreviation"`
	Description  string    `json:"description"`
	IsActive     bool      `json:"is_active"`
	CreatedAt    time.Time `json:"created_at"`
}

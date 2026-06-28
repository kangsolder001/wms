package entity

import "time"

type Item struct {
	ID            string    `json:"id"`
	SKU           string    `json:"sku"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	Category      string    `json:"category"`
	UnitOfMeasure string    `json:"unit_of_measure"`
	Weight        float64   `json:"weight"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

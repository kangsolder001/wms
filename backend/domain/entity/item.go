package entity

import (
	"fmt"
	"time"
)

type Item struct {
	ID            string    `json:"id"`
	SKU           string    `json:"sku"`
	Name          string    `json:"name"`
	Description   string    `json:"description"`
	CategoryID    string    `json:"category_id"`
	Category      string    `json:"category"`
	Barcode       string    `json:"barcode"`
	UnitOfMeasure string    `json:"unit_of_measure"`
	Weight        float64   `json:"weight"`
	Length        float64   `json:"length"`
	Width         float64   `json:"width"`
	Height        float64   `json:"height"`
	IsActive      bool      `json:"is_active"`
	CreatedAt     time.Time `json:"created_at"`
	UpdatedAt     time.Time `json:"updated_at"`
}

func GenerateBarcode(sku, name, category string) string {
	return fmt.Sprintf("%s|%s|%s", sku, name, category)
}

package entity

import "time"

type StockTransfer struct {
	ID             string    `json:"id"`
	TransferNumber string    `json:"transfer_number"`
	FromLocationID string    `json:"from_location_id"`
	ToLocationID   string    `json:"to_location_id"`
	ItemID         string    `json:"item_id"`
	Quantity       float64   `json:"quantity"`
	Status         string    `json:"status"`
	Notes          string    `json:"notes"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

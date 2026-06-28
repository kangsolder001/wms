package entity

import "time"

type Stock struct {
	ID               string     `json:"id"`
	ItemID           string     `json:"item_id"`
	LocationID       string     `json:"location_id"`
	Quantity         float64    `json:"quantity"`
	ReservedQuantity float64    `json:"reserved_quantity"`
	BatchNumber      string     `json:"batch_number"`
	ExpiryDate       *time.Time `json:"expiry_date,omitempty"`
	UpdatedAt        time.Time  `json:"updated_at"`
}

type StockMovement struct {
	ID             string    `json:"id"`
	ItemID         string    `json:"item_id"`
	FromLocationID *string   `json:"from_location_id,omitempty"`
	ToLocationID   *string   `json:"to_location_id,omitempty"`
	Quantity       float64   `json:"quantity"`
	MovementType   string    `json:"movement_type"`
	ReferenceType  string    `json:"reference_type"`
	ReferenceID    string    `json:"reference_id"`
	Notes          string    `json:"notes"`
	CreatedBy      string    `json:"created_by"`
	CreatedAt      time.Time `json:"created_at"`
}

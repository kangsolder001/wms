package dto

type AdjustStockRequest struct {
	ItemID     string  `json:"item_id" validate:"required"`
	LocationID string  `json:"location_id" validate:"required"`
	Quantity   float64 `json:"quantity" validate:"required"`
	Notes      string  `json:"notes"`
}

type StockResponse struct {
	ID               string  `json:"id"`
	ItemID           string  `json:"item_id"`
	LocationID       string  `json:"location_id"`
	Quantity         float64 `json:"quantity"`
	ReservedQuantity float64 `json:"reserved_quantity"`
	BatchNumber      string  `json:"batch_number"`
	ItemSKU          string  `json:"item_sku,omitempty"`
	ItemName         string  `json:"item_name,omitempty"`
	LocationCode     string  `json:"location_code,omitempty"`
}

type StockMovementResponse struct {
	ID             string  `json:"id"`
	ItemID         string  `json:"item_id"`
	FromLocationID *string `json:"from_location_id,omitempty"`
	ToLocationID   *string `json:"to_location_id,omitempty"`
	Quantity       float64 `json:"quantity"`
	MovementType   string  `json:"movement_type"`
	ReferenceType  string  `json:"reference_type"`
	ReferenceID    string  `json:"reference_id"`
	Notes          string  `json:"notes"`
	CreatedBy      string  `json:"created_by"`
}

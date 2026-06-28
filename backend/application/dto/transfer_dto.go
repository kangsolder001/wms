package dto

type CreateTransferRequest struct {
	FromLocationID string  `json:"from_location_id" validate:"required"`
	ToLocationID   string  `json:"to_location_id" validate:"required"`
	ItemID         string  `json:"item_id" validate:"required"`
	Quantity       float64 `json:"quantity" validate:"required,gt=0"`
	Notes          string  `json:"notes"`
}

type TransferResponse struct {
	ID             string  `json:"id"`
	TransferNumber string  `json:"transfer_number"`
	FromLocationID string  `json:"from_location_id"`
	ToLocationID   string  `json:"to_location_id"`
	ItemID         string  `json:"item_id"`
	Quantity       float64 `json:"quantity"`
	Status         string  `json:"status"`
	Notes          string  `json:"notes"`
}

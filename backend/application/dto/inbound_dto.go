package dto

import "time"

type CreatePurchaseOrderRequest struct {
	SupplierName string                    `json:"supplier_name" validate:"required"`
	ExpectedDate string                    `json:"expected_date"`
	Notes        string                    `json:"notes"`
	Items        []CreatePOItemRequest     `json:"items" validate:"required,dive"`
}

type CreatePOItemRequest struct {
	ItemID           string  `json:"item_id" validate:"required"`
	ExpectedQuantity float64 `json:"expected_quantity" validate:"required,gt=0"`
	UnitPrice        float64 `json:"unit_price"`
}

type ReceiveGoodsRequest struct {
	GRNNumber string                 `json:"grn_number" validate:"required"`
	Notes     string                 `json:"notes"`
	Items     []ReceiveItemRequest   `json:"items" validate:"required,dive"`
}

type ReceiveItemRequest struct {
	ItemID      string  `json:"item_id" validate:"required"`
	Quantity    float64 `json:"quantity" validate:"required,gt=0"`
	BatchNumber string  `json:"batch_number"`
	LocationID  string  `json:"location_id" validate:"required"`
}

type PurchaseOrderResponse struct {
	ID           string                    `json:"id"`
	PONumber     string                    `json:"po_number"`
	SupplierName string                    `json:"supplier_name"`
	Status       string                    `json:"status"`
	ExpectedDate *time.Time                `json:"expected_date,omitempty"`
	Notes        string                    `json:"notes"`
	CreatedBy    string                    `json:"created_by"`
	CreatedByName string                   `json:"created_by_name,omitempty"`
	CreatedAt    time.Time                 `json:"created_at"`
	Items        []POItemResponse          `json:"items,omitempty"`
}

type POItemResponse struct {
	ID               string  `json:"id"`
	ItemID           string  `json:"item_id"`
	ExpectedQuantity float64 `json:"expected_quantity"`
	ReceivedQuantity float64 `json:"received_quantity"`
	UnitPrice        float64 `json:"unit_price"`
}

type GoodsReceiptResponse struct {
	ID         string                   `json:"id"`
	GRNNumber  string                   `json:"grn_number"`
	POID       string                   `json:"po_id"`
	ReceivedBy string                   `json:"received_by"`
	ReceivedAt time.Time                `json:"received_at"`
	Notes      string                   `json:"notes"`
	Items      []GRNItemResponse        `json:"items,omitempty"`
}

type GRNItemResponse struct {
	ID          string  `json:"id"`
	ItemID      string  `json:"item_id"`
	Quantity    float64 `json:"quantity"`
	BatchNumber string  `json:"batch_number"`
	LocationID  string  `json:"location_id"`
}

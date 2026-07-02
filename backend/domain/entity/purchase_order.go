package entity

import "time"

type PurchaseOrder struct {
	ID           string                 `json:"id"`
	PONumber     string                 `json:"po_number"`
	SupplierName string                 `json:"supplier_name"`
	Status       string                 `json:"status"`
	ExpectedDate *time.Time             `json:"expected_date,omitempty"`
	Notes        string                 `json:"notes"`
	CreatedBy    string                 `json:"created_by"`
	CreatedByName string                `json:"created_by_name,omitempty"`
	CreatedAt    time.Time              `json:"created_at"`
	UpdatedAt    time.Time              `json:"updated_at"`
	Items        []PurchaseOrderItem    `json:"items,omitempty"`
}

type PurchaseOrderItem struct {
	ID                string  `json:"id"`
	POID              string  `json:"po_id"`
	ItemID            string  `json:"item_id"`
	ExpectedQuantity  float64 `json:"expected_quantity"`
	ReceivedQuantity  float64 `json:"received_quantity"`
	UnitPrice         float64 `json:"unit_price"`
}

type GoodsReceipt struct {
	ID        string              `json:"id"`
	GRNNumber string              `json:"grn_number"`
	POID      string              `json:"po_id"`
	ReceivedBy string             `json:"received_by"`
	ReceivedAt time.Time          `json:"received_at"`
	Notes     string              `json:"notes"`
	Items     []GoodsReceiptItem  `json:"items,omitempty"`
}

type GoodsReceiptItem struct {
	ID         string  `json:"id"`
	GRNID      string  `json:"grn_id"`
	ItemID     string  `json:"item_id"`
	Quantity   float64 `json:"quantity"`
	BatchNumber string `json:"batch_number"`
	LocationID string  `json:"location_id"`
}

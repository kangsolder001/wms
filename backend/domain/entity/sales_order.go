package entity

import "time"

type SalesOrder struct {
	ID           string              `json:"id"`
	SONumber     string              `json:"so_number"`
	CustomerName string              `json:"customer_name"`
	Status       string              `json:"status"`
	Notes        string              `json:"notes"`
	CreatedBy    string              `json:"created_by"`
	CreatedAt    time.Time           `json:"created_at"`
	Items        []SalesOrderItem    `json:"items,omitempty"`
}

type SalesOrderItem struct {
	ID               string  `json:"id"`
	SOID             string  `json:"so_id"`
	ItemID           string  `json:"item_id"`
	OrderedQuantity  float64 `json:"ordered_quantity"`
	PickedQuantity   float64 `json:"picked_quantity"`
}

type PickList struct {
	ID        string     `json:"id"`
	SOID      string     `json:"so_id"`
	Status    string     `json:"status"`
	PickedBy  *string    `json:"picked_by,omitempty"`
	PickedAt  *time.Time `json:"picked_at,omitempty"`
	CreatedAt time.Time  `json:"created_at"`
}

type Shipment struct {
	ID             string     `json:"id"`
	ShipmentNumber string     `json:"shipment_number"`
	SOID           string     `json:"so_id"`
	Carrier        string     `json:"carrier"`
	TrackingNumber string     `json:"tracking_number"`
	Status         string     `json:"status"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
	CreatedBy      string     `json:"created_by"`
}

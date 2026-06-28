package dto

import "time"

type CreateSalesOrderRequest struct {
	CustomerName string                  `json:"customer_name" validate:"required"`
	Notes        string                  `json:"notes"`
	Items        []CreateSOItemRequest   `json:"items" validate:"required,dive"`
}

type CreateSOItemRequest struct {
	ItemID          string  `json:"item_id" validate:"required"`
	OrderedQuantity float64 `json:"ordered_quantity" validate:"required,gt=0"`
}

type PickListRequest struct {
	Items []PickItemRequest `json:"items" validate:"required,dive"`
}

type PickItemRequest struct {
	ItemID   string  `json:"item_id" validate:"required"`
	Quantity float64 `json:"quantity" validate:"required,gt=0"`
}

type ShipRequest struct {
	Carrier        string `json:"carrier" validate:"required"`
	TrackingNumber string `json:"tracking_number"`
}

type SalesOrderResponse struct {
	ID           string              `json:"id"`
	SONumber     string              `json:"so_number"`
	CustomerName string              `json:"customer_name"`
	Status       string              `json:"status"`
	Notes        string              `json:"notes"`
	CreatedAt    time.Time           `json:"created_at"`
	Items        []SOItemResponse    `json:"items,omitempty"`
}

type SOItemResponse struct {
	ID              string  `json:"id"`
	ItemID          string  `json:"item_id"`
	OrderedQuantity float64 `json:"ordered_quantity"`
	PickedQuantity  float64 `json:"picked_quantity"`
}

type ShipmentResponse struct {
	ID             string     `json:"id"`
	ShipmentNumber string     `json:"shipment_number"`
	SOID           string     `json:"so_id"`
	Carrier        string     `json:"carrier"`
	TrackingNumber string     `json:"tracking_number"`
	Status         string     `json:"status"`
	ShippedAt      *time.Time `json:"shipped_at,omitempty"`
}

package dto

type CreateItemRequest struct {
	CategoryID    string  `json:"category_id" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description"`
	UnitOfMeasure string  `json:"unit_of_measure" validate:"required"`
	Weight        float64 `json:"weight"`
	Length        float64 `json:"length"`
	Width         float64 `json:"width"`
	Height        float64 `json:"height"`
}

type UpdateItemRequest struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	CategoryID    *string  `json:"category_id"`
	UnitOfMeasure *string  `json:"unit_of_measure"`
	Weight        *float64 `json:"weight"`
	Length        *float64 `json:"length"`
	Width         *float64 `json:"width"`
	Height        *float64 `json:"height"`
	IsActive      *bool    `json:"is_active"`
}

type ItemResponse struct {
	ID            string  `json:"id"`
	SKU           string  `json:"sku"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	CategoryID    string  `json:"category_id"`
	Barcode       string  `json:"barcode"`
	UnitOfMeasure string  `json:"unit_of_measure"`
	Weight        float64 `json:"weight"`
	Length        float64 `json:"length"`
	Width         float64 `json:"width"`
	Height        float64 `json:"height"`
	IsActive      bool    `json:"is_active"`
}

type GenerateSKURequest struct {
	CategoryID string `json:"category_id" validate:"required"`
}

type GenerateSKUResponse struct {
	SKU string `json:"sku"`
}

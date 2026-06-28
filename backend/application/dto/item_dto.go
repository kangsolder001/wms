package dto

type CreateItemRequest struct {
	SKU           string  `json:"sku" validate:"required"`
	Name          string  `json:"name" validate:"required"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	UnitOfMeasure string  `json:"unit_of_measure" validate:"required"`
	Weight        float64 `json:"weight"`
}

type UpdateItemRequest struct {
	Name          *string  `json:"name"`
	Description   *string  `json:"description"`
	Category      *string  `json:"category"`
	UnitOfMeasure *string  `json:"unit_of_measure"`
	Weight        *float64 `json:"weight"`
	IsActive      *bool    `json:"is_active"`
}

type ItemResponse struct {
	ID            string  `json:"id"`
	SKU           string  `json:"sku"`
	Name          string  `json:"name"`
	Description   string  `json:"description"`
	Category      string  `json:"category"`
	UnitOfMeasure string  `json:"unit_of_measure"`
	Weight        float64 `json:"weight"`
	IsActive      bool    `json:"is_active"`
}

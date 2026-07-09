package dto

type CreateCategoryRequest struct {
	Name         string `json:"name" validate:"required"`
	Abbreviation string `json:"abbreviation" validate:"required"`
	Description  string `json:"description"`
}

type UpdateCategoryRequest struct {
	Name         *string `json:"name"`
	Abbreviation *string `json:"abbreviation"`
	Description  *string `json:"description"`
	IsActive     *bool   `json:"is_active"`
}

type CategoryResponse struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Abbreviation string `json:"abbreviation"`
	Description  string `json:"description"`
	IsActive     bool   `json:"is_active"`
}

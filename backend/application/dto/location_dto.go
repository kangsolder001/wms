package dto

type CreateLocationRequest struct {
	Code     string  `json:"code" validate:"required"`
	Name     string  `json:"name" validate:"required"`
	Zone     string  `json:"zone"`
	Aisle    string  `json:"aisle"`
	Rack     string  `json:"rack"`
	Level    string  `json:"level"`
	Bin      string  `json:"bin"`
	Type     string  `json:"type" validate:"required,oneof=storage receiving shipping staging"`
	Capacity float64 `json:"capacity"`
}

type UpdateLocationRequest struct {
	Name     *string  `json:"name"`
	Zone     *string  `json:"zone"`
	Aisle    *string  `json:"aisle"`
	Rack     *string  `json:"rack"`
	Level    *string  `json:"level"`
	Bin      *string  `json:"bin"`
	Type     *string  `json:"type"`
	Capacity *float64 `json:"capacity"`
	IsActive *bool    `json:"is_active"`
}

type LocationResponse struct {
	ID        string  `json:"id"`
	Code      string  `json:"code"`
	Name      string  `json:"name"`
	Zone      string  `json:"zone"`
	Aisle     string  `json:"aisle"`
	Rack      string  `json:"rack"`
	Level     string  `json:"level"`
	Bin       string  `json:"bin"`
	Type      string  `json:"type"`
	Capacity  float64 `json:"capacity"`
	IsActive  bool    `json:"is_active"`
}

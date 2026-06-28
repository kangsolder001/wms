package entity

import "time"

type Location struct {
	ID        string    `json:"id"`
	Code      string    `json:"code"`
	Name      string    `json:"name"`
	Zone      string    `json:"zone"`
	Aisle     string    `json:"aisle"`
	Rack      string    `json:"rack"`
	Level     string    `json:"level"`
	Bin       string    `json:"bin"`
	Type      string    `json:"type"`
	Capacity  float64   `json:"capacity"`
	IsActive  bool      `json:"is_active"`
	CreatedAt time.Time `json:"created_at"`
}

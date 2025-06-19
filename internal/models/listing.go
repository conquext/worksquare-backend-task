package models

import (
	"strings"
)

// Listing represents a housing listing
type Listing struct {
	ID         int      `json:"id"`
	Title      string   `json:"title"`
	Price      string   `json:"price"`
	Bedrooms   int      `json:"bedrooms"`
	Bathrooms  int      `json:"bathrooms"`
	Location   string   `json:"location"`
	Status     []string `json:"status"`
	Image      string   `json:"image"`
}

type PaginationQuery struct {
	Page  int `json:"page" query:"page" validate:"min=1"`
	Limit int `json:"limit" query:"limit" validate:"min=1,max=100"`
}

func (p *PaginationQuery) SetDefaults() {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 10
	}
	if p.Limit > 100 {
		p.Limit = 100
	}
}
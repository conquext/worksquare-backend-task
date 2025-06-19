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

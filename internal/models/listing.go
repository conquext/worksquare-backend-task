package models

import (
	"strconv"
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

// GetPropertyType returns the property type from the status array
func (l *Listing) GetPropertyType() string {
	if len(l.Status) > 0 {
		return l.Status[0]
	}
	return ""
}

// GetListingType returns the listing type (For Rent, For Lease, etc.) from the status array
func (l *Listing) GetListingType() string {
	if len(l.Status) > 1 {
		return l.Status[1]
	}
	return ""
}

// GetPriceNumeric returns the numeric price value, removing currency symbols and formatting
func (l *Listing) GetPriceNumeric() float64 {
	// Remove currency symbol and commas
	priceStr := strings.ReplaceAll(l.Price, "₦", "")
	priceStr = strings.ReplaceAll(priceStr, ",", "")
	
	// Handle special cases like "/ week", "/ night"
	if strings.Contains(priceStr, "/") {
		priceStr = strings.Split(priceStr, "/")[0]
	}
	
	priceStr = strings.TrimSpace(priceStr)
	
	price, err := strconv.ParseFloat(priceStr, 64)
	if err != nil {
		return 0
	}
	return price
}

// GetCity extracts the city from the location string
func (l *Listing) GetCity() string {
	parts := strings.Split(l.Location, ",")
	if len(parts) > 0 {
		return strings.TrimSpace(parts[len(parts)-1])
	}
	return l.Location
}

// GetArea extracts the area from the location string
func (l *Listing) GetArea() string {
	parts := strings.Split(l.Location, ",")
	if len(parts) > 1 {
		return strings.TrimSpace(parts[0])
	}
	return l.Location
}

// ListingFilter represents filtering options for listings
type ListingFilter struct {
	Location     string `json:"location" query:"location"`
	PropertyType string `json:"property_type" query:"property_type"`
	MinPrice     *int   `json:"min_price" query:"min_price"`
	MaxPrice     *int   `json:"max_price" query:"max_price"`
	MinBedrooms  *int   `json:"min_bedrooms" query:"min_bedrooms"`
	MaxBedrooms  *int   `json:"max_bedrooms" query:"max_bedrooms"`
	MinBathrooms *int   `json:"min_bathrooms" query:"min_bathrooms"`
	MaxBathrooms *int   `json:"max_bathrooms" query:"max_bathrooms"`
	City         string `json:"city" query:"city"`
}

// PaginationQuery represents pagination parameters
type PaginationQuery struct {
	Page  int `json:"page" query:"page" validate:"min=1"`
	Limit int `json:"limit" query:"limit" validate:"min=1,max=100"`
}

// SetDefaults sets default values for pagination
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
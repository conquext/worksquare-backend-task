package services

import (
	"fmt"

	"housing-api/internal/models"
	"housing-api/internal/repositories"
	"housing-api/pkg/pagination"
)

// ListingService handles business logic for listings
type ListingService struct {
	repo *repositories.ListingRepository
}

// NewListingService creates a new listing service
func NewListingService() (*ListingService, error) {
	repo, err := repositories.NewListingRepository()
	if err != nil {
		return nil, fmt.Errorf("failed to create listing repository: %w", err)
	}

	return &ListingService{
		repo: repo,
	}, nil
}

// GetListings returns paginated listings with optional filtering
func (s *ListingService) GetListings(filter models.ListingFilter, paginationQuery models.PaginationQuery) (*models.PaginatedResponse, error) {
	paginationQuery.SetDefaults()

	// Get paginated listings
	listings, total, err := s.repo.GetPaginated(filter, paginationQuery.Page, paginationQuery.Limit)
	if err != nil {
		return nil, fmt.Errorf("failed to get listings: %w", err)
	}

	// Convert to interface slice
	items := make([]interface{}, len(listings))
	for i, listing := range listings {
		items[i] = listing
	}

	// Calculate pagination metadata
	meta := pagination.CalculateMetadata(paginationQuery.Page, paginationQuery.Limit, total)

	return &models.PaginatedResponse{
		Items: items,
		Meta:  meta,
	}, nil
}

// GetListingByID returns a single listing by ID
func (s *ListingService) GetListingByID(id int) (*models.Listing, error) {
	listing, err := s.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("listing not found: %w", err)
	}

	return listing, nil
}

// GetFiltersMetadata returns metadata for filtering (unique locations, property types)
func (s *ListingService) GetFiltersMetadata() (map[string]interface{}, error) {
	locations := s.repo.GetUniqueLocations()
	propertyTypes := s.repo.GetUniquePropertyTypes()

	metadata := map[string]interface{}{
		"locations":      locations,
		"property_types": propertyTypes,
		"filters": map[string]interface{}{
			"location": map[string]interface{}{
				"type":        "string",
				"description": "Filter by location (partial match)",
				"example":     "Lagos",
			},
			"property_type": map[string]interface{}{
				"type":        "string",
				"description": "Filter by property type (exact match)",
				"options":     propertyTypes,
				"example":     "House",
			},
			"city": map[string]interface{}{
				"type":        "string",
				"description": "Filter by city (partial match)",
				"example":     "Lagos",
			},
			"min_price": map[string]interface{}{
				"type":        "integer",
				"description": "Minimum price filter",
				"example":     1000000,
			},
			"max_price": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum price filter",
				"example":     5000000,
			},
			"min_bedrooms": map[string]interface{}{
				"type":        "integer",
				"description": "Minimum number of bedrooms",
				"example":     2,
			},
			"max_bedrooms": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum number of bedrooms",
				"example":     5,
			},
			"min_bathrooms": map[string]interface{}{
				"type":        "integer",
				"description": "Minimum number of bathrooms",
				"example":     1,
			},
			"max_bathrooms": map[string]interface{}{
				"type":        "integer",
				"description": "Maximum number of bathrooms",
				"example":     4,
			},
		},
		"pagination": map[string]interface{}{
			"page": map[string]interface{}{
				"type":        "integer",
				"description": "Page number (starts from 1)",
				"default":     1,
				"minimum":     1,
			},
			"limit": map[string]interface{}{
				"type":        "integer",
				"description": "Number of items per page",
				"default":     10,
				"minimum":     1,
				"maximum":     100,
			},
		},
	}

	return metadata, nil
}

// SearchListings provides advanced search functionality
func (s *ListingService) SearchListings(query string, filter models.ListingFilter, paginationQuery models.PaginationQuery) (*models.PaginatedResponse, error) {
	if query != "" && filter.Location == "" {
		filter.Location = query
	}

	return s.GetListings(filter, paginationQuery)
}

// GetListingStats returns statistics about listings
func (s *ListingService) GetListingStats() (map[string]interface{}, error) {
	allListings, err := s.repo.GetAll(models.ListingFilter{})
	if err != nil {
		return nil, fmt.Errorf("failed to get listings for stats: %w", err)
	}

	totalListings := len(allListings)
	propertyTypes := make(map[string]int)
	cities := make(map[string]int)
	priceRanges := map[string]int{
		"under_1m":    0,
		"1m_to_2m":    0,
		"2m_to_3m":    0,
		"3m_to_5m":    0,
		"above_5m":    0,
	}

	var totalPrice, minPrice, maxPrice float64
	minPrice = float64(^uint(0) >> 1) // Max int value

	for i, listing := range allListings {
		// Count property types
		propertyType := listing.GetPropertyType()
		propertyTypes[propertyType]++

		// Count cities
		city := listing.GetCity()
		cities[city]++

		// Price statistics
		price := listing.GetPriceNumeric()
		if price > 0 {
			totalPrice += price
			if price < minPrice {
				minPrice = price
			}
			if price > maxPrice {
				maxPrice = price
			}

			// Price ranges
			if price < 1000000 {
				priceRanges["under_1m"]++
			} else if price < 2000000 {
				priceRanges["1m_to_2m"]++
			} else if price < 3000000 {
				priceRanges["2m_to_3m"]++
			} else if price < 5000000 {
				priceRanges["3m_to_5m"]++
			} else {
				priceRanges["above_5m"]++
			}
		}

		// Handle first iteration for minPrice
		if i == 0 && minPrice == float64(^uint(0)>>1) {
			minPrice = price
		}
	}

	avgPrice := float64(0)
	if totalListings > 0 {
		avgPrice = totalPrice / float64(totalListings)
	}

	stats := map[string]interface{}{
		"total_listings":  totalListings,
		"property_types":  propertyTypes,
		"cities":          cities,
		"price_ranges":    priceRanges,
		"price_stats": map[string]interface{}{
			"average": avgPrice,
			"minimum": minPrice,
			"maximum": maxPrice,
		},
	}

	return stats, nil
}
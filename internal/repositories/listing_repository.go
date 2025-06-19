package repositories

import (
	"encoding/json"
	"fmt"
	"os"
	"sort"
	"strings"

	"housing-api/internal/models"
	"housing-api/internal/utils"
)

// ListingRepository handles listing data operations
type ListingRepository struct {
	listings []models.Listing
	filePath string
}

func NewListingRepository() (*ListingRepository, error) {
	filePath := utils.GetDataFilePath("listings.json")

	// Check if file exists
	if _, err := os.Stat(filePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("listings file does not exist")
	}
	repo := &ListingRepository{
		filePath: filePath,
	}

	if err := repo.loadListings(); err != nil {
		return nil, fmt.Errorf("failed to load listing data: %w", err)
	}

	return repo, nil
}

// loadListings loads listings from JSON file
func (r *ListingRepository) loadListings() error {
	file, err := os.ReadFile(r.filePath)
	if err != nil {
		return fmt.Errorf("failed to read listings file: %w", err)
	}

	if err := json.Unmarshal(file, &r.listings); err != nil {
		return fmt.Errorf("failed to unmarshal listings: %w", err)
	}

	return nil
}

// ReloadListings reloads listings from JSON file (useful for updates)
func (r *ListingRepository) ReloadListings() error {
	return r.loadListings()
}

// GetAll returns all listings with optional filtering
func (r *ListingRepository) GetAll(filter models.ListingFilter) ([]models.Listing, error) {
	var filtered []models.Listing

	for _, listing := range r.listings {
		if r.matchesFilter(listing, filter) {
			filtered = append(filtered, listing)
		}
	}

	return filtered, nil
}

func (r *ListingRepository) GetByID(id int) (*models.Listing, error) {
	for _, listing := range r.listings {
		if listing.ID == id {
			return &listing, nil
		}
	}
	return nil, fmt.Errorf("listing with ID %d not found", id)
}

// GetPaginated returns paginated listings with sorting support
func (r *ListingRepository) GetPaginated(filter models.ListingFilter, page, limit int) ([]models.Listing, int64, error) {
	// Get all filtered listings
	filtered, err := r.GetAll(filter)
	if err != nil {
		return nil, 0, err
	}

	total := int64(len(filtered))

	// Sort by ID for consistent pagination
	sort.Slice(filtered, func(i, j int) bool {
		return filtered[i].ID < filtered[j].ID
	})

	// Calculate pagination
	offset := (page - 1) * limit
	if offset >= len(filtered) {
		return []models.Listing{}, total, nil
	}

	end := offset + limit
	if end > len(filtered) {
		end = len(filtered)
	}

	return filtered[offset:end], total, nil
}

// GetUniqueLocations returns all unique locations (cities)
func (r *ListingRepository) GetUniqueLocations() []string {
	locationMap := make(map[string]bool)
	var locations []string

	for _, listing := range r.listings {
		city := listing.GetCity()
		if city != "" && !locationMap[city] {
			locationMap[city] = true
			locations = append(locations, city)
		}
	}

	// Sort locations alphabetically
	sort.Strings(locations)
	return locations
}

// GetUniquePropertyTypes returns all unique property types
func (r *ListingRepository) GetUniquePropertyTypes() []string {
	typeMap := make(map[string]bool)
	var types []string

	for _, listing := range r.listings {
		propertyType := listing.GetPropertyType()
		if propertyType != "" && !typeMap[propertyType] {
			typeMap[propertyType] = true
			types = append(types, propertyType)
		}
	}

	// Sort property types alphabetically
	sort.Strings(types)
	return types
}

// GetPriceRange returns the minimum and maximum prices in the dataset
func (r *ListingRepository) GetPriceRange() (float64, float64) {
	if len(r.listings) == 0 {
		return 0, 0
	}

	minPrice := r.listings[0].GetPriceNumeric()
	maxPrice := minPrice

	for _, listing := range r.listings {
		price := listing.GetPriceNumeric()
		if price > 0 { // Only consider valid prices
			if price < minPrice {
				minPrice = price
			}
			if price > maxPrice {
				maxPrice = price
			}
		}
	}

	return minPrice, maxPrice
}

// GetBedroomRange returns the minimum and maximum number of bedrooms
func (r *ListingRepository) GetBedroomRange() (int, int) {
	if len(r.listings) == 0 {
		return 0, 0
	}

	minBedrooms := r.listings[0].Bedrooms
	maxBedrooms := minBedrooms

	for _, listing := range r.listings {
		if listing.Bedrooms < minBedrooms {
			minBedrooms = listing.Bedrooms
		}
		if listing.Bedrooms > maxBedrooms {
			maxBedrooms = listing.Bedrooms
		}
	}

	return minBedrooms, maxBedrooms
}

// GetBathroomRange returns the minimum and maximum number of bathrooms
func (r *ListingRepository) GetBathroomRange() (int, int) {
	if len(r.listings) == 0 {
		return 0, 0
	}

	minBathrooms := r.listings[0].Bathrooms
	maxBathrooms := minBathrooms

	for _, listing := range r.listings {
		if listing.Bathrooms < minBathrooms {
			minBathrooms = listing.Bathrooms
		}
		if listing.Bathrooms > maxBathrooms {
			maxBathrooms = listing.Bathrooms
		}
	}

	return minBathrooms, maxBathrooms
}

// GetListingsByPropertyType returns listings grouped by property type
func (r *ListingRepository) GetListingsByPropertyType() map[string][]models.Listing {
	grouped := make(map[string][]models.Listing)

	for _, listing := range r.listings {
		propertyType := listing.GetPropertyType()
		if propertyType != "" {
			grouped[propertyType] = append(grouped[propertyType], listing)
		}
	}

	return grouped
}

// GetListingsByCity returns listings grouped by city
func (r *ListingRepository) GetListingsByCity() map[string][]models.Listing {
	grouped := make(map[string][]models.Listing)

	for _, listing := range r.listings {
		city := listing.GetCity()
		if city != "" {
			grouped[city] = append(grouped[city], listing)
		}
	}

	return grouped
}

// GetTotalCount returns the total number of listings
func (r *ListingRepository) GetTotalCount() int {
	return len(r.listings)
}

// SearchListings performs a text-based search across multiple fields
func (r *ListingRepository) SearchListings(query string) []models.Listing {
	if query == "" {
		return r.listings
	}

	query = strings.ToLower(query)
	var results []models.Listing

	for _, listing := range r.listings {
		// Search in title, location, and property type
		if r.containsQuery(listing.Title, query) ||
			r.containsQuery(listing.Location, query) ||
			r.containsQuery(listing.GetPropertyType(), query) {
			results = append(results, listing)
		}
	}

	return results
}

// containsQuery checks if a field contains the search query (case-insensitive)
func (r *ListingRepository) containsQuery(field, query string) bool {
	return strings.Contains(strings.ToLower(field), query)
}

// matchesFilter checks if a listing matches the given filter criteria
func (r *ListingRepository) matchesFilter(listing models.Listing, filter models.ListingFilter) bool {
	// Location filter (case-insensitive, partial match in full location string)
	if filter.Location != "" {
		if !strings.Contains(strings.ToLower(listing.Location), strings.ToLower(filter.Location)) {
			return false
		}
	}

	// City filter (case-insensitive, partial match in city name)
	if filter.City != "" {
		if !strings.Contains(strings.ToLower(listing.GetCity()), strings.ToLower(filter.City)) {
			return false
		}
	}

	// Property type filter (case-insensitive, exact match)
	if filter.PropertyType != "" {
		if !strings.EqualFold(listing.GetPropertyType(), filter.PropertyType) {
			return false
		}
	}

	// Price range filter
	price := listing.GetPriceNumeric()
	if price > 0 { // Only apply price filters to listings with valid prices
		if filter.MinPrice != nil && price < float64(*filter.MinPrice) {
			return false
		}
		if filter.MaxPrice != nil && price > float64(*filter.MaxPrice) {
			return false
		}
	}

	// Bedrooms range filter
	if filter.MinBedrooms != nil && listing.Bedrooms < *filter.MinBedrooms {
		return false
	}
	if filter.MaxBedrooms != nil && listing.Bedrooms > *filter.MaxBedrooms {
		return false
	}

	// Bathrooms range filter
	if filter.MinBathrooms != nil && listing.Bathrooms < *filter.MinBathrooms {
		return false
	}
	if filter.MaxBathrooms != nil && listing.Bathrooms > *filter.MaxBathrooms {
		return false
	}

	return true
}

// GetSimilarListings returns listings similar to the given listing
func (r *ListingRepository) GetSimilarListings(targetListing models.Listing, limit int) []models.Listing {
	var similar []models.Listing

	for _, listing := range r.listings {
		// Skip the same listing
		if listing.ID == targetListing.ID {
			continue
		}

		// Consider listings similar if they have:
		// 1. Same city
		// 2. Same property type
		// 3. Similar number of bedrooms (within 1)
		// 4. Similar price range (within 20%)
		if r.isSimilar(listing, targetListing) {
			similar = append(similar, listing)
		}

		// Stop when we have enough similar listings
		if len(similar) >= limit {
			break
		}
	}

	return similar
}

// isSimilar checks if two listings are similar based on key criteria
func (r *ListingRepository) isSimilar(listing1, listing2 models.Listing) bool {
	// Same city
	if !strings.EqualFold(listing1.GetCity(), listing2.GetCity()) {
		return false
	}

	// Same property type
	if !strings.EqualFold(listing1.GetPropertyType(), listing2.GetPropertyType()) {
		return false
	}

	// Similar bedrooms (within 1)
	bedroomDiff := listing1.Bedrooms - listing2.Bedrooms
	if bedroomDiff < -1 || bedroomDiff > 1 {
		return false
	}

	// Similar price range (within 20%)
	price1 := listing1.GetPriceNumeric()
	price2 := listing2.GetPriceNumeric()
	
	if price1 > 0 && price2 > 0 {
		priceDiff := (price1 - price2) / price2
		if priceDiff < -0.2 || priceDiff > 0.2 {
			return false
		}
	}

	return true
}
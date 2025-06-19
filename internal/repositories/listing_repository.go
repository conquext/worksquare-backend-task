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
	
	repo.loadListings()

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

// GetPaginated returns paginated listings with sorting support
func (r *ListingRepository) GetPaginated(page, limit int) ([]models.Listing, int64, error) {
	all_listings := r.listings
	
	total := int64(len(all_listings))

	// Sort by ID for consistent pagination
	sort.Slice(all_listings, func(i, j int) bool {
		return all_listings[i].ID < all_listings[j].ID
	})

	// Calculate pagination
	offset := (page - 1) * limit
	if offset >= len(all_listings) {
		return []models.Listing{}, total, nil
	}

	end := offset + limit
	if end > len(all_listings) {
		end = len(all_listings)
	}

	return all_listings[offset:end], total, nil
}
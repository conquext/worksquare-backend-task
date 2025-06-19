package unit

import (
	"testing"

	"housing-api/internal/models"
	"housing-api/internal/services"

	"github.com/stretchr/testify/assert"
)

func TestListingService_GetListings(t *testing.T) {
	// Initialize service
	service, err := services.NewListingService()
	assert.NoError(t, err)
	assert.NotNil(t, service)

	// Test getting listings with default pagination
	paginationQuery := models.PaginationQuery{
		Page:  1,
		Limit: 10,
	}
	paginationQuery.SetDefaults()

	filter := models.ListingFilter{}

	result, err := service.GetListings(filter, paginationQuery)
	assert.NoError(t, err)
	assert.NotNil(t, result)
	assert.LessOrEqual(t, len(result.Items), 10)
}

func TestListingService_GetListingByID(t *testing.T) {
	service, err := services.NewListingService()
	assert.NoError(t, err)

	// Test valid ID
	listing, err := service.GetListingByID(1)
	assert.NoError(t, err)
	assert.NotNil(t, listing)
	assert.Equal(t, 1, listing.ID)

	// Test invalid ID
	listing, err = service.GetListingByID(9999)
	assert.Error(t, err)
	assert.Nil(t, listing)
}

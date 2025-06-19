package controllers

import (
	"strconv"

	"housing-api/internal/models"
	"housing-api/internal/services"
	"housing-api/pkg/response"

	"github.com/gofiber/fiber/v2"
)

// ListingController handles listing-related HTTP requests
type ListingController struct {
	listingService *services.ListingService
}

// NewListingController creates a new listing controller
func NewListingController() (*ListingController, error) {
	listingService, err := services.NewListingService()
	if err != nil {
		return nil, err
	}

	return &ListingController{
		listingService: listingService,
	}, nil
}

// GetListings godoc
// @Summary Get all listings
// @Description Get all housing listings with optional filtering and pagination
// @Tags listings
// @Accept json
// @Produce json
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param location query string false "Filter by location"
// @Param property_type query string false "Filter by property type"
// @Param city query string false "Filter by city"
// @Param min_price query int false "Minimum price"
// @Param max_price query int false "Maximum price"
// @Param min_bedrooms query int false "Minimum bedrooms"
// @Param max_bedrooms query int false "Maximum bedrooms"
// @Param min_bathrooms query int false "Minimum bathrooms"
// @Param max_bathrooms query int false "Maximum bathrooms"
// @Success 200 {object} models.APIResponse{data=models.PaginatedResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /listings [get]
func (c *ListingController) GetListings(ctx *fiber.Ctx) error {
	// Parse pagination query
	var paginationQuery models.PaginationQuery
	if err := ctx.QueryParser(&paginationQuery); err != nil {
		return response.BadRequest(ctx, "Invalid pagination parameters", err)
	}

	// Parse filter query
	var filter models.ListingFilter
	if err := ctx.QueryParser(&filter); err != nil {
		return response.BadRequest(ctx, "Invalid filter parameters", err)
	}

	// Get listings
	result, err := c.listingService.GetListings(filter, paginationQuery)
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get listings", err)
	}

	return response.Success(ctx, "Listings retrieved successfully", result)
}

// GetListingByID godoc
// @Summary Get listing by ID
// @Description Get a single housing listing by its ID
// @Tags listings
// @Accept json
// @Produce json
// @Param id path int true "Listing ID"
// @Success 200 {object} models.APIResponse{data=models.Listing}
// @Failure 400 {object} models.APIResponse
// @Failure 404 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /listings/{id} [get]
func (c *ListingController) GetListingByID(ctx *fiber.Ctx) error {
	// Parse ID parameter
	idParam := ctx.Params("id")
	id, err := strconv.Atoi(idParam)
	if err != nil {
		return response.BadRequest(ctx, "Invalid listing ID", err)
	}

	// Get listing
	listing, err := c.listingService.GetListingByID(id)
	if err != nil {
		return response.NotFound(ctx, "Listing not found", err)
	}

	return response.Success(ctx, "Listing retrieved successfully", listing)
}

// SearchListings godoc
// @Summary Search listings
// @Description Search housing listings with query string
// @Tags listings
// @Accept json
// @Produce json
// @Param q query string true "Search query"
// @Param page query int false "Page number" default(1)
// @Param limit query int false "Items per page" default(10)
// @Param location query string false "Filter by location"
// @Param property_type query string false "Filter by property type"
// @Param city query string false "Filter by city"
// @Success 200 {object} models.APIResponse{data=models.PaginatedResponse}
// @Failure 400 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /listings/search [get]
func (c *ListingController) SearchListings(ctx *fiber.Ctx) error {
	// Get search query
	query := ctx.Query("q")
	if query == "" {
		return response.BadRequest(ctx, "Search query is required", nil)
	}

	// Parse pagination query
	var paginationQuery models.PaginationQuery
	if err := ctx.QueryParser(&paginationQuery); err != nil {
		return response.BadRequest(ctx, "Invalid pagination parameters", err)
	}

	// Parse filter query
	var filter models.ListingFilter
	if err := ctx.QueryParser(&filter); err != nil {
		return response.BadRequest(ctx, "Invalid filter parameters", err)
	}

	// Search listings
	result, err := c.listingService.SearchListings(query, filter, paginationQuery)
	if err != nil {
		return response.InternalServerError(ctx, "Failed to search listings", err)
	}

	return response.Success(ctx, "Search completed successfully", result)
}

// GetFiltersMetadata godoc
// @Summary Get filter metadata
// @Description Get available filter options and metadata
// @Tags listings
// @Accept json
// @Produce json
// @Success 200 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /listings/filters [get]
func (c *ListingController) GetFiltersMetadata(ctx *fiber.Ctx) error {
	metadata, err := c.listingService.GetFiltersMetadata()
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get filter metadata", err)
	}

	return response.Success(ctx, "Filter metadata retrieved successfully", metadata)
}

// GetListingStats godoc
// @Summary Get listing statistics
// @Description Get statistics about listings (protected route)
// @Tags listings
// @Accept json
// @Produce json
// @Security BearerAuth
// @Success 200 {object} models.APIResponse
// @Failure 401 {object} models.APIResponse
// @Failure 500 {object} models.APIResponse
// @Router /listings/stats [get]
func (c *ListingController) GetListingStats(ctx *fiber.Ctx) error {
	stats, err := c.listingService.GetListingStats()
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get listing statistics", err)
	}

	return response.Success(ctx, "Listing statistics retrieved successfully", stats)
}
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

func (c *ListingController) GetFiltersMetadata(ctx *fiber.Ctx) error {
	metadata, err := c.listingService.GetFiltersMetadata()
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get filter metadata", err)
	}

	return response.Success(ctx, "Filter metadata retrieved successfully", metadata)
}

func (c *ListingController) GetListingStats(ctx *fiber.Ctx) error {
	stats, err := c.listingService.GetListingStats()
	if err != nil {
		return response.InternalServerError(ctx, "Failed to get listing statistics", err)
	}

	return response.Success(ctx, "Listing statistics retrieved successfully", stats)
}
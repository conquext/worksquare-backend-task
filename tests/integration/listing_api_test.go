package integration

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	"housing-api/api/routes"
	"housing-api/internal/config"
	"housing-api/internal/models"

	"github.com/gofiber/fiber/v2"
	"github.com/stretchr/testify/assert"
)

func setupTestApp() *fiber.App {
	app := fiber.New()
	cfg, _ := config.Load()
	routes.Setup(app, cfg)
	return app
}

func TestGetListings(t *testing.T) {
	app := setupTestApp()

	req := httptest.NewRequest("GET", "/api/v1/listings", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	var response models.APIResponse
	err = json.NewDecoder(resp.Body).Decode(&response)
	assert.NoError(t, err)
	assert.True(t, response.Success)
}

func TestGetListingByID(t *testing.T) {
	app := setupTestApp()

	// Test valid ID
	req := httptest.NewRequest("GET", "/api/v1/listings/1", nil)
	resp, err := app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	// Test invalid ID
	req = httptest.NewRequest("GET", "/api/v1/listings/9999", nil)
	resp, err = app.Test(req)

	assert.NoError(t, err)
	assert.Equal(t, http.StatusNotFound, resp.StatusCode)
}
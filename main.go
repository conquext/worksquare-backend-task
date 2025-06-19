package main

import (
	"housing-api/cmd/server"
	"log"
)

// @title Worksquare Housing Listings API
// @version 1.0
// @description RESTful API for housing listings with JWT authentication, pagination, and filtering
// @termsOfService http://swagger.io/terms/

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:3000
// @BasePath /api/v1

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.

func main() {
	if err := server.Run(); err != nil {
		log.Fatal("Failed to start server:", err)
	}
}
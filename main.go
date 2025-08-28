package main

import (
	initializers "crackers/d2delight.com/initializers"
	"crackers/d2delight.com/routes"
	"crackers/d2delight.com/validator"

	"github.com/gin-gonic/gin"
)

func init() {
	validator.CustomValidator()
	initializers.LoadEnvVariables()
	initializers.CreateDatabaseIfNotExists()
	initializers.DbConnect()
}

// @title My API
// @version 1.0
// @description This is my application API documentation
// @termsOfService http://example.com/terms/

// @contact.name API Support
// @contact.url http://www.example.com/support
// @contact.email support@example.com

// @license.name MIT
// @license.url https://opensource.org/licenses/MIT

// @host localhost:8000
// @BasePath /api

func main() {
	router := gin.Default()
	routes.RegisterRoutes(router)
	router.Run()
}

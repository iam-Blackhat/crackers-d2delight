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
func main() {
	router := gin.Default()
	routes.RegisterRoutes(router)
	router.Run()
}

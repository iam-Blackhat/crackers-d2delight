package routes

import (
	"crackers/d2delight.com/controllers"
	"crackers/d2delight.com/middleware"

	_ "crackers/d2delight.com/docs" // swagger docs

	swaggerFiles "github.com/swaggo/files" // swagger embed files
	ginSwagger "github.com/swaggo/gin-swagger"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {

	// Swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// Health check (keep outside /api for simplicity)
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// All APIs under /api
	api := r.Group("/api")
	{
		// Auth routes
		auth := api.Group("/auth")
		{
			auth.POST("/login", controllers.Login)
			auth.POST("/register", controllers.Create)
		}

		// User routes (protected)
		users := api.Group("/users", middleware.AuthRequired())
		{
			users.GET("/", controllers.GetUsers)
			users.GET("/:id", controllers.GetUserByID)
			users.PUT("/:id", controllers.Update)
			users.DELETE("/:id", controllers.Delete)
		}

		// Customer Profile Routes
		customerProfiles := api.Group("/customer-profiles")
		{
			customerProfiles.POST("/", controllers.CreateCustomerProfile)
			customerProfiles.GET("/", middleware.AuthRequired(), controllers.GetCustomerProfiles)
			customerProfiles.GET("/:id", middleware.AuthRequired(), controllers.GetCustomerProfileByID)
			customerProfiles.PUT("/:id", controllers.UpdateCustomerProfile)
			customerProfiles.DELETE("/:id", middleware.AuthRequired(), controllers.DeleteCustomerProfile)
		}

		// Category Routes
		categories := api.Group("/categories", middleware.AuthRequired())
		{
			categories.POST("/", controllers.CreateCategory)
			categories.GET("/", controllers.GetCategories)
			categories.GET("/:id", controllers.GetCategory)
			categories.PUT("/:id", controllers.UpdateCategory)
			categories.DELETE("/:id", controllers.DeleteCategory)
		}

		// Product Routes
		products := api.Group("/products")
		{
			products.GET("/", controllers.GetProducts)
			products.GET("/:id", controllers.GetProductByID)
			products.POST("/", controllers.CreateProduct)
			products.PUT("/:id", controllers.UpdateProduct)
			products.DELETE("/:id", controllers.DeleteProduct)
		}

		// Order Routes
		orders := api.Group("/orders")
		{
			orders.POST("/", controllers.CreateOrder)
			orders.GET("/", middleware.AuthRequired(), controllers.GetOrders)
			orders.GET("/:id", middleware.AuthRequired(), controllers.GetOrderByID)
			orders.PUT("/:id", controllers.UpdateOrder)
			orders.DELETE("/:id", middleware.AuthRequired(), controllers.DeleteOrder)
		}

		// Role Routes
		roles := api.Group("/roles", middleware.AuthRequired())
		{
			roles.POST("/", controllers.CreateRole)
			roles.GET("/", controllers.GetRoles)
			roles.GET("/:id", controllers.GetRoleByID)
			roles.PUT("/:id", controllers.UpdateRole)
			roles.DELETE("/:id", controllers.DeleteRole)
		}
	}
}

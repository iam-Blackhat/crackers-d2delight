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
			auth.POST("/register/verify-otp", controllers.VerifyRegisterOTP)
		}

		// User routes (protected)
		users := api.Group("/users", middleware.AuthRequired())
		{
			users.GET("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetUsers)
			users.GET("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetUserByID)
			users.PUT("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.Update)
			users.DELETE("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.Delete)
		}

		// Customer Profile Routes
		customerProfiles := api.Group("/customer-profiles")
		{
			customerProfiles.POST("/", middleware.AuthRequired(), middleware.RoleRequired("CUSTOMER"), controllers.CreateCustomerProfile)
			customerProfiles.GET("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetCustomerProfiles)
			customerProfiles.GET("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetCustomerProfileByID)
			customerProfiles.PUT("/:id", middleware.AuthRequired(), middleware.RoleRequired("CUSTOMER"), controllers.UpdateCustomerProfile)
			customerProfiles.DELETE("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.DeleteCustomerProfile)
		}

		// Category Routes
		categories := api.Group("/categories", middleware.AuthRequired())
		{
			categories.POST("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.CreateCategory)
			categories.GET("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetCategories)
			categories.GET("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetCategory)
			categories.PUT("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.UpdateCategory)
			categories.DELETE("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.DeleteCategory)
		}

		// Product Routes
		products := api.Group("/products")
		{
			products.GET("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetProducts)
			products.GET("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetProductByID)
			products.POST("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.CreateProduct)
			products.PUT("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.UpdateProduct)
			products.DELETE("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.DeleteProduct)
		}

		// Order Routes
		orders := api.Group("/orders")
		{
			orders.POST("/", middleware.AuthRequired(), middleware.RoleRequired("CUSTOMER"), controllers.CreateOrder)
			orders.GET("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), middleware.AuthRequired(), controllers.GetOrders)
			orders.GET("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), middleware.AuthRequired(), controllers.GetOrderByID)
			orders.PUT("/:id", middleware.AuthRequired(), middleware.RoleRequired("CUSTOMER"), controllers.UpdateOrder)
			orders.DELETE("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), middleware.AuthRequired(), controllers.DeleteOrder)
		}

		// Role Routes
		roles := api.Group("/roles", middleware.AuthRequired())
		{
			roles.POST("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.CreateRole)
			roles.GET("/", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetRoles)
			roles.GET("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.GetRoleByID)
			roles.PUT("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.UpdateRole)
			roles.DELETE("/:id", middleware.AuthRequired(), middleware.RoleRequired("SUPER ADMIN", "ADMIN"), controllers.DeleteRole)
		}
	}
}

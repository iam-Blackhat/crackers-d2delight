package routes

import (
	"crackers/d2delight.com/controllers"
	"crackers/d2delight.com/middleware"

	"github.com/gin-gonic/gin"
)

func RegisterRoutes(r *gin.Engine) {
	// Health check
	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "ok"})
	})

	// Auth routes
	auth := r.Group("/auth")
	{
		auth.POST("/login", controllers.Login)
		auth.POST("/register", controllers.Create)
	}

	// User routes (protected)
	users := r.Group("/users", middleware.AuthRequired())
	{
		users.GET("/", controllers.GetUsers)
		users.GET("/:id", controllers.GetUserByID)
		users.PUT("/:id", controllers.Update)
		users.DELETE("/:id", controllers.Delete)
	}

	// Customer Routes
	customers := r.Group("/customers")
	{
		customers.POST("/", controllers.CreateCustomer)
		customers.GET("/", middleware.AuthRequired(), controllers.GetCustomers)
		customers.GET("/:id", middleware.AuthRequired(), controllers.GetCustomerByID)
		customers.PUT("/:id", controllers.UpdateCustomer)
		customers.DELETE("/:id", middleware.AuthRequired(), controllers.DeleteCustomer)
	}

	categories := r.Group("/categories", middleware.AuthRequired())
	{
		categories.POST("/", controllers.CreateCategory)
		categories.GET("/", controllers.GetCategories)
		categories.GET("/:id", controllers.GetCategory)
		categories.PUT("/:id", controllers.UpdateCategory)
		categories.DELETE("/:id", controllers.DeleteCategory)
	}

	orders := r.Group("/orders")
	{
		orders.POST("/", controllers.CreateOrder)
		orders.GET("/", middleware.AuthRequired(), controllers.GetOrders)
		orders.GET("/:id", middleware.AuthRequired(), controllers.GetOrderByID)
		orders.PUT("/:id", controllers.UpdateOrder)
		orders.DELETE("/:id", middleware.AuthRequired(), controllers.DeleteOrder)
	}
}

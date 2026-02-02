package product

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	products := r.Group("/products")
	{
		products.POST("", handler.Create)       // Create new product
		products.GET("", handler.GetAll)        // Get products with filters & pagination
		products.GET("/:id", handler.GetByID)   // Get detail product
		products.PUT("/:id", handler.Update)    // Update product info
		products.DELETE("/:id", handler.Delete) // Delete product
	}
}

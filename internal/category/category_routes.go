// internal/category/routes.go
package category

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	categories := r.Group("/categories")
	{
		categories.POST("", handler.Create)
		categories.GET("", handler.GetAll)
		categories.GET("/:id", handler.GetByID)
		categories.PUT("/:id", handler.Update)
		categories.DELETE("/:id", handler.Delete)
	}
}

package customer

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	customers := r.Group("/customers")
	{
		customers.POST("", handler.Create)
		customers.GET("", handler.GetAll)
		customers.GET("/:id", handler.GetByID)
		customers.PUT("/:id", handler.Update)
		customers.DELETE("/:id", handler.Delete)
	}
}

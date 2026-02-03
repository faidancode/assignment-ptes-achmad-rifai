package order

import "github.com/gin-gonic/gin"

func RegisterRoutes(r *gin.RouterGroup, handler *Handler) {
	orders := r.Group("/orders")
	{
		orders.POST("", handler.Create)
		orders.GET("", handler.GetAll)
		orders.GET("/:id", handler.GetByID)
		orders.DELETE("/:id", handler.Delete)
	}
}

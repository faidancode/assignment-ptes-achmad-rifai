package dashboard

import "github.com/gin-gonic/gin"

// RegisterRoutes sekarang hanya menerima router group dan handler yang sudah di-inject
func RegisterRoutes(r *gin.RouterGroup, h *Handler) {
	dashboardGroup := r.Group("/dashboard")
	{
		dashboardGroup.GET("/products", h.GetProductReport)
		dashboardGroup.GET("/top-customers", h.GetTopCustomers)
		dashboardGroup.GET("/overview", h.GetFullDashboard)
	}
}

package dashboard

import (
	"assignment-ptes-achmad-rifai/internal/pkg/response"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

func (h *Handler) GetProductReport(c *gin.Context) {
	res, err := h.service.GetProductDashboard(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "DASHBOARD_ERROR", "Failed to fetch dashboard data", err.Error())
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

func (h *Handler) GetTopCustomers(c *gin.Context) {
	// default value
	limit := int32(10)

	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}

	res, err := h.service.GetTopCustomers(c.Request.Context(), limit)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"DASHBOARD_ERROR",
			"Failed to fetch top customers data",
			err.Error(),
		)
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

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

// GetProductReport godoc
// @Summary      Get product dashboard report
// @Description  Retrieve summary statistics for products, such as total products, active/inactive status, etc.
// @Tags         dashboard
// @Produce      json
// @Success      200      {object}  ProductReportResponse
// @Failure      500      {object}  map[string]string
// @Router       /dashboard/products [get]
func (h *Handler) GetProductReport(c *gin.Context) {
	res, err := h.service.GetProductDashboard(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "DASHBOARD_ERROR", "Failed to fetch dashboard data", err.Error())
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

// GetTopCustomers godoc
// @Summary      Get top performing customers
// @Description  Retrieve a list of customers with the highest transaction volume or spending
// @Tags         dashboard
// @Produce      json
// @Param        limit    query     int  false  "Limit the number of customers returned (default: 5)"
// @Success      200      {array}   TopCustomerResponse
// @Failure      500      {object}  map[string]string
// @Router       /dashboard/top-customers [get]
func (h *Handler) GetTopCustomers(c *gin.Context) {
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

// GetFullDashboard(Overview) godoc
// @Summary      Get complete dashboard overview
// @Description  Retrieve a comprehensive report including financial summaries, top customers, and product stats
// @Tags         dashboard
// @Produce      json
// @Param        limit    query     int  false  "Limit for sub-lists (top customers/recent items)"
// @Success      200      {object}  DashboardReportResponse
// @Failure      500      {object}  map[string]string
// @Router       /dashboard/overview [get]
func (h *Handler) GetFullDashboard(c *gin.Context) {
	limit := int32(10)
	if l := c.Query("limit"); l != "" {
		if parsed, err := strconv.Atoi(l); err == nil && parsed > 0 {
			limit = int32(parsed)
		}
	}

	// Memanggil fungsi concurrency
	res, err := h.service.GetCompleteDashboard(c.Request.Context(), limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "DASHBOARD_ERROR", "Failed to aggregate dashboard", err.Error())
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

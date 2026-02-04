package order

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

// Create godoc
// @Summary      Create a new order
// @Description  Place a new order with multiple items. Calculates total price and quantity automatically.
// @Tags         orders
// @Accept       json
// @Produce      json
// @Param        request body      CreateOrderRequest  true  "Order Request Body"
// @Success      201      {object}  OrderResponse
// @Failure      400      {object}  map[string]string "Invalid input or empty items"
// @Failure      404      {object}  map[string]string "Customer or Product not found"
// @Router       /orders [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateOrderRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "CREATE_ERROR", "Failed to create order", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, res, nil)
}

// GetAll godoc
// @Summary      List all orders
// @Description  Retrieve a paginated list of orders with basic customer info
// @Tags         orders
// @Produce      json
// @Param        page       query    int  false  "Page number"
// @Param        page_size  query    int  false  "Items per page"
// @Success      200      {array}   OrderResponse
// @Failure      500      {object}  map[string]string
// @Router       /orders [get]
func (h *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	params := ListParams{
		Page:     page,
		PageSize: pageSize,
	}
	res, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "FETCH_ERROR", "Failed to fetch orders", err.Error())
		return
	}
	response.Success(c, http.StatusOK, res, nil)
}

// GetByID godoc
// @Summary      Get order details
// @Description  Retrieve full order details including all item descriptions and category names
// @Tags         orders
// @Produce      json
// @Param        id       path      string  true  "Order ID"
// @Success      200      {object}  OrderResponse
// @Failure      404      {object}  map[string]string "Order not found"
// @Router       /orders/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Order not found", err.Error())
		return
	}
	response.Success(c, http.StatusOK, res, nil)
}

// Delete godoc
// @Summary      Delete an order
// @Description  Remove an order record and its associated items
// @Tags         orders
// @Produce      json
// @Param        id       path      string  true  "Order ID"
// @Success      204      {object}  nil
// @Failure      404      {object}  map[string]string
// @Router       /orders/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "DELETE_ERROR", "Failed to delete order", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Order deleted successfully", nil)
}

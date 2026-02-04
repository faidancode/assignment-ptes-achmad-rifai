package customer

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
// @Summary      Create a new customer
// @Description  Register a new customer with a unique email address
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        request body      CreateCustomerRequest  true  "Customer Request"
// @Success      201      {object}  CustomerResponse
// @Failure      400      {object}  map[string]string "Invalid input or email format"
// @Failure      409      {object}  map[string]string "Email already exists"
// @Router       /customers [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "CREATE_ERROR", "Failed to create customer", err.Error())
		return
	}
	response.Success(c, http.StatusCreated, res, nil)
}

// GetAll godoc
// @Summary      List all customers
// @Description  Retrieve a list of all registered customers
// @Tags         customers
// @Produce      json
// @Param        query    query    ListParams  false  "Pagination Query"
// @Success      200      {array}   CustomerResponse
// @Failure      500      {object}  map[string]string
// @Router       /customers [get]
func (h *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	params := ListParams{
		Page:     page,
		PageSize: pageSize,
	}
	res, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "FETCH_ERROR", "Failed to fetch customers", err.Error())
		return
	}
	response.Success(c, http.StatusOK, res, nil)
}

// GetByID godoc
// @Summary      Get customer details
// @Description  Retrieve specific customer information by their unique ID
// @Tags         customers
// @Produce      json
// @Param        id       path      string  true  "Customer ID"
// @Success      200      {object}  CustomerResponse
// @Failure      404      {object}  map[string]string "Customer not found"
// @Router       /customers/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Customer not found", err.Error())
		return
	}
	response.Success(c, http.StatusOK, res, nil)
}

// Update godoc
// @Summary      Update customer information
// @Description  Update name or email for an existing customer
// @Tags         customers
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Customer ID"
// @Param        request  body      UpdateCustomerRequest  true  "Update Request Body"
// @Success      200      {object}  CustomerResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /customers/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateCustomerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, http.StatusBadRequest, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	res, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "UPDATE_ERROR", "Failed to update customer", err.Error())
		return
	}
	response.Success(c, http.StatusOK, res, nil)
}

// Delete godoc
// @Summary      Delete a customer
// @Description  Permanently remove a customer from the database
// @Tags         customers
// @Produce      json
// @Param        id       path      string  true  "Customer ID"
// @Success      204      {object}  nil
// @Failure      404      {object}  map[string]string
// @Router       /customers/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "DELETE_ERROR", "Failed to delete customer", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Customer deleted successfully", nil)
}

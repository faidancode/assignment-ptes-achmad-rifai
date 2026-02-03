package customer

import (
	"assignment-ptes-achmad-rifai/internal/pkg/response"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	service Service
}

func NewHandler(service Service) *Handler {
	return &Handler{service: service}
}

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

func (h *Handler) GetAll(c *gin.Context) {
	res, err := h.service.List(c.Request.Context())
	if err != nil {
		response.Error(c, http.StatusInternalServerError, "FETCH_ERROR", "Failed to fetch customers", err.Error())
		return
	}
	response.Success(c, http.StatusOK, res, nil)
}

func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")
	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(c, http.StatusNotFound, "NOT_FOUND", "Customer not found", err.Error())
		return
	}
	response.Success(c, http.StatusOK, res, nil)
}

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

func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")
	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(c, http.StatusInternalServerError, "DELETE_ERROR", "Failed to delete customer", err.Error())
		return
	}
	response.Success(c, http.StatusOK, "Customer deleted successfully", nil)
}

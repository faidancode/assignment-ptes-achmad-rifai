package product

import (
	"assignment-ptes-achmad-rifai/internal/pkg/response"
	"errors"
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
// @Summary      Create a new product
// @Description  Create a new product with category association and stock details
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        request body      CreateProductRequest  true  "Product Request"
// @Success      201      {object}  ProductResponse
// @Failure      400      {object}  map[string]string
// @Router       /products [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateProductRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(c, 400, "VALIDATION_ERROR", "Invalid request body", err.Error())
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		response.Error(c, 500, "CREATE_ERROR", "Failed to create product", err.Error())
		return
	}

	response.Success(c, http.StatusCreated, res, nil)
}

// GetAll godoc
// @Summary      List products
// @Description  Get a list of products with advanced filters (price, stock, category)
// @Tags         products
// @Produce      json
// @Param        query    query    ListParams  false  "Filter & Pagination Query"
// @Success      200      {array}   ProductResponse
// @Router       /products [get]
func (h *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))

	// Tangkap filter dari query params
	name := c.Query("name")
	categoryID := c.Query("category_id")
	minPriceStr := c.Query("min_price")
	maxPriceStr := c.Query("max_price")
	sortBy := c.DefaultQuery("sort", "name_asc") // Default sort

	params := ListParams{
		Page:     page,
		PageSize: pageSize,
		Sort:     &sortBy,
	}

	// Mapping string ke tipe data yang sesuai (pointer)
	if name != "" {
		params.Name = &name
	}
	if categoryID != "" {
		params.Category = &categoryID
	}

	if minPrice, err := strconv.ParseFloat(minPriceStr, 64); err == nil {
		params.MinPrice = &minPrice
	}
	if maxPrice, err := strconv.ParseFloat(maxPriceStr, 64); err == nil {
		params.MaxPrice = &maxPrice
	}

	data, total, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		response.Error(c, 500, "LIST_ERROR", "Failed to list products", err.Error())
		return
	}

	response.Success(c, 200, data, &response.PaginationMeta{
		Total:      total,
		Page:       page,
		PageSize:   pageSize,
		TotalPages: int((total + int64(pageSize) - 1) / int64(pageSize)),
	})
}

// GetByID godoc
// @Summary      Get product detail
// @Description  Retrieve product information including its category details
// @Tags         products
// @Produce      json
// @Param        id       path      string  true  "Product ID"
// @Success      200      {object}  ProductResponse
// @Failure      404      {object}  map[string]string
// @Router       /products/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			response.Error(c, 404, "NOT_FOUND", err.Error(), nil)
			return
		}
		response.Error(c, 500, "GET_ERROR", "Failed to get product", err.Error())
		return
	}

	response.Success(c, 200, res, nil)
}

// Update godoc
// @Summary      Update product
// @Description  Update product price, stock, or other details by ID
// @Tags         products
// @Accept       json
// @Produce      json
// @Param        id       path      string                true  "Product ID"
// @Param        request  body      UpdateProductRequest  true  "Update Request Body"
// @Success      200      {object}  ProductResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /products/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")
	var req UpdateProductRequest
	res, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			response.Error(c, 404, "NOT_FOUND", err.Error(), nil)
			return
		}
		response.Error(c, 500, "GET_ERROR", "Failed to get product", err.Error())
		return
	}

	response.Success(c, 200, res, nil)
}

// Delete godoc
// @Summary      Delete product
// @Description  Remove a product record from the database
// @Tags         products
// @Produce      json
// @Param        id       path      string  true  "Product ID"
// @Success      204      {object}  nil
// @Failure      404      {object}  map[string]string
// @Router       /products/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"DELETE_ERROR",
			"Failed to delete product",
			err.Error(),
		)
		return
	}

	response.Success(c, http.StatusOK, nil, nil)
}

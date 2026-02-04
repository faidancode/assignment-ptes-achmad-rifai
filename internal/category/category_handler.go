package category

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
	return &Handler{
		service: service,
	}
}

// Create godoc
// @Summary      Create a new category
// @Description  Create a new category with name and description
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        request body      CreateCategoryRequest  true  "Category Request"
// @Success      201      {object}  CategoryResponse
// @Failure      400      {object}  map[string]string
// @Router       /categories [post]
func (h *Handler) Create(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Invalid request body",
			err.Error(),
		)
		return
	}

	res, err := h.service.Create(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, ErrInvalidCategoryName) {
			response.Error(
				c,
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				err.Error(),
				nil,
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"CREATE_ERROR",
			"Failed to create category",
			err.Error(),
		)
		return
	}

	response.Success(c, http.StatusCreated, res, nil)
}

// GetAll godoc
// @Summary      Get all categories
// @Description  Retrieve a list of all categories
// @Tags         categories
// @Produce      json
// @Param        query    query    ListParams  false  "Pagination Query"
// @Success      200      {array}   CategoryResponse
// @Router       /categories [get]
func (h *Handler) GetAll(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "10"))
	params := ListParams{
		Page:     page,
		PageSize: pageSize,
	}
	res, err := h.service.List(c.Request.Context(), params)
	if err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"FETCH_ERROR",
			"Failed to fetch categories",
			err.Error(),
		)
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

// GetByID godoc
// @Summary      Get category by ID
// @Description  Retrieve a single category by its unique ID
// @Tags         categories
// @Produce      json
// @Param        id       path      string  true  "Category ID"
// @Success      200      {object}  CategoryResponse
// @Failure      404      {object}  map[string]string
// @Router       /categories/{id} [get]
func (h *Handler) GetByID(c *gin.Context) {
	id := c.Param("id")

	res, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		response.Error(
			c,
			http.StatusNotFound,
			"NOT_FOUND",
			"Category not found",
			nil,
		)
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

// Update godoc
// @Summary      Update category
// @Description  Update category name or description by ID
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id       path      string                 true  "Category ID"
// @Param        request  body      UpdateCategoryRequest  true  "Update Request"
// @Success      200      {object}  CategoryResponse
// @Failure      400      {object}  map[string]string
// @Failure      404      {object}  map[string]string
// @Router       /categories/{id} [put]
func (h *Handler) Update(c *gin.Context) {
	id := c.Param("id")

	var req UpdateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		response.Error(
			c,
			http.StatusBadRequest,
			"VALIDATION_ERROR",
			"Invalid request body",
			err.Error(),
		)
		return
	}

	res, err := h.service.Update(c.Request.Context(), id, req)
	if err != nil {
		if errors.Is(err, ErrInvalidCategoryName) {
			response.Error(
				c,
				http.StatusBadRequest,
				"VALIDATION_ERROR",
				err.Error(),
				nil,
			)
			return
		}

		response.Error(
			c,
			http.StatusInternalServerError,
			"UPDATE_ERROR",
			"Failed to update category",
			err.Error(),
		)
		return
	}

	response.Success(c, http.StatusOK, res, nil)
}

// Delete godoc
// @Summary      Delete category
// @Description  Remove a category from the system by ID
// @Tags         categories
// @Produce      json
// @Param        id       path      string  true  "Category ID"
// @Success      204      {object}  nil
// @Failure      404      {object}  map[string]string
// @Router       /categories/{id} [delete]
func (h *Handler) Delete(c *gin.Context) {
	id := c.Param("id")

	if err := h.service.Delete(c.Request.Context(), id); err != nil {
		response.Error(
			c,
			http.StatusInternalServerError,
			"DELETE_ERROR",
			"Failed to delete category",
			err.Error(),
		)
		return
	}

	response.Success(c, http.StatusOK, nil, nil)
}

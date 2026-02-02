package category_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"assignment-ptes-achmad-rifai/internal/category"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ==================== FAKE SERVICE ====================

type fakeCategoryService struct {
	CreateFn  func(ctx context.Context, req category.CreateCategoryRequest) (category.CategoryResponse, error)
	ListFn    func(ctx context.Context) ([]category.CategoryResponse, error)
	GetByIDFn func(ctx context.Context, id string) (category.CategoryResponse, error)
	UpdateFn  func(ctx context.Context, id string, req category.UpdateCategoryRequest) (category.CategoryResponse, error)
	DeleteFn  func(ctx context.Context, id string) error
}

func (f *fakeCategoryService) Create(ctx context.Context, req category.CreateCategoryRequest) (category.CategoryResponse, error) {
	return f.CreateFn(ctx, req)
}

func (f *fakeCategoryService) List(ctx context.Context) ([]category.CategoryResponse, error) {
	return f.ListFn(ctx)
}

func (f *fakeCategoryService) GetByID(ctx context.Context, id string) (category.CategoryResponse, error) {
	return f.GetByIDFn(ctx, id)
}

func (f *fakeCategoryService) Update(ctx context.Context, id string, req category.UpdateCategoryRequest) (category.CategoryResponse, error) {
	return f.UpdateFn(ctx, id, req)
}

func (f *fakeCategoryService) Delete(ctx context.Context, id string) error {
	return f.DeleteFn(ctx, id)
}

// ==================== HELPERS ====================

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// ==================== TESTS ====================

func TestHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		desc := "Tech products"
		svc := &fakeCategoryService{
			CreateFn: func(ctx context.Context, req category.CreateCategoryRequest) (category.CategoryResponse, error) {
				assert.Equal(t, "Electronics", req.Name)
				assert.Equal(t, &desc, req.Description)
				return category.CategoryResponse{
					ID:          "uuid-1",
					Name:        req.Name,
					Description: *req.Description,
				}, nil
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.POST("/categories", handler.Create)

		reqBody := category.CreateCategoryRequest{
			Name:        "Electronics",
			Description: &desc,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("validation error - missing name", func(t *testing.T) {
		svc := &fakeCategoryService{}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.POST("/categories", handler.Create)

		reqBody := map[string]interface{}{
			"description": "No name",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid category name", func(t *testing.T) {
		desc := "Test"
		svc := &fakeCategoryService{
			CreateFn: func(ctx context.Context, req category.CreateCategoryRequest) (category.CategoryResponse, error) {
				return category.CategoryResponse{}, category.ErrInvalidCategoryName
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.POST("/categories", handler.Create)

		reqBody := category.CreateCategoryRequest{
			Name:        "Invalid@Name",
			Description: &desc,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		desc := "Test"
		svc := &fakeCategoryService{
			CreateFn: func(ctx context.Context, req category.CreateCategoryRequest) (category.CategoryResponse, error) {
				return category.CategoryResponse{}, errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.POST("/categories", handler.Create)

		reqBody := category.CreateCategoryRequest{
			Name:        "Electronics",
			Description: &desc,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/categories", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCategoryService{
			ListFn: func(ctx context.Context) ([]category.CategoryResponse, error) {
				return []category.CategoryResponse{
					{ID: "1", Name: "Electronics", Description: "Tech"},
					{ID: "2", Name: "Books", Description: "Reading"},
				}, nil
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.GET("/categories", handler.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeCategoryService{
			ListFn: func(ctx context.Context) ([]category.CategoryResponse, error) {
				return nil, errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.GET("/categories", handler.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/categories", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCategoryService{
			GetByIDFn: func(ctx context.Context, id string) (category.CategoryResponse, error) {
				assert.Equal(t, "uuid-1", id)
				return category.CategoryResponse{
					ID:          "uuid-1",
					Name:        "Electronics",
					Description: "Tech products",
				}, nil
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.GET("/categories/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/categories/uuid-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeCategoryService{
			GetByIDFn: func(ctx context.Context, id string) (category.CategoryResponse, error) {
				return category.CategoryResponse{}, category.ErrCategoryNotFound
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.GET("/categories/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/categories/uuid-999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		desc := "Updated description"
		svc := &fakeCategoryService{
			UpdateFn: func(ctx context.Context, id string, req category.UpdateCategoryRequest) (category.CategoryResponse, error) {
				assert.Equal(t, "uuid-1", id)
				assert.Equal(t, "Updated", req.Name)
				return category.CategoryResponse{
					ID:          id,
					Name:        req.Name,
					Description: *req.Description,
				}, nil
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.PUT("/categories/:id", handler.Update)

		reqBody := category.UpdateCategoryRequest{
			Name:        "Updated",
			Description: &desc,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/categories/uuid-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		svc := &fakeCategoryService{}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.PUT("/categories/:id", handler.Update)

		reqBody := map[string]interface{}{
			"description": "No name",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/categories/uuid-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("invalid category name", func(t *testing.T) {
		desc := "Test"
		svc := &fakeCategoryService{
			UpdateFn: func(ctx context.Context, id string, req category.UpdateCategoryRequest) (category.CategoryResponse, error) {
				return category.CategoryResponse{}, category.ErrInvalidCategoryName
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.PUT("/categories/:id", handler.Update)

		reqBody := category.UpdateCategoryRequest{
			Name:        "Invalid@",
			Description: &desc,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/categories/uuid-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		desc := "Test"
		svc := &fakeCategoryService{
			UpdateFn: func(ctx context.Context, id string, req category.UpdateCategoryRequest) (category.CategoryResponse, error) {
				return category.CategoryResponse{}, errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.PUT("/categories/:id", handler.Update)

		reqBody := category.UpdateCategoryRequest{
			Name:        "Updated",
			Description: &desc,
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/categories/uuid-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCategoryService{
			DeleteFn: func(ctx context.Context, id string) error {
				assert.Equal(t, "uuid-1", id)
				return nil
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.DELETE("/categories/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/categories/uuid-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeCategoryService{
			DeleteFn: func(ctx context.Context, id string) error {
				return errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := category.NewHandler(svc)
		r.DELETE("/categories/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/categories/uuid-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

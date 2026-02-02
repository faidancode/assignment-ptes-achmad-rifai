// internal/product/product_handler_test.go
package product_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"assignment-ptes-achmad-rifai/internal/product"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ==================== FAKE SERVICE ====================

type fakeProductService struct {
	CreateFn  func(ctx context.Context, req product.CreateProductRequest) (product.ProductResponse, error)
	ListFn    func(ctx context.Context, params product.ListParams) ([]product.ProductResponse, int64, error)
	GetByIDFn func(ctx context.Context, id string) (product.ProductResponse, error)
	UpdateFn  func(ctx context.Context, id string, req product.UpdateProductRequest) (product.ProductResponse, error)
	DeleteFn  func(ctx context.Context, id string) error
}

func (f *fakeProductService) Create(ctx context.Context, req product.CreateProductRequest) (product.ProductResponse, error) {
	return f.CreateFn(ctx, req)
}
func (f *fakeProductService) List(ctx context.Context, p product.ListParams) ([]product.ProductResponse, int64, error) {
	return f.ListFn(ctx, p)
}
func (f *fakeProductService) GetByID(ctx context.Context, id string) (product.ProductResponse, error) {
	return f.GetByIDFn(ctx, id)
}
func (f *fakeProductService) Update(ctx context.Context, id string, req product.UpdateProductRequest) (product.ProductResponse, error) {
	return f.UpdateFn(ctx, id, req)
}
func (f *fakeProductService) Delete(ctx context.Context, id string) error {
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
		svc := &fakeProductService{
			CreateFn: func(ctx context.Context, req product.CreateProductRequest) (product.ProductResponse, error) {
				return product.ProductResponse{ID: "p-1", Name: req.Name, Price: req.Price}, nil
			},
		}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.POST("/products", handler.Create)

		reqBody, _ := json.Marshal(product.CreateProductRequest{Name: "Laptop", Price: 15000, CategoryID: "cat-123"})
		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)
		t.Logf("Response Status: %d", w.Code)
		t.Logf("Response Body: %s", w.Body.String())

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("error - bad request", func(t *testing.T) {
		svc := &fakeProductService{}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.POST("/products", handler.Create)

		req := httptest.NewRequest(http.MethodPost, "/products", bytes.NewReader([]byte("invalid-json")))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeProductService{
			ListFn: func(ctx context.Context, params product.ListParams) ([]product.ProductResponse, int64, error) {
				assert.Equal(t, 2, params.Page)
				return []product.ProductResponse{{ID: "1", Name: "P1"}}, 10, nil
			},
		}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.GET("/products", handler.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/products?page=2&pageSize=5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error - service failure", func(t *testing.T) {
		svc := &fakeProductService{
			ListFn: func(ctx context.Context, p product.ListParams) ([]product.ProductResponse, int64, error) {
				return nil, 0, errors.New("db error")
			},
		}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.GET("/products", handler.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/products", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeProductService{
			GetByIDFn: func(ctx context.Context, id string) (product.ProductResponse, error) {
				return product.ProductResponse{ID: id, Name: "Coffee"}, nil
			},
		}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.GET("/products/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/products/uuid-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error - not found", func(t *testing.T) {
		svc := &fakeProductService{
			GetByIDFn: func(ctx context.Context, id string) (product.ProductResponse, error) {
				return product.ProductResponse{}, product.ErrProductNotFound
			},
		}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.GET("/products/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/products/none", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeProductService{
			UpdateFn: func(ctx context.Context, id string, req product.UpdateProductRequest) (product.ProductResponse, error) {
				return product.ProductResponse{ID: id, Name: "Updated"}, nil
			},
		}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.PUT("/products/:id", handler.Update)

		reqBody, _ := json.Marshal(product.UpdateProductRequest{Name: "Updated"})
		req := httptest.NewRequest(http.MethodPut, "/products/1", bytes.NewReader(reqBody))
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeProductService{
			DeleteFn: func(ctx context.Context, id string) error {
				return nil
			},
		}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.DELETE("/products/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error - service error", func(t *testing.T) {
		svc := &fakeProductService{
			DeleteFn: func(ctx context.Context, id string) error {
				return errors.New("cannot delete")
			},
		}
		r := setupTestRouter()
		handler := product.NewHandler(svc)
		r.DELETE("/products/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/products/1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

package customer_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"assignment-ptes-achmad-rifai/internal/customer"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ========== FAKE SERVICE ==========

type fakeCustomerService struct {
	CreateFn  func(ctx context.Context, req customer.CreateCustomerRequest) (customer.CustomerResponse, error)
	ListFn    func(ctx context.Context) ([]customer.CustomerResponse, error)
	GetByIDFn func(ctx context.Context, id string) (customer.CustomerResponse, error)
	UpdateFn  func(ctx context.Context, id string, req customer.UpdateCustomerRequest) (customer.CustomerResponse, error)
	DeleteFn  func(ctx context.Context, id string) error
}

func (f *fakeCustomerService) Create(ctx context.Context, req customer.CreateCustomerRequest) (customer.CustomerResponse, error) {
	return f.CreateFn(ctx, req)
}

func (f *fakeCustomerService) List(ctx context.Context) ([]customer.CustomerResponse, error) {
	return f.ListFn(ctx)
}

func (f *fakeCustomerService) GetByID(ctx context.Context, id string) (customer.CustomerResponse, error) {
	return f.GetByIDFn(ctx, id)
}

func (f *fakeCustomerService) Update(ctx context.Context, id string, req customer.UpdateCustomerRequest) (customer.CustomerResponse, error) {
	return f.UpdateFn(ctx, id, req)
}

func (f *fakeCustomerService) Delete(ctx context.Context, id string) error {
	return f.DeleteFn(ctx, id)
}

// ========== HELPERS ==========

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	return gin.New()
}

// ========== TESTS ==========

func TestHandler_Create(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCustomerService{
			CreateFn: func(ctx context.Context, req customer.CreateCustomerRequest) (customer.CustomerResponse, error) {
				assert.Equal(t, "John Doe", req.Name)
				assert.Equal(t, "john@example.com", req.Email)
				return customer.CustomerResponse{
					ID:    "uuid-1",
					Name:  req.Name,
					Email: req.Email,
				}, nil
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.POST("/customers", handler.Create)

		reqBody := customer.CreateCustomerRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("validation error - missing name", func(t *testing.T) {
		svc := &fakeCustomerService{}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.POST("/customers", handler.Create)

		reqBody := map[string]interface{}{
			"email": "john@example.com",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeCustomerService{
			CreateFn: func(ctx context.Context, req customer.CreateCustomerRequest) (customer.CustomerResponse, error) {
				return customer.CustomerResponse{}, errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.POST("/customers", handler.Create)

		reqBody := customer.CreateCustomerRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPost, "/customers", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCustomerService{
			ListFn: func(ctx context.Context) ([]customer.CustomerResponse, error) {
				return []customer.CustomerResponse{
					{ID: "uuid-1", Name: "John Doe", Email: "john@example.com"},
					{ID: "uuid-2", Name: "Jane Doe", Email: "jane@example.com"},
				}, nil
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.GET("/customers", handler.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeCustomerService{
			ListFn: func(ctx context.Context) ([]customer.CustomerResponse, error) {
				return nil, errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.GET("/customers", handler.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/customers", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCustomerService{
			GetByIDFn: func(ctx context.Context, id string) (customer.CustomerResponse, error) {
				assert.Equal(t, "uuid-1", id)
				return customer.CustomerResponse{
					ID:    "uuid-1",
					Name:  "John Doe",
					Email: "john@example.com",
				}, nil
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.GET("/customers/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/customers/uuid-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("not found", func(t *testing.T) {
		svc := &fakeCustomerService{
			GetByIDFn: func(ctx context.Context, id string) (customer.CustomerResponse, error) {
				return customer.CustomerResponse{}, customer.ErrCustomerNotFound
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.GET("/customers/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/customers/uuid-999", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

func TestHandler_Update(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCustomerService{
			UpdateFn: func(ctx context.Context, id string, req customer.UpdateCustomerRequest) (customer.CustomerResponse, error) {
				assert.Equal(t, "uuid-1", id)
				assert.Equal(t, "Updated Name", req.Name)
				assert.Equal(t, "updated@example.com", req.Email)
				return customer.CustomerResponse{
					ID:    id,
					Name:  req.Name,
					Email: req.Email,
				}, nil
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.PUT("/customers/:id", handler.Update)

		reqBody := customer.UpdateCustomerRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/customers/uuid-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("validation error", func(t *testing.T) {
		svc := &fakeCustomerService{}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.PUT("/customers/:id", handler.Update)

		reqBody := map[string]interface{}{
			"email": "updated@example.com",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/customers/uuid-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeCustomerService{
			UpdateFn: func(ctx context.Context, id string, req customer.UpdateCustomerRequest) (customer.CustomerResponse, error) {
				return customer.CustomerResponse{}, errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.PUT("/customers/:id", handler.Update)

		reqBody := customer.UpdateCustomerRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}
		body, _ := json.Marshal(reqBody)

		req := httptest.NewRequest(http.MethodPut, "/customers/uuid-1", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_Delete(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeCustomerService{
			DeleteFn: func(ctx context.Context, id string) error {
				assert.Equal(t, "uuid-1", id)
				return nil
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.DELETE("/customers/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/customers/uuid-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeCustomerService{
			DeleteFn: func(ctx context.Context, id string) error {
				return errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := customer.NewHandler(svc)
		r.DELETE("/customers/:id", handler.Delete)

		req := httptest.NewRequest(http.MethodDelete, "/customers/uuid-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

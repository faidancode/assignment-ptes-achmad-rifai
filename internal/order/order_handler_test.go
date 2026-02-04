package order_test

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"assignment-ptes-achmad-rifai/internal/order"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ========== FAKE SERVICE ==========

type fakeOrderService struct {
	CreateFn  func(ctx context.Context, req order.CreateOrderRequest) (order.OrderResponse, error)
	ListFn    func(ctx context.Context, p order.ListParams) ([]order.OrderResponse, error)
	GetByIDFn func(ctx context.Context, id string) (order.OrderResponse, error)
	DeleteFn  func(ctx context.Context, id string) error
}

func (f *fakeOrderService) Create(ctx context.Context, req order.CreateOrderRequest) (order.OrderResponse, error) {
	return f.CreateFn(ctx, req)
}
func (f *fakeOrderService) List(ctx context.Context, p order.ListParams) ([]order.OrderResponse, error) {
	return f.ListFn(ctx, p)
}
func (f *fakeOrderService) GetByID(ctx context.Context, id string) (order.OrderResponse, error) {
	return f.GetByIDFn(ctx, id)
}
func (f *fakeOrderService) Delete(ctx context.Context, id string) error {
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
		svc := &fakeOrderService{
			CreateFn: func(ctx context.Context, req order.CreateOrderRequest) (order.OrderResponse, error) {
				return order.OrderResponse{ID: "order-1", TotalPrice: 100}, nil
			},
		}

		r := setupTestRouter()
		handler := order.NewHandler(svc)
		r.POST("/orders", handler.Create)

		reqBody := order.CreateOrderRequest{
			CustomerID: "cust-1",
			Items:      []order.OrderItemRequest{{ProductID: "p1", Quantity: 1, UnitPrice: 100}},
		}
		body, _ := json.Marshal(reqBody)
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusCreated, w.Code)
	})

	t.Run("validation error - empty items", func(t *testing.T) {
		// Berikan dummy function agar tidak nil pointer dereference jika terpanggil
		svc := &fakeOrderService{
			CreateFn: func(ctx context.Context, req order.CreateOrderRequest) (order.OrderResponse, error) {
				return order.OrderResponse{}, nil
			},
		}

		r := setupTestRouter()
		handler := order.NewHandler(svc)
		r.POST("/orders", handler.Create)

		// Skenario: Items kosong padahal di DTO ada tag `binding:"required"`
		reqBody := order.CreateOrderRequest{
			CustomerID: "cust-1",
			Items:      []order.OrderItemRequest{}, // Ini akan mentrigger error binding
		}

		body, _ := json.Marshal(reqBody)
		// PASTIKAN body tidak nil
		req := httptest.NewRequest(http.MethodPost, "/orders", bytes.NewReader(body))
		req.Header.Set("Content-Type", "application/json")

		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusBadRequest, w.Code)
	})
}

func TestHandler_GetAll(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeOrderService{
			ListFn: func(ctx context.Context, p order.ListParams) ([]order.OrderResponse, error) {
				return []order.OrderResponse{
					{
						ID:            "order-1",
						CustomerID:    "cust-1",
						TotalQuantity: 2,
						TotalPrice:    100000,
					},
					{
						ID:            "order-2",
						CustomerID:    "cust-2",
						TotalQuantity: 1,
						TotalPrice:    50000,
					},
				}, nil
			},
		}

		r := setupTestRouter()
		handler := order.NewHandler(svc)
		r.GET("/orders", handler.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		// Opsional: Cek isi body jika perlu
		var response map[string]interface{}
		json.Unmarshal(w.Body.Bytes(), &response)
		data := response["data"].([]interface{})
		assert.Len(t, data, 2)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeOrderService{
			ListFn: func(ctx context.Context, p order.ListParams) ([]order.OrderResponse, error) {
				return nil, errors.New("database connection lost")
			},
		}

		r := setupTestRouter()
		handler := order.NewHandler(svc)
		r.GET("/orders", handler.GetAll)

		req := httptest.NewRequest(http.MethodGet, "/orders", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_GetByID(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeOrderService{
			GetByIDFn: func(ctx context.Context, id string) (order.OrderResponse, error) {
				return order.OrderResponse{ID: id}, nil
			},
		}

		r := setupTestRouter()
		handler := order.NewHandler(svc)
		r.GET("/orders/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/orders/uuid-1", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("service error", func(t *testing.T) {
		svc := &fakeOrderService{
			GetByIDFn: func(ctx context.Context, id string) (order.OrderResponse, error) {
				return order.OrderResponse{}, errors.New("not found")
			},
		}

		r := setupTestRouter()
		handler := order.NewHandler(svc)
		r.GET("/orders/:id", handler.GetByID)

		req := httptest.NewRequest(http.MethodGet, "/orders/uuid-99", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusNotFound, w.Code)
	})
}

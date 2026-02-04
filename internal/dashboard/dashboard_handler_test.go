package dashboard_test

import (
	"assignment-ptes-achmad-rifai/internal/dashboard"
	"assignment-ptes-achmad-rifai/internal/pkg/response"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// ========== FAKE SERVICE ==========

type fakeDashboardService struct {
	GetProductDashboardFn  func(ctx context.Context) (dashboard.ProductReportResponse, error)
	GetTopCustomersFn      func(ctx context.Context, limit int32) ([]dashboard.TopCustomerResponse, error)
	GetCompleteDashboardFn func(ctx context.Context, limit int32) (dashboard.DashboardReportResponse, error)
}

func (f *fakeDashboardService) GetProductDashboard(ctx context.Context) (dashboard.ProductReportResponse, error) {
	return f.GetProductDashboardFn(ctx)
}

func (f *fakeDashboardService) GetTopCustomers(ctx context.Context, limit int32) ([]dashboard.TopCustomerResponse, error) {
	return f.GetTopCustomersFn(ctx, limit)
}

func (f *fakeDashboardService) GetCompleteDashboard(ctx context.Context, limit int32) (dashboard.DashboardReportResponse, error) {
	return f.GetCompleteDashboardFn(ctx, limit)
}

// ========== HELPERS ==========

func setupTestRouter() *gin.Engine {
	gin.SetMode(gin.TestMode)
	r := gin.New()
	return r
}

// ========== TESTS ==========

func TestHandler_GetProductReport(t *testing.T) {
	t.Run("success", func(t *testing.T) {
		svc := &fakeDashboardService{
			GetProductDashboardFn: func(ctx context.Context) (dashboard.ProductReportResponse, error) {
				return dashboard.ProductReportResponse{TotalProducts: 100}, nil
			},
		}

		r := setupTestRouter()
		handler := dashboard.NewHandler(svc)
		r.GET("/dashboard/products", handler.GetProductReport)

		req := httptest.NewRequest(http.MethodGet, "/dashboard/products", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("error service", func(t *testing.T) {
		svc := &fakeDashboardService{
			GetProductDashboardFn: func(ctx context.Context) (dashboard.ProductReportResponse, error) {
				return dashboard.ProductReportResponse{}, errors.New("db error")
			},
		}

		r := setupTestRouter()
		handler := dashboard.NewHandler(svc)
		r.GET("/dashboard/products", handler.GetProductReport)

		req := httptest.NewRequest(http.MethodGet, "/dashboard/products", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

func TestHandler_GetTopCustomers(t *testing.T) {
	t.Run("success with default limit", func(t *testing.T) {
		svc := &fakeDashboardService{
			GetTopCustomersFn: func(ctx context.Context, limit int32) ([]dashboard.TopCustomerResponse, error) {
				assert.Equal(t, int32(10), limit) // Default limit check
				return []dashboard.TopCustomerResponse{{Name: "Customer A"}}, nil
			},
		}

		r := setupTestRouter()
		handler := dashboard.NewHandler(svc)
		r.GET("/dashboard/top-customers", handler.GetTopCustomers)

		req := httptest.NewRequest(http.MethodGet, "/dashboard/top-customers", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})

	t.Run("success with custom limit query", func(t *testing.T) {
		svc := &fakeDashboardService{
			GetTopCustomersFn: func(ctx context.Context, limit int32) ([]dashboard.TopCustomerResponse, error) {
				assert.Equal(t, int32(5), limit) // Custom limit check
				return []dashboard.TopCustomerResponse{}, nil
			},
		}

		r := setupTestRouter()
		handler := dashboard.NewHandler(svc)
		r.GET("/dashboard/top-customers", handler.GetTopCustomers)

		req := httptest.NewRequest(http.MethodGet, "/dashboard/top-customers?limit=5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)
	})
}

func TestHandler_GetFullDashboard(t *testing.T) {
	t.Run("success - parallel aggregation", func(t *testing.T) {
		svc := &fakeDashboardService{
			GetCompleteDashboardFn: func(ctx context.Context, limit int32) (dashboard.DashboardReportResponse, error) {
				return dashboard.DashboardReportResponse{
					ProductReport: dashboard.ProductReportResponse{TotalProducts: 50},
					TopCustomers:  []dashboard.TopCustomerResponse{{Name: "Best Buyer"}},
				}, nil
			},
		}

		r := setupTestRouter()
		handler := dashboard.NewHandler(svc)
		r.GET("/dashboard/full", handler.GetFullDashboard)

		req := httptest.NewRequest(http.MethodGet, "/dashboard/full?limit=5", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusOK, w.Code)

		var res response.ApiEnvelope
		err := json.Unmarshal(w.Body.Bytes(), &res)
		assert.NoError(t, err)
		assert.True(t, res.Ok)
		assert.NotNil(t, res.Data)

		// Jika ingin cek isi datanya lebih detail:
		dataMap := res.Data.(map[string]interface{})
		productReport := dataMap["product_report"].(map[string]interface{})
		assert.Equal(t, float64(50), productReport["total_products"])
	})

	t.Run("error - failed to aggregate", func(t *testing.T) {
		svc := &fakeDashboardService{
			GetCompleteDashboardFn: func(ctx context.Context, limit int32) (dashboard.DashboardReportResponse, error) {
				return dashboard.DashboardReportResponse{}, errors.New("concurrency error")
			},
		}

		r := setupTestRouter()
		handler := dashboard.NewHandler(svc)
		r.GET("/dashboard/full", handler.GetFullDashboard)

		req := httptest.NewRequest(http.MethodGet, "/dashboard/full", nil)
		w := httptest.NewRecorder()
		r.ServeHTTP(w, req)

		assert.Equal(t, http.StatusInternalServerError, w.Code)
	})
}

package dashboard_test

import (
	"context"
	"encoding/json"
	"errors"
	"sync"
	"testing"
	"time"

	"assignment-ptes-achmad-rifai/internal/dashboard"
	mockDashboard "assignment-ptes-achmad-rifai/internal/dashboard/mocks"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"assignment-ptes-achmad-rifai/internal/shared/database/helper"

	"github.com/go-redis/redismock/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupServiceTest(t *testing.T) (dashboard.Service, *mockDashboard.MockRepository, redismock.ClientMock) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	// Mock Redis
	dbRedis, redisMock := redismock.NewClientMock()

	// Mock Repo
	repo := mockDashboard.NewMockRepository(ctrl)

	// Create Service
	svc := dashboard.NewService(repo, dbRedis)

	return svc, repo, redisMock
}

func TestService_GetProductDashboard(t *testing.T) {
	ctx := context.Background()
	cacheKey := dashboard.ProductReportKey

	t.Run("Hit Cache - Harus ambil data dari Redis", func(t *testing.T) {
		svc, repo, redisMock := setupServiceTest(t)
		expectedResp := dashboard.ProductReportResponse{
			TotalProducts: 10,
		}
		jsonResp, _ := json.Marshal(expectedResp)

		redisMock.ExpectGet(cacheKey).SetVal(string(jsonResp))

		result, err := svc.GetProductDashboard(ctx)

		assert.NoError(t, err)
		assert.NotNil(t, result.CachedAt)
		assert.Equal(t, expectedResp.TotalProducts, result.TotalProducts)

		repo.EXPECT().GetProductReport(ctx).Times(0)
	})

	t.Run("Miss Cache - Harus ambil dari DB dan simpan ke Redis", func(t *testing.T) {
		svc, repo, redisMock := setupServiceTest(t)
		redisMock.ExpectGet(cacheKey).RedisNil()

		mockReport := dbgen.GetProductDashboardReportRow{TotalProducts: 50}
		mockRecent := []dbgen.GetRecentProductsRow{}

		repo.EXPECT().GetProductReport(ctx).Return(mockReport, nil)
		repo.EXPECT().GetRecentProducts(ctx, int32(5)).Return(mockRecent, nil)

		redisMock.ExpectSet(cacheKey, gomock.Any(), 5*time.Minute).SetVal("OK")

		result, err := svc.GetProductDashboard(ctx)

		assert.NoError(t, err)
		assert.Equal(t, int64(50), result.TotalProducts)
		assert.Nil(t, result.CachedAt) // Dari DB, CachedAt harusnya nil
	})
}

func TestService_GetTopCustomers(t *testing.T) {
	ctx := context.Background()

	t.Run("Success - Get Top Customers", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		mockRows := []dbgen.GetTopCustomersRow{
			{ID: "uuid-1", Name: "Rifai", Email: "rifai@test.com", TotalSpent: decimal.NewFromInt(500000), TotalOrders: 5},
		}

		repo.EXPECT().GetTopCustomers(ctx, int32(5)).Return(mockRows, nil)

		result, err := svc.GetTopCustomers(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 500000.0, result[0].TotalSpent)
		assert.Equal(t, "Rifai", result[0].Name)
	})

	t.Run("Negative - Database Error", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		limit := int32(5)

		// Skenario: Repository mengembalikan error (misal: koneksi putus)
		repo.EXPECT().
			GetTopCustomers(ctx, limit).
			Return(nil, errors.New("database connection lost"))

		result, err := svc.GetTopCustomers(ctx, limit)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database connection lost", err.Error())
	})

	t.Run("Negative - No Data Found", func(t *testing.T) {
		limit := int32(5)
		svc, repo, _ := setupServiceTest(t)

		repo.EXPECT().
			GetTopCustomers(ctx, limit).
			Return([]dbgen.GetTopCustomersRow{}, nil)

		result, err := svc.GetTopCustomers(ctx, limit)

		assert.NoError(t, err)
		assert.Len(t, result, 0)
	})
}

func TestService_GetCompleteDashboard(t *testing.T) {
	ctx := context.Background()
	cacheKey := dashboard.ProductReportKey

	t.Run("Positive - Hybrid Performance (Redis Hit + DB Top Customers)", func(t *testing.T) {
		limit := int32(5)
		svc, repo, redisMock := setupServiceTest(t)

		productData := dashboard.ProductReportResponse{TotalProducts: 100}
		jsonProd, _ := json.Marshal(productData)
		redisMock.ExpectGet(cacheKey).SetVal(string(jsonProd))

		repo.EXPECT().
			GetTopCustomers(ctx, limit).
			Return([]dbgen.GetTopCustomersRow{
				{Name: "Budi", TotalSpent: helper.Float64ToDecimal(1000000)},
			}, nil)

		result, err := svc.GetCompleteDashboard(ctx, limit)

		assert.NoError(t, err)
		assert.NotNil(t, result.ProductReport.CachedAt)
		assert.Equal(t, "Budi", result.TopCustomers[0].Name)

		assert.Nil(t, redisMock.ExpectationsWereMet())
	})

	t.Run("Negative - Concurrency Error (One Fails)", func(t *testing.T) {
		limit := int32(5)
		svc, repo, redisMock := setupServiceTest(t)

		// Redis Hit (Product OK)
		redisMock.ExpectGet(cacheKey).SetVal("{}")

		// Top Customers fail
		repo.EXPECT().
			GetTopCustomers(ctx, limit).
			Return(nil, errors.New("database down"))

		_, err := svc.GetCompleteDashboard(ctx, limit)

		// errgroup catch goroutine error
		assert.Error(t, err)
		assert.Equal(t, "database down", err.Error())
	})
}

func TestService_Singleflight_Proof(t *testing.T) {
	ctx := context.Background()
	cacheKey := dashboard.ProductReportKey

	t.Run("should only call repository once when multiple concurrent requests happen", func(t *testing.T) {
		svc, repo, redisMock := setupServiceTest(t)

		// (Cache Miss)
		redisMock.ExpectGet(cacheKey).RedisNil()

		// Mock Repo with delay
		repo.EXPECT().GetProductReport(gomock.Any()).DoAndReturn(func(ctx context.Context) (dbgen.GetProductDashboardReportRow, error) {
			time.Sleep(100 * time.Millisecond) // Hold request
			return dbgen.GetProductDashboardReportRow{TotalProducts: 100}, nil
		}).Times(1)

		repo.EXPECT().GetRecentProducts(gomock.Any(), gomock.Any()).Return([]dbgen.GetRecentProductsRow{}, nil).Times(1)

		// Mock Redis Set once
		redisMock.ExpectSet(
			cacheKey,
			"",
			time.Minute*5,
		).SetVal("OK")

		// Concurrent Exec
		const numRequests = 10
		var wg sync.WaitGroup
		wg.Add(numRequests)

		// Validation
		results := make(chan error, numRequests)

		for i := 0; i < numRequests; i++ {
			go func() {
				defer wg.Done()
				_, err := svc.GetProductDashboard(ctx)
				results <- err
			}()
		}

		wg.Wait()
		close(results)

		// Verif
		for err := range results {
			assert.NoError(t, err)
		}
	})
}

package dashboard

import (
	"context"
	"encoding/json"
	"errors"
	"testing"
	"time"

	mock "assignment-ptes-achmad-rifai/internal/dashboard/mocks"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"

	"github.com/go-redis/redismock/v9"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestService_GetProductDashboard(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	db, redisMock := redismock.NewClientMock()

	service := NewService(mockRepo, db)
	ctx := context.Background()
	cacheKey := "dashboard:product:report"

	t.Run("Hit Cache - Harus ambil data dari Redis", func(t *testing.T) {
		expectedResp := ProductReportResponse{
			TotalProducts: 10,
		}
		jsonResp, _ := json.Marshal(expectedResp)

		// Ekspektasi: Redis punya data
		redisMock.ExpectGet(cacheKey).SetVal(string(jsonResp))

		// Eksekusi
		result, err := service.GetProductDashboard(ctx)

		// Verifikasi
		assert.NoError(t, err)
		assert.NotNil(t, result.CachedAt) // Pastikan field CachedAt terisi
		assert.Equal(t, expectedResp.TotalProducts, result.TotalProducts)

		// Pastikan Repo TIDAK dipanggil
		mockRepo.EXPECT().GetProductReport(ctx).Times(0)
	})

	t.Run("Miss Cache - Harus ambil dari DB dan simpan ke Redis", func(t *testing.T) {
		// Ekspektasi: Redis kosong
		redisMock.ExpectGet(cacheKey).RedisNil()

		// Mock data dari DB
		mockReport := dbgen.GetProductDashboardReportRow{TotalProducts: 50}
		mockRecent := []dbgen.GetRecentProductsRow{}

		mockRepo.EXPECT().GetProductReport(ctx).Return(mockReport, nil)
		mockRepo.EXPECT().GetRecentProducts(ctx, int32(5)).Return(mockRecent, nil)

		// Ekspektasi: Simpan ke Redis setelah ambil dari DB
		redisMock.ExpectSet(cacheKey, gomock.Any(), 5*time.Minute).SetVal("OK")

		// Eksekusi
		result, err := service.GetProductDashboard(ctx)

		// Verifikasi
		assert.NoError(t, err)
		assert.Equal(t, int64(50), result.TotalProducts)
		assert.Nil(t, result.CachedAt) // Dari DB, CachedAt harusnya nil
	})
}

func TestService_GetTopCustomers(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mock.NewMockRepository(ctrl)
	// Kita asumsikan dashboard service menggunakan repo yang sama
	service := NewService(mockRepo, nil)
	ctx := context.Background()

	t.Run("Success - Get Top Customers", func(t *testing.T) {
		mockRows := []dbgen.GetTopCustomersRow{
			{ID: "uuid-1", Name: "Rifai", Email: "rifai@test.com", TotalSpent: decimal.NewFromInt(500000), TotalOrders: 5},
		}

		mockRepo.EXPECT().GetTopCustomers(ctx, int32(5)).Return(mockRows, nil)

		result, err := service.GetTopCustomers(ctx, 5)

		assert.NoError(t, err)
		assert.Equal(t, 500000.0, result[0].TotalSpent)
		assert.Equal(t, "Rifai", result[0].Name)
	})

	t.Run("Negative - Database Error", func(t *testing.T) {
		limit := int32(5)

		// Skenario: Repository mengembalikan error (misal: koneksi putus)
		mockRepo.EXPECT().
			GetTopCustomers(ctx, limit).
			Return(nil, errors.New("database connection lost"))

		result, err := service.GetTopCustomers(ctx, limit)

		assert.Error(t, err)
		assert.Nil(t, result)
		assert.Equal(t, "database connection lost", err.Error())
	})

	t.Run("Negative - No Data Found", func(t *testing.T) {
		limit := int32(5)

		// Skenario: Database normal, tapi tidak ada data order sama sekali (array kosong)
		mockRepo.EXPECT().
			GetTopCustomers(ctx, limit).
			Return([]dbgen.GetTopCustomersRow{}, nil)

		result, err := service.GetTopCustomers(ctx, limit)

		assert.NoError(t, err)
		assert.Len(t, result, 0) // Hasil harus array kosong, bukan error
	})
}

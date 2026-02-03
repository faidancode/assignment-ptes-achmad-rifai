package dashboard

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	mock "assignment-ptes-achmad-rifai/internal/dashboard/mocks"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"

	"github.com/go-redis/redismock/v9"
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

package dashboard

import (
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
)

type Service interface {
	GetProductDashboard(ctx context.Context) (ProductReportResponse, error)
	GetTopCustomers(ctx context.Context, limit int32) ([]TopCustomerResponse, error)
}

type service struct {
	repo Repository
	rdb  *redis.Client
}

func NewService(repo Repository, rdb *redis.Client) Service {
	return &service{repo: repo, rdb: rdb}
}

func (s *service) GetProductDashboard(ctx context.Context) (ProductReportResponse, error) {
	cacheKey := "dashboard:product:report"

	// 1. Coba ambil dari Redis
	cachedData, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp ProductReportResponse
		if err := json.Unmarshal([]byte(cachedData), &resp); err == nil {
			now := time.Now()
			resp.CachedAt = &now
			return resp, nil
		}
	}

	// 2. Jika tidak ada di cache, ambil dari DB
	report, err := s.repo.GetProductReport(ctx)
	if err != nil {
		return ProductReportResponse{}, err
	}

	recent, err := s.repo.GetRecentProducts(ctx, 5) // Limit 5 produk terbaru
	if err != nil {
		return ProductReportResponse{}, err
	}

	// 3. Mapping data
	recentResp := make([]RecentProductResponse, 0)
	for _, p := range recent {
		price, _ := p.Price.Float64()
		recentResp = append(recentResp, RecentProductResponse{
			ID: p.ID, Name: p.Name, Price: price, StockQuantity: p.StockQuantity, CreatedAt: p.CreatedAt,
		})
	}

	avgPrice, _ := report.AvgPrice.Float64()
	finalResp := ProductReportResponse{
		TotalProducts:  report.TotalProducts,
		TotalStock:     report.TotalStock,
		AveragePrice:   avgPrice,
		RecentProducts: recentResp,
	}

	// 4. Simpan ke Redis (TTL 5 menit)
	jsonData, _ := json.Marshal(finalResp)
	s.rdb.Set(ctx, cacheKey, jsonData, 5*time.Minute)

	return finalResp, nil
}

// internal/modules/dashboard/service.go

func (s *service) GetTopCustomers(ctx context.Context, limit int32) ([]TopCustomerResponse, error) {
	// 1. Ambil data dari Repository (hasil generate SQLC)
	rows, err := s.repo.GetTopCustomers(ctx, limit)
	if err != nil {
		return nil, err
	}

	// 2. Mapping ke DTO
	var resp []TopCustomerResponse
	for _, r := range rows {
		spent, _ := r.TotalSpent.Float64()
		resp = append(resp, TopCustomerResponse{
			ID:          r.ID,
			Name:        r.Name,
			Email:       r.Email,
			TotalSpent:  spent,
			TotalOrders: r.TotalOrders,
		})
	}

	return resp, nil
}

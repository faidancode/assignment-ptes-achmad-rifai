package dashboard

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
	"encoding/json"
	"time"

	"github.com/redis/go-redis/v9"
	"golang.org/x/sync/errgroup"
	"golang.org/x/sync/singleflight"
)

const (
	ProductReportKey = "dashboard:product:report"
	TopCustomerKey   = "dashboard:customer:top"
)

//go:generate mockgen -source=dashboard_service.go -destination=mocks/dashboard_service_mock.go -package=mock

type Service interface {
	GetProductDashboard(ctx context.Context) (ProductReportResponse, error)
	GetTopCustomers(ctx context.Context, limit int32) ([]TopCustomerResponse, error)
	GetCompleteDashboard(ctx context.Context, limit int32) (DashboardReportResponse, error)
}

type service struct {
	repo Repository
	rdb  *redis.Client
	sf   *singleflight.Group
}

func NewService(repo Repository, rdb *redis.Client) Service {
	return &service{repo: repo, rdb: rdb, sf: &singleflight.Group{}}
}

func (s *service) GetProductDashboard(ctx context.Context) (ProductReportResponse, error) {
	cacheKey := ProductReportKey

	// 1. Hit Redis
	if cached, err := s.rdb.Get(ctx, cacheKey).Result(); err == nil {
		var resp ProductReportResponse
		if json.Unmarshal([]byte(cached), &resp) == nil {
			now := time.Now()
			resp.CachedAt = &now
			return resp, nil
		}
	}

	// 2. Gunakan Singleflight untuk mencegah Cache Stampede
	v, err, _ := s.sf.Do(cacheKey, func() (interface{}, error) {
		// Eksekusi paralel di dalam singleflight
		var g errgroup.Group
		var report dbgen.GetProductDashboardReportRow
		var recent []dbgen.GetRecentProductsRow

		g.Go(func() error {
			var err error
			report, err = s.repo.GetProductReport(ctx)
			return err
		})

		g.Go(func() error {
			var err error
			recent, err = s.repo.GetRecentProducts(ctx, 5)
			return err
		})

		if err := g.Wait(); err != nil {
			return nil, err
		}

		// Mapping
		recentResp := make([]RecentProductResponse, 0, len(recent))
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

		// Simpan ke Cache
		if jsonData, err := json.Marshal(finalResp); err == nil {
			s.rdb.Set(ctx, cacheKey, jsonData, 5*time.Minute)
		}

		return finalResp, nil
	})

	if err != nil {
		return ProductReportResponse{}, err
	}

	return v.(ProductReportResponse), nil
}

func (s *service) GetTopCustomers(ctx context.Context, limit int32) ([]TopCustomerResponse, error) {
	cacheKey := TopCustomerKey

	// 1. Coba ambil dari Redis
	cachedData, err := s.rdb.Get(ctx, cacheKey).Result()
	if err == nil {
		var resp []TopCustomerResponse
		if err := json.Unmarshal([]byte(cachedData), &resp); err == nil {
			return resp, nil
		}
	}

	// 2. Gunakan Singleflight untuk mencegah bentrokan request ke DB
	v, err, _ := s.sf.Do(cacheKey, func() (interface{}, error) {
		rows, err := s.repo.GetTopCustomers(ctx, limit)
		if err != nil {
			return nil, err
		}

		// Alokasi capacity yang pas agar lebih cepat
		resp := make([]TopCustomerResponse, 0, len(rows))
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

		// 3. Simpan ke Redis (TTL 5 menit)
		if jsonData, err := json.Marshal(resp); err == nil {
			s.rdb.Set(ctx, cacheKey, jsonData, 5*time.Minute)
		}

		return resp, nil
	})

	if err != nil {
		return nil, err
	}

	return v.([]TopCustomerResponse), nil
}

func (s *service) GetCompleteDashboard(ctx context.Context, topCustomerLimit int32) (DashboardReportResponse, error) {
	var g errgroup.Group
	var productReport ProductReportResponse // Pakai tipe data aslimu
	var topCustomers []TopCustomerResponse

	// Goroutine 1: Redis + DB Fallback
	g.Go(func() error {
		var err error
		productReport, err = s.GetProductDashboard(ctx)
		return err
	})

	// Goroutine 2: Heavy DB Query
	g.Go(func() error {
		var err error
		topCustomers, err = s.GetTopCustomers(ctx, topCustomerLimit)
		return err
	})

	if err := g.Wait(); err != nil {
		return DashboardReportResponse{}, err
	}

	return DashboardReportResponse{
		ProductReport: productReport,
		TopCustomers:  topCustomers,
	}, nil
}

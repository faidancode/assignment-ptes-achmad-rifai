package dashboard

import (
	"time"
)

type ProductReportResponse struct {
	TotalProducts  int64                   `json:"total_products"`
	TotalStock     int64                   `json:"total_stock"`
	AveragePrice   float64                 `json:"average_price"`
	RecentProducts []RecentProductResponse `json:"recent_products"`
	CachedAt       *time.Time              `json:"cached_at,omitempty"` // Penanda jika data dari cache
}

type RecentProductResponse struct {
	ID            string    `json:"id"`
	Name          string    `json:"name"`
	Price         float64   `json:"price"`
	StockQuantity int32     `json:"stock_quantity"`
	CreatedAt     time.Time `json:"created_at"`
}

type TopCustomerResponse struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Email       string  `json:"email"`
	TotalSpent  float64 `json:"total_spent"`
	TotalOrders int64   `json:"total_orders"`
}

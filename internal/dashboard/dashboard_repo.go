package dashboard

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
)

//go:generate mockgen -source=dashboard_repo.go -destination=mocks/dashboard_repo_mock.go -package=mock
type Repository interface {
	GetProductReport(ctx context.Context) (dbgen.GetProductDashboardReportRow, error)
	GetRecentProducts(ctx context.Context, limit int32) ([]dbgen.GetRecentProductsRow, error)

	GetTopCustomers(ctx context.Context, limit int32) ([]dbgen.GetTopCustomersRow, error)
}

type repository struct {
	q *dbgen.Queries
}

func NewRepository(q *dbgen.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) GetProductReport(ctx context.Context) (dbgen.GetProductDashboardReportRow, error) {
	return r.q.GetProductDashboardReport(ctx)
}

func (r *repository) GetRecentProducts(ctx context.Context, limit int32) ([]dbgen.GetRecentProductsRow, error) {
	return r.q.GetRecentProducts(ctx, limit)
}

func (r *repository) GetTopCustomers(ctx context.Context, limit int32) ([]dbgen.GetTopCustomersRow, error) {
	return r.q.GetTopCustomers(ctx, limit)
}

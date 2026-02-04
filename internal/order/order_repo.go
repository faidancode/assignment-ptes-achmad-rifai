package order

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
	"database/sql"
)

//go:generate mockgen -source=order_repo.go -destination=mocks/order_repo_mock.go -package=mock
type Repository interface {
	// Transaction helpers
	WithTx(tx dbgen.DBTX) Repository

	CreateOrder(ctx context.Context, params dbgen.CreateOrderParams) error
	CreateOrderItem(ctx context.Context, params dbgen.CreateOrderItemParams) error
	GetOrders(ctx context.Context, params dbgen.GetOrdersParams) ([]dbgen.GetOrdersRow, error)
	GetByID(ctx context.Context, id string) (dbgen.GetOrderByIDRow, error)
	GetItemsByOrderID(ctx context.Context, orderID string) ([]dbgen.OrderItem, error)
	Delete(ctx context.Context, id string) error
}

type repository struct {
	q *dbgen.Queries
}

func NewRepository(q *dbgen.Queries) Repository {
	return &repository{q: q}
}

func (r *repository) WithTx(tx dbgen.DBTX) Repository {
	if sqlTx, ok := tx.(*sql.Tx); ok {
		return &repository{
			q: r.q.WithTx(sqlTx),
		}
	}

	return r
}

// ... implementasi method lainnya (CreateOrder, GetOrders, dll) tetap sama

func (r *repository) CreateOrder(ctx context.Context, params dbgen.CreateOrderParams) error {
	return r.q.CreateOrder(ctx, params)
}

func (r *repository) CreateOrderItem(ctx context.Context, params dbgen.CreateOrderItemParams) error {
	return r.q.CreateOrderItem(ctx, params)
}

func (r *repository) GetOrders(ctx context.Context, params dbgen.GetOrdersParams) ([]dbgen.GetOrdersRow, error) {
	return r.q.GetOrders(ctx, params)
}

func (r *repository) GetByID(ctx context.Context, id string) (dbgen.GetOrderByIDRow, error) {
	return r.q.GetOrderByID(ctx, id)
}

func (r *repository) GetItemsByOrderID(ctx context.Context, orderID string) ([]dbgen.OrderItem, error) {
	return r.q.GetOrderItemsByOrderID(ctx, orderID)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.q.DeleteOrder(ctx, id)
}

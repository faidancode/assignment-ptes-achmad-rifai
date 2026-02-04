package customer

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
)

//go:generate mockgen -source=customer_repo.go -destination=mocks/customer_repo_mock.go -package=mock
type Repository interface {
	Create(ctx context.Context, params dbgen.CreateCustomerParams) error
	GetCustomers(ctx context.Context, params dbgen.GetCustomersParams) ([]dbgen.GetCustomersRow, error)
	GetByID(ctx context.Context, id string) (dbgen.GetCustomerByIDRow, error)
	Update(ctx context.Context, params dbgen.UpdateCustomerParams) error
	Delete(ctx context.Context, id string) error
}

type repository struct {
	q *dbgen.Queries
}

// Ubah parameter dari *sql.DB menjadi *dbgen.Queries
func NewRepository(q *dbgen.Queries) Repository {
	return &repository{
		q: q,
	}
}

func (r *repository) Create(ctx context.Context, params dbgen.CreateCustomerParams) error {
	return r.q.CreateCustomer(ctx, params)
}

func (r *repository) GetCustomers(ctx context.Context, params dbgen.GetCustomersParams) ([]dbgen.GetCustomersRow, error) {
	return r.q.GetCustomers(ctx, params)
}

func (r *repository) GetByID(ctx context.Context, id string) (dbgen.GetCustomerByIDRow, error) {
	return r.q.GetCustomerByID(ctx, id)
}

func (r *repository) Update(ctx context.Context, params dbgen.UpdateCustomerParams) error {
	return r.q.UpdateCustomer(ctx, params)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.q.DeleteCustomer(ctx, id)
}

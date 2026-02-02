package product

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=product_repo.go -destination=mocks/product_repo_mock.go -package=mock
type Repository interface {
	Create(ctx context.Context, params dbgen.CreateProductParams) error
	GetByID(ctx context.Context, id string) (dbgen.GetProductByIDRow, error)
	List(ctx context.Context, params dbgen.ListProductsParams) ([]dbgen.ListProductsRow, error)
	Count(ctx context.Context, params dbgen.CountProductsParams) (int64, error)
	Update(ctx context.Context, params dbgen.UpdateProductParams) error
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

func (r *repository) Create(ctx context.Context, params dbgen.CreateProductParams) error {
	params.ID = uuid.New().String()
	return r.q.CreateProduct(ctx, params)
}

func (r *repository) GetByID(ctx context.Context, id string) (dbgen.GetProductByIDRow, error) {
	return r.q.GetProductByID(ctx, id)
}

func (r *repository) List(
	ctx context.Context,
	params dbgen.ListProductsParams,
) ([]dbgen.ListProductsRow, error) {
	return r.q.ListProducts(ctx, params)
}

func (r *repository) Count(
	ctx context.Context,
	params dbgen.CountProductsParams,
) (int64, error) {
	return r.q.CountProducts(ctx, params)
}

func (r *repository) Update(ctx context.Context, params dbgen.UpdateProductParams) error {
	return r.q.UpdateProduct(ctx, params)
}

func (r *repository) Delete(ctx context.Context, id string) error {
	return r.q.DeleteProduct(ctx, id)
}

package category

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"

	"github.com/google/uuid"
)

/*
Repository
*/
//go:generate mockgen -source=category_repo.go -destination=mocks/category_repo_mock.go -package=mock
type Repository interface {
	Create(ctx context.Context, params dbgen.CreateCategoryParams) error
	GetCategories(ctx context.Context, params dbgen.GetCategoriesParams) ([]dbgen.GetCategoriesRow, error)
	GetByID(ctx context.Context, id string) (dbgen.GetCategoryByIDRow, error)
	Update(ctx context.Context, params dbgen.UpdateCategoryParams) error
	Delete(ctx context.Context, id string) error
}

/*
sqlc implementation
*/

type repository struct {
	q *dbgen.Queries
}

// Ubah parameter dari *sql.DB menjadi *dbgen.Queries
func NewRepository(q *dbgen.Queries) Repository {
	return &repository{
		q: q,
	}
}

func (r *repository) Create(
	ctx context.Context,
	params dbgen.CreateCategoryParams,
) error {
	params.ID = uuid.New().String()
	return r.q.CreateCategory(ctx, params)
}

func (r *repository) GetCategories(
	ctx context.Context, params dbgen.GetCategoriesParams,
) ([]dbgen.GetCategoriesRow, error) {
	return r.q.GetCategories(ctx, params)
}

func (r *repository) GetByID(
	ctx context.Context,
	id string,
) (dbgen.GetCategoryByIDRow, error) {
	return r.q.GetCategoryByID(ctx, id)
}

func (r *repository) Update(
	ctx context.Context,
	params dbgen.UpdateCategoryParams,
) error {
	return r.q.UpdateCategory(ctx, params)
}

func (r *repository) Delete(
	ctx context.Context,
	id string,
) error {
	return r.q.DeleteCategory(ctx, id)
}

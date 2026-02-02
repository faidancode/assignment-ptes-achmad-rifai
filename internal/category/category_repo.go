package category

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
	"database/sql"
)

/*
Repository
*/
//go:generate mockgen -source=category_repo.go -destination=mocks/category_repo_mock.go -package=mock
type Repository interface {
	Create(ctx context.Context, params dbgen.CreateCategoryParams) error
	GetCategories(ctx context.Context) ([]dbgen.Category, error)
	GetByID(ctx context.Context, id string) (dbgen.Category, error)
	Update(ctx context.Context, params dbgen.UpdateCategoryParams) error
	Delete(ctx context.Context, id string) error
}

/*
sqlc implementation
*/

type repository struct {
	q *dbgen.Queries
}

func NewRepository(db *sql.DB) Repository {
	return &repository{
		q: dbgen.New(db),
	}
}

func (r *repository) Create(
	ctx context.Context,
	params dbgen.CreateCategoryParams,
) error {
	return r.q.CreateCategory(ctx, params)
}

func (r *repository) GetCategories(
	ctx context.Context,
) ([]dbgen.Category, error) {
	return r.q.GetCategories(ctx)
}

func (r *repository) GetByID(
	ctx context.Context,
	id string,
) (dbgen.Category, error) {
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

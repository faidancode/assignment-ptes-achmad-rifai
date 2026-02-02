package category

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"assignment-ptes-achmad-rifai/internal/shared/database/helper"
	"context"

	"github.com/google/uuid"
)

//go:generate mockgen -source=category_service.go -destination=mocks/category_service_mock.go -package=mock
type Service interface {
	Create(ctx context.Context, req CreateCategoryRequest) (CategoryResponse, error)
	List(ctx context.Context) ([]CategoryResponse, error)
	GetByID(ctx context.Context, id string) (CategoryResponse, error)
	Update(ctx context.Context, id string, req UpdateCategoryRequest) (CategoryResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{
		repo: repo,
	}
}

func (s *service) Create(
	ctx context.Context,
	req CreateCategoryRequest,
) (CategoryResponse, error) {

	id := uuid.NewString()

	params := dbgen.CreateCategoryParams{
		ID:          id,
		Name:        req.Name,
		Description: helper.NewNullString(req.Description),
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return CategoryResponse{}, err
	}

	// 🔥 mapping langsung, tanpa query tambahan
	return CategoryResponse{
		ID:          id,
		Name:        req.Name,
		Description: helper.StringValue(req.Description),
	}, nil
}

func (s *service) List(ctx context.Context) ([]CategoryResponse, error) {
	rows, err := s.repo.GetCategories(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]CategoryResponse, 0, len(rows))
	for _, row := range rows {
		// Panggil helper dengan parameter mentah
		res = append(res, toResponse(row.ID, row.Name))
	}

	return res, nil
}

func (s *service) GetByID(
	ctx context.Context,
	id string,
) (CategoryResponse, error) {

	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return CategoryResponse{}, err
	}

	return mapToResponse(cat), nil
}

func (s *service) Update(
	ctx context.Context,
	id string,
	req UpdateCategoryRequest,
) (CategoryResponse, error) {

	if err := s.repo.Update(ctx, dbgen.UpdateCategoryParams{
		ID:          id,
		Name:        req.Name,
		Description: helper.NewNullString(req.Description),
	}); err != nil {
		return CategoryResponse{}, err
	}

	cat, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return CategoryResponse{}, err
	}

	return mapToResponse(cat), nil
}

func (s *service) Delete(
	ctx context.Context,
	id string,
) error {
	return s.repo.Delete(ctx, id)
}

/*
Helper
*/

func mapToResponse(cat dbgen.GetCategoryByIDRow) CategoryResponse {
	return CategoryResponse{
		ID:   cat.ID,
		Name: cat.Name,
	}
}

func toResponse(id string, name string) CategoryResponse {
	return CategoryResponse{
		ID:   id,
		Name: name,
	}
}

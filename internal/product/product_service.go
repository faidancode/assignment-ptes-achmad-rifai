package product

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"assignment-ptes-achmad-rifai/internal/shared/database/helper"
	"context"
	"database/sql"
)

//go:generate mockgen -source=product_service.go -destination=mocks/product_service_mock.go -package=mock
type Service interface {
	Create(ctx context.Context, req CreateProductRequest) (ProductResponse, error)
	List(ctx context.Context, params ListParams) ([]ProductResponse, int64, error)
	GetByID(ctx context.Context, id string) (ProductResponse, error)
	Update(ctx context.Context, id string, req UpdateProductRequest) (ProductResponse, error)
	Delete(ctx context.Context, id string) error
}
type ListParams struct {
	Page     int
	PageSize int
	Name     *string
	Category *string
	Sort     *string
}
type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}
func (s *service) Create(
	ctx context.Context,
	req CreateProductRequest,
) (ProductResponse, error) {

	params := dbgen.CreateProductParams{
		Name:          req.Name,
		Description:   helper.NewNullString(req.Description),
		Price:         helper.ToDecimal(req.Price),
		CategoryID:    req.CategoryID,
		StockQuantity: int32(req.StockQuantity),
		IsActive:      helper.BoolValue(req.IsActive, true),
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return ProductResponse{}, err
	}

	return ProductResponse{
		Name:          req.Name,
		Description:   helper.StringValue(req.Description),
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		IsActive:      helper.BoolValue(req.IsActive, true),
		Category: CategoryResponse{
			ID: req.CategoryID,
		},
	}, nil
}
func (s *service) List(
	ctx context.Context,
	p ListParams,
) ([]ProductResponse, int64, error) {

	limit := int32(p.PageSize)
	offset := int32((p.Page - 1) * p.PageSize)

	rows, err := s.repo.List(ctx, dbgen.ListProductsParams{
		SearchName: helper.NewNullString(p.Name),
		CategoryID: helper.NewNullString(p.Category),
		Limit:      limit,
		Offset:     offset,
	})
	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, dbgen.CountProductsParams{
		SearchName: helper.NewNullString(p.Name),
		CategoryID: helper.NewNullString(p.Category),
	})
	if err != nil {
		return nil, 0, err
	}

	res := make([]ProductResponse, 0, len(rows))
	for _, r := range rows {
		res = append(res, mapListToResponse(r))
	}

	return res, total, nil
}
func (s *service) GetByID(ctx context.Context, id string) (ProductResponse, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		if err == sql.ErrNoRows {
			return ProductResponse{}, ErrProductNotFound
		}
		return ProductResponse{}, err
	}

	return mapDetailToResponse(row), nil
}
func (s *service) Update(
	ctx context.Context,
	id string,
	req UpdateProductRequest,
) (ProductResponse, error) {

	params := dbgen.UpdateProductParams{
		ID:            id,
		Name:          req.Name,
		Description:   helper.NewNullString(req.Description),
		Price:         helper.ToDecimal(req.Price),
		CategoryID:    req.CategoryID,
		StockQuantity: int32(req.StockQuantity),
		IsActive:      helper.BoolValue(req.IsActive, true),
	}

	if err := s.repo.Update(ctx, params); err != nil {
		return ProductResponse{}, err
	}

	return s.GetByID(ctx, id)
}
func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

// Mapper khusus untuk hasil List
func mapListToResponse(r dbgen.ListProductsRow) ProductResponse {
	return ProductResponse{
		ID:            r.ID,
		Name:          r.Name,
		Description:   r.Description.String,
		Price:         helper.FloatFromDecimal(r.Price),
		StockQuantity: int(r.StockQuantity),
		IsActive:      r.IsActive,
		Category: CategoryResponse{
			ID:   r.CategoryID,
			Name: r.CategoryName,
		},
	}
}

// Mapper khusus untuk hasil GetByID
func mapDetailToResponse(r dbgen.GetProductByIDRow) ProductResponse {
	return ProductResponse{
		ID:            r.ID,
		Name:          r.Name,
		Description:   r.Description.String,
		Price:         helper.FloatFromDecimal(r.Price),
		StockQuantity: int(r.StockQuantity),
		IsActive:      r.IsActive,
		Category: CategoryResponse{
			ID:   r.CategoryID,
			Name: r.CategoryName,
		},
	}
}

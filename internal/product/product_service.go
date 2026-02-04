package product

import (
	"assignment-ptes-achmad-rifai/internal/dashboard"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"assignment-ptes-achmad-rifai/internal/shared/database/helper"
	"context"
	"database/sql"
	"log"

	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
)

//go:generate mockgen -source=product_service.go -destination=mocks/product_service_mock.go -package=mock
type Service interface {
	Create(ctx context.Context, req CreateProductRequest) (ProductResponse, error)
	List(ctx context.Context, params ListParams) ([]ProductResponse, int64, error)
	GetByID(ctx context.Context, id string) (ProductResponse, error)
	Update(ctx context.Context, id string, req UpdateProductRequest) (ProductResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
	rdb  *redis.Client
}

func NewService(repo Repository, rdb *redis.Client) Service {
	return &service{repo: repo, rdb: rdb}
}
func (s *service) Create(
	ctx context.Context,
	req CreateProductRequest,
) (ProductResponse, error) {
	newUUID, err := uuid.NewV7()
	if err != nil {
		return ProductResponse{}, err
	}

	productID := newUUID.String()
	params := dbgen.CreateProductParams{
		ID:            productID,
		Name:          req.Name,
		Description:   helper.StringToNull(req.Description),
		Price:         helper.Float64ToDecimal(req.Price),
		CategoryID:    req.CategoryID,
		StockQuantity: int32(req.StockQuantity),
		IsActive:      helper.BoolPtrValue(req.IsActive, true),
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return ProductResponse{}, err
	}

	dashboardCacheKey := dashboard.ProductReportKey
	if err := s.rdb.Del(ctx, dashboardCacheKey).Err(); err != nil {
		log.Printf("failed to invalidate dashboard product cache: %v", err)
	}

	return ProductResponse{
		ID:            productID,
		Name:          req.Name,
		Description:   helper.StringPtrValue(req.Description),
		Price:         req.Price,
		StockQuantity: req.StockQuantity,
		IsActive:      helper.BoolPtrValue(req.IsActive, true),
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
		SearchName: helper.StringPtrValue(p.Name),          // "" = no filter
		CategoryID: helper.StringPtrValue(p.Category),      // "" = no filter
		MinPrice:   helper.Float64PtrToDecimal(p.MinPrice), // 0 = no filter
		MaxPrice:   helper.Float64PtrToDecimal(p.MaxPrice),
		MinStock:   helper.Int32PtrValue(p.MinStock), // 0 = no filter
		MaxStock:   helper.Int32PtrValue(p.MaxStock),
		OrderBy:    helper.StringPtrValue(p.Sort),
		Limit:      limit,
		Offset:     offset,
	})

	if err != nil {
		return nil, 0, err
	}

	total, err := s.repo.Count(ctx, dbgen.CountProductsParams{
		SearchName: helper.StringPtrValue(p.Name),          // "" = no filter
		CategoryID: helper.StringPtrValue(p.Category),      // "" = no filter
		MinPrice:   helper.Float64PtrToDecimal(p.MinPrice), // 0 = no filter
		MaxPrice:   helper.Float64PtrToDecimal(p.MaxPrice),
		MinStock:   helper.Int32PtrValue(p.MinStock), // 0 = no filter
		MaxStock:   helper.Int32PtrValue(p.MaxStock),
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
		Description:   helper.StringToNull(req.Description),
		Price:         helper.Float64ToDecimal(req.Price),
		CategoryID:    req.CategoryID,
		StockQuantity: int32(req.StockQuantity),
		IsActive:      helper.BoolPtrValue(req.IsActive, true),
	}

	if err := s.repo.Update(ctx, params); err != nil {
		return ProductResponse{}, err
	}

	dashboardCacheKey := dashboard.ProductReportKey
	if err := s.rdb.Del(ctx, dashboardCacheKey).Err(); err != nil {
		log.Printf("failed to invalidate dashboard product cache: %v", err)
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
		Price:         helper.DecimalToFloat64(r.Price),
		StockQuantity: int(r.StockQuantity),
		TotalSold:     int(r.TotalSold),
		IsActive:      r.IsActive,
		Category: CategoryResponse{
			ID:          r.CategoryID,
			Name:        r.CategoryName,
			Description: r.CategoryDescription.String,
		},
	}
}

// Mapper khusus untuk hasil GetByID
func mapDetailToResponse(r dbgen.GetProductByIDRow) ProductResponse {
	return ProductResponse{
		ID:            r.ID,
		Name:          r.Name,
		Description:   r.Description.String,
		Price:         helper.DecimalToFloat64(r.Price),
		StockQuantity: int(r.StockQuantity),
		IsActive:      r.IsActive,
		Category: CategoryResponse{
			ID:          r.CategoryID,
			Name:        r.CategoryName,
			Description: r.CategoryDescription.String,
		},
	}
}

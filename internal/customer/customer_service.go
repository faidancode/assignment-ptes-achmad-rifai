package customer

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=customer_service.go -destination=mocks/customer_service_mock.go -package=mock
type Service interface {
	Create(ctx context.Context, req CreateCustomerRequest) (CustomerResponse, error)
	List(ctx context.Context, p ListParams) ([]CustomerResponse, error)
	GetByID(ctx context.Context, id string) (CustomerResponse, error)
	Update(ctx context.Context, id string, req UpdateCustomerRequest) (CustomerResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	repo Repository
}

func NewService(repo Repository) Service {
	return &service{repo: repo}
}

func (s *service) Create(ctx context.Context, req CreateCustomerRequest) (CustomerResponse, error) {
	id := uuid.NewString()
	now := time.Now()

	params := dbgen.CreateCustomerParams{
		ID:        id,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
	}

	if err := s.repo.Create(ctx, params); err != nil {
		return CustomerResponse{}, err
	}

	return CustomerResponse{
		ID:        id,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: now,
	}, nil
}

func (s *service) List(ctx context.Context, p ListParams) ([]CustomerResponse, error) {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.PageSize <= 0 {
		p.PageSize = 10
	}

	limit := int32(p.PageSize)
	offset := int32((p.Page - 1) * p.PageSize)
	rows, err := s.repo.GetCustomers(ctx, dbgen.GetCustomersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	res := make([]CustomerResponse, 0, len(rows))
	for _, row := range rows {
		res = append(res, mapToListResponse(row))
	}
	return res, nil
}

func (s *service) GetByID(ctx context.Context, id string) (CustomerResponse, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return CustomerResponse{}, ErrCustomerNotFound
	}
	return mapToDetailResponse(row), nil
}

func (s *service) Update(ctx context.Context, id string, req UpdateCustomerRequest) (CustomerResponse, error) {
	// Check existence
	existing, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return CustomerResponse{}, ErrCustomerNotFound
	}

	params := dbgen.UpdateCustomerParams{
		ID:    id,
		Name:  req.Name,
		Email: req.Email,
	}

	if err := s.repo.Update(ctx, params); err != nil {
		return CustomerResponse{}, err
	}

	return CustomerResponse{
		ID:        id,
		Name:      req.Name,
		Email:     req.Email,
		CreatedAt: existing.CreatedAt,
	}, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

func mapToListResponse(row dbgen.GetCustomersRow) CustomerResponse {
	return CustomerResponse{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
	}
}

func mapToDetailResponse(row dbgen.GetCustomerByIDRow) CustomerResponse {
	return CustomerResponse{
		ID:        row.ID,
		Name:      row.Name,
		Email:     row.Email,
		CreatedAt: row.CreatedAt,
	}
}

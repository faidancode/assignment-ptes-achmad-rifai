package order

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"assignment-ptes-achmad-rifai/internal/shared/database/helper"
	"context"
	"database/sql"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=order_service.go -destination=mocks/order_service_mock.go -package=mock

type Service interface {
	Create(ctx context.Context, req CreateOrderRequest) (OrderResponse, error)
	List(ctx context.Context) ([]OrderResponse, error)
	GetByID(ctx context.Context, id string) (OrderResponse, error)
	Delete(ctx context.Context, id string) error
}

type service struct {
	db   *sql.DB // Diperlukan untuk memulai transaksi
	repo Repository
}

func NewService(db *sql.DB, repo Repository) Service {
	return &service{
		db:   db,
		repo: repo,
	}
}

func (s *service) Create(ctx context.Context, req CreateOrderRequest) (OrderResponse, error) {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return OrderResponse{}, err
	}
	defer tx.Rollback()

	txRepo := s.repo.WithTx(tx)

	orderID := uuid.NewString()
	now := time.Now()
	var totalQty int
	var totalPrice float64

	for _, item := range req.Items {
		totalQty += item.Quantity
		totalPrice += float64(item.Quantity) * item.UnitPrice
	}

	orderParams := dbgen.CreateOrderParams{
		ID:            orderID,
		CustomerID:    req.CustomerID,
		TotalQuantity: int32(totalQty),
		TotalPrice:    helper.Float64ToDecimal(totalPrice),
		CreatedAt:     now,
	}

	if err := txRepo.CreateOrder(ctx, orderParams); err != nil {
		return OrderResponse{}, err
	}

	itemResponses := make([]OrderItemResponse, 0)
	for _, item := range req.Items {
		itemID := uuid.NewString()
		itemParams := dbgen.CreateOrderItemParams{
			ID:        itemID,
			OrderID:   orderID,
			ProductID: item.ProductID,
			Quantity:  int32(item.Quantity),
			UnitPrice: helper.Float64ToDecimal(item.UnitPrice),
		}

		if err := txRepo.CreateOrderItem(ctx, itemParams); err != nil {
			return OrderResponse{}, err
		}

		itemResponses = append(itemResponses, OrderItemResponse{
			ID: itemID, ProductID: item.ProductID, Quantity: item.Quantity, UnitPrice: item.UnitPrice,
		})
	}

	if err := tx.Commit(); err != nil {
		return OrderResponse{}, err
	}

	return OrderResponse{
		ID:            orderID,
		CustomerID:    req.CustomerID,
		TotalQuantity: totalQty,
		TotalPrice:    totalPrice,
		CreatedAt:     now,
		Items:         itemResponses,
	}, nil
}

func (s *service) List(ctx context.Context) ([]OrderResponse, error) {
	rows, err := s.repo.GetOrders(ctx)
	if err != nil {
		return nil, err
	}

	res := make([]OrderResponse, 0, len(rows))
	for _, row := range rows {
		res = append(res, OrderResponse{
			ID:            row.ID,
			CustomerID:    row.CustomerID,
			TotalQuantity: int(row.TotalQuantity),
			TotalPrice:    helper.DecimalToFloat64(row.TotalPrice),
			CreatedAt:     row.CreatedAt,
		})
	}
	return res, nil
}

func (s *service) GetByID(ctx context.Context, id string) (OrderResponse, error) {
	row, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	items, err := s.repo.GetItemsByOrderID(ctx, id)
	if err != nil {
		// ERROR CEK DI SINI: Pastikan error ini dikembalikan (return)
		return OrderResponse{}, err
	}
	itemResponses := make([]OrderItemResponse, 0)
	for _, item := range items {
		itemResponses = append(itemResponses, OrderItemResponse{
			ID: item.ID, ProductID: item.ProductID, Quantity: int(item.Quantity), UnitPrice: helper.DecimalToFloat64(item.UnitPrice),
		})
	}

	return OrderResponse{
		ID:            row.ID,
		CustomerID:    row.CustomerID,
		TotalQuantity: int(row.TotalQuantity),
		TotalPrice:    helper.DecimalToFloat64(row.TotalPrice),
		CreatedAt:     row.CreatedAt,
		Items:         itemResponses,
	}, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

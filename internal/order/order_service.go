package order

import (
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"assignment-ptes-achmad-rifai/internal/shared/database/helper"
	"context"
	"database/sql"
	"encoding/json"
	"log"
	"time"

	"github.com/google/uuid"
)

//go:generate mockgen -source=order_service.go -destination=mocks/order_service_mock.go -package=mock

type Service interface {
	Create(ctx context.Context, req CreateOrderRequest) (OrderResponse, error)
	List(ctx context.Context, params ListParams) ([]OrderResponse, error)
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
		TotalQuantity: int32(totalQty),
		TotalPrice:    totalPrice,
		CreatedAt:     now,
		Items:         itemResponses,
	}, nil
}

func (s *service) List(ctx context.Context, p ListParams) ([]OrderResponse, error) {
	limit := int32(p.PageSize)
	offset := int32((p.Page - 1) * p.PageSize)
	rows, err := s.repo.GetOrders(ctx, dbgen.GetOrdersParams{
		Limit:  limit,
		Offset: offset,
	})
	if err != nil {
		return nil, err
	}

	var resp []OrderResponse
	for _, r := range rows {
		var items []OrderItemResponse
		if len(r.Items) > 0 {
			if err := json.Unmarshal(r.Items, &items); err != nil {
				log.Printf("error unmarshal items for order %s: %v", r.ID, err)
			}
		}

		totalPrice, _ := r.TotalPrice.Float64()

		resp = append(resp, OrderResponse{
			ID:            r.ID,
			CustomerID:    r.CustomerID,
			CustomerName:  r.CustomerName,
			CustomerEmail: r.CustomerEmail,
			TotalQuantity: r.TotalQuantity,
			TotalPrice:    totalPrice,
			CreatedAt:     r.CreatedAt,
			Items:         items,
		})
	}
	return resp, nil
}

func (s *service) GetByID(ctx context.Context, id string) (OrderResponse, error) {
	r, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return OrderResponse{}, err
	}

	var items []OrderItemResponse
	if len(r.Items) > 0 {
		if err := json.Unmarshal(r.Items, &items); err != nil {
			log.Printf("error unmarshal items for order %s: %v", r.ID, err)
		}
	}

	return OrderResponse{
		ID:            r.ID,
		CustomerID:    r.CustomerID,
		CustomerName:  r.CustomerName,
		CustomerEmail: r.CustomerEmail,
		TotalQuantity: int32(r.TotalQuantity),
		TotalPrice:    helper.DecimalToFloat64(r.TotalPrice),
		CreatedAt:     r.CreatedAt,
		Items:         items,
	}, nil
}

func (s *service) Delete(ctx context.Context, id string) error {
	return s.repo.Delete(ctx, id)
}

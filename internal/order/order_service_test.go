package order_test

import (
	"context"
	"database/sql"
	"errors"
	"testing"
	"time"

	"assignment-ptes-achmad-rifai/internal/order"
	mockOrder "assignment-ptes-achmad-rifai/internal/order/mocks"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/google/uuid"
	"github.com/shopspring/decimal"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupServiceTest(t *testing.T) (order.Service, *mockOrder.MockRepository, sqlmock.Sqlmock) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}

	t.Cleanup(func() {
		db.Close()
	})

	repo := mockOrder.NewMockRepository(ctrl)
	svc := order.NewService(db, repo)

	return svc, repo, mock
}

func TestService_Create_WithTransaction(t *testing.T) {
	ctx := context.Background()

	t.Run("success_create_order", func(t *testing.T) {
		svc, repo, mock := setupServiceTest(t) // Ambil mock dari setup

		customerID := uuid.NewString()
		productID := uuid.NewString()
		req := order.CreateOrderRequest{
			CustomerID: customerID,
			Items: []order.OrderItemRequest{
				{ProductID: productID, Quantity: 2, UnitPrice: 50000},
			},
		}

		// --- SQL Mock Expectations ---
		mock.ExpectBegin()
		mock.ExpectCommit()

		// --- Repo Mock Expectations ---
		// Kita butuh transaksi dummy untuk dikembalikan oleh BeginTx mock
		dbTmp, _, _ := sqlmock.New()
		defer dbTmp.Close()

		repo.EXPECT().WithTx(gomock.Any()).Return(repo).AnyTimes()
		repo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)
		repo.EXPECT().CreateOrderItem(gomock.Any(), gomock.Any()).Return(nil)

		// Execute
		res, err := svc.Create(ctx, req)

		// Assert
		assert.NoError(t, err)
		assert.Equal(t, 2, res.TotalQuantity)
		assert.Equal(t, float64(100000), res.TotalPrice)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error_create_item_failed_should_rollback", func(t *testing.T) {
		svc, repo, mock := setupServiceTest(t)

		req := order.CreateOrderRequest{
			CustomerID: uuid.NewString(),
			Items:      []order.OrderItemRequest{{ProductID: "p1", Quantity: 1, UnitPrice: 100}},
		}

		// --- SQL Mock Expectations ---
		mock.ExpectBegin()
		mock.ExpectRollback()

		dbTmp, _, _ := sqlmock.New()
		defer dbTmp.Close()

		repo.EXPECT().WithTx(gomock.Any()).Return(repo).AnyTimes()
		repo.EXPECT().CreateOrder(gomock.Any(), gomock.Any()).Return(nil)

		// Simulasi error pada item
		repo.EXPECT().CreateOrderItem(gomock.Any(), gomock.Any()).Return(assert.AnError)

		_, err := svc.Create(ctx, req)

		assert.Error(t, err)
		assert.NoError(t, mock.ExpectationsWereMet())
	})

	t.Run("error_begin_tx_failed", func(t *testing.T) {
		svc, _, _ := setupServiceTest(t)

		req := order.CreateOrderRequest{
			CustomerID: "c1",
			Items:      []order.OrderItemRequest{{ProductID: "p1", Quantity: 1, UnitPrice: 10}},
		}

		_, err := svc.Create(ctx, req)

		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		// Sesuaikan dengan setupServiceTest yang mengembalikan (svc, repo, mock)
		svc, repo, _ := setupServiceTest(t)

		rows := []dbgen.Order{
			{
				ID:            "o1",
				CustomerID:    "c1",
				TotalQuantity: 1,
				// Pastikan menggunakan decimal.Decimal untuk field TotalPrice
				TotalPrice: decimal.NewFromFloat(150000),
				CreatedAt:  time.Now(),
			},
		}

		repo.EXPECT().GetOrders(ctx).Return(rows, nil)

		res, err := svc.List(ctx)
		assert.NoError(t, err)
		assert.Len(t, res, 1)
		assert.Equal(t, "o1", res[0].ID)
		assert.Equal(t, float64(150000), res[0].TotalPrice)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)

		repo.EXPECT().GetOrders(ctx).Return(nil, errors.New("db error"))

		_, err := svc.List(ctx)
		assert.Error(t, err)
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)

		// Mock Data Parent Order
		repo.EXPECT().GetByID(ctx, id).Return(dbgen.Order{
			ID:         id,
			TotalPrice: decimal.NewFromFloat(200000),
		}, nil)

		// Mock Data Child Items
		repo.EXPECT().GetItemsByOrderID(ctx, id).Return([]dbgen.OrderItem{
			{
				ID:        uuid.NewString(),
				OrderID:   id,
				UnitPrice: decimal.NewFromFloat(100000),
				Quantity:  2,
			},
		}, nil)

		res, err := svc.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
		assert.Len(t, res.Items, 1)
		assert.Equal(t, float64(200000), res.TotalPrice)
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)

		// Sqlc biasanya mengembalikan sql.ErrNoRows jika tidak ketemu
		repo.EXPECT().GetByID(ctx, id).Return(dbgen.Order{}, sql.ErrNoRows)

		_, err := svc.GetByID(ctx, id)
		assert.Error(t, err)
	})

	t.Run("repo error on items", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)

		repo.EXPECT().GetByID(ctx, id).Return(dbgen.Order{ID: id}, nil)
		// Simulasi error saat mengambil detail items
		repo.EXPECT().GetItemsByOrderID(ctx, id).Return(nil, errors.New("db error"))

		_, err := svc.GetByID(ctx, id)
		assert.Error(t, err)
	})
}

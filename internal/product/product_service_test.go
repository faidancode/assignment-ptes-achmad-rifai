// products/product.service_test.go

package product_test

import (
	"assignment-ptes-achmad-rifai/internal/dashboard"
	"assignment-ptes-achmad-rifai/internal/product"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
	"database/sql"
	"errors"
	"testing"

	mockProduct "assignment-ptes-achmad-rifai/internal/product/mocks"

	"github.com/go-redis/redismock/v9"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func setupServiceTest(t *testing.T) (product.Service, *mockProduct.MockRepository, redismock.ClientMock) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	// Mock Redis
	dbRedis, redisMock := redismock.NewClientMock()

	// Mock Repo
	repo := mockProduct.NewMockRepository(ctrl)

	// Create Service
	svc := product.NewService(repo, dbRedis)

	return svc, repo, redisMock
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()
	req := product.CreateProductRequest{
		Name:  "Indomie",
		Price: 3500,
	}

	t.Run("success", func(t *testing.T) {
		svc, repo, redisMock := setupServiceTest(t)
		repo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(nil)

		redisMock.ExpectDel(dashboard.ProductReportKey).SetVal(1)

		res, err := svc.Create(ctx, req)
		assert.NoError(t, err)
		assert.Equal(t, req.Name, res.Name)
	})

	t.Run("error database", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		repo.EXPECT().
			Create(gomock.Any(), gomock.Any()).
			Return(errors.New("db error"))

		_, err := svc.Create(ctx, req)
		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	ctx := context.Background()
	p := product.ListParams{Page: 1, PageSize: 10}

	t.Run("success", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]dbgen.ListProductsRow{{ID: "1", Name: "P1"}}, nil)
		repo.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(1), nil)

		res, total, err := svc.List(ctx, p)
		assert.NoError(t, err)
		assert.Equal(t, int64(1), total)
		assert.Len(t, res, 1)
	})

	t.Run("error count", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		repo.EXPECT().List(gomock.Any(), gomock.Any()).Return([]dbgen.ListProductsRow{}, nil)
		repo.EXPECT().Count(gomock.Any(), gomock.Any()).Return(int64(0), errors.New("count failed"))

		_, _, err := svc.List(ctx, p)
		assert.Error(t, err)
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	id := "uuid-1"

	t.Run("success", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		repo.EXPECT().GetByID(ctx, id).Return(dbgen.GetProductByIDRow{ID: id, Name: "P1"}, nil)

		res, err := svc.GetByID(ctx, id)
		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		repo.EXPECT().GetByID(ctx, id).Return(dbgen.GetProductByIDRow{}, sql.ErrNoRows)

		_, err := svc.GetByID(ctx, id)
		assert.ErrorIs(t, err, product.ErrProductNotFound)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	id := "uuid-1"
	req := product.UpdateProductRequest{Name: "New Name"}

	t.Run("success", func(t *testing.T) {
		svc, repo, redisMock := setupServiceTest(t)
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		repo.EXPECT().GetByID(gomock.Any(), id).Return(dbgen.GetProductByIDRow{ID: id, Name: "New Name"}, nil)
		redisMock.ExpectDel(dashboard.ProductReportKey).SetVal(1)
		res, err := svc.Update(ctx, id, req)
		assert.NoError(t, err)
		assert.Equal(t, "New Name", res.Name)
	})

	t.Run("not found after update", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		repo.EXPECT().Update(gomock.Any(), gomock.Any()).Return(nil)
		repo.EXPECT().GetByID(gomock.Any(), id).Return(dbgen.GetProductByIDRow{}, sql.ErrNoRows)

		_, err := svc.Update(ctx, id, req)
		assert.ErrorIs(t, err, product.ErrProductNotFound)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	id := "uuid-1"

	t.Run("success", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		repo.EXPECT().Delete(ctx, id).Return(nil)

		err := svc.Delete(ctx, id)
		assert.NoError(t, err)
	})

	t.Run("failed delete", func(t *testing.T) {
		svc, repo, _ := setupServiceTest(t)
		repo.EXPECT().Delete(ctx, id).Return(errors.New("constraint error"))

		err := svc.Delete(ctx, id)
		assert.Error(t, err)
	})
}

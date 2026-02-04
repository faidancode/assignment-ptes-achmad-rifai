package customer_test

import (
	"assignment-ptes-achmad-rifai/internal/customer"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"context"
	"errors"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mockCustomer "assignment-ptes-achmad-rifai/internal/customer/mocks"
)

func setupServiceTest(t *testing.T) (customer.Service, *mockCustomer.MockRepository) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mockCustomer.NewMockRepository(ctrl)

	svc := customer.NewService(repo)

	return svc, repo
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		req := customer.CreateCustomerRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		// Expect Create called with matching params
		repo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(dbgen.CreateCustomerParams{})).
			DoAndReturn(func(_ context.Context, p dbgen.CreateCustomerParams) error {
				assert.NotEmpty(t, p.ID)
				assert.Equal(t, "John Doe", p.Name)
				assert.Equal(t, "john@example.com", p.Email)
				assert.WithinDuration(t, time.Now(), p.CreatedAt, time.Second*5)
				return nil
			})

		res, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "John Doe", res.Name)
		assert.Equal(t, "john@example.com", res.Email)
		assert.NotEmpty(t, res.ID)
		assert.WithinDuration(t, time.Now(), res.CreatedAt, time.Second*5)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		req := customer.CreateCustomerRequest{
			Name:  "John Doe",
			Email: "john@example.com",
		}

		repo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(dbgen.CreateCustomerParams{})).
			Return(errors.New("db error"))

		_, err := svc.Create(ctx, req)

		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	ctx := context.Background()
	p := customer.ListParams{Page: 1, PageSize: 10}

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		rows := []dbgen.GetCustomersRow{
			{
				ID:        uuid.NewString(),
				Name:      "User 1",
				Email:     "user1@example.com",
				CreatedAt: time.Now(),
			},
			{
				ID:        uuid.NewString(),
				Name:      "User 2",
				Email:     "user2@example.com",
				CreatedAt: time.Now(),
			},
		}

		repo.EXPECT().
			GetCustomers(ctx, p).
			Return(rows, nil)

		res, err := svc.List(ctx, p)

		assert.NoError(t, err)
		assert.Len(t, res, 2)
		assert.Equal(t, "User 1", res[0].Name)
		assert.Equal(t, "User 2", res[1].Name)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)
		p := customer.ListParams{Page: 1, PageSize: 10}

		repo.EXPECT().
			GetCustomers(ctx, p).
			Return(nil, errors.New("db error"))

		_, err := svc.List(ctx, p)

		assert.Error(t, err)
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		row := dbgen.GetCustomerByIDRow{
			ID:        id,
			Name:      "Jane Doe",
			Email:     "jane@example.com",
			CreatedAt: time.Now(),
		}

		repo.EXPECT().
			GetByID(ctx, id).
			Return(row, nil)

		res, err := svc.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
		assert.Equal(t, "Jane Doe", res.Name)
		assert.Equal(t, "jane@example.com", res.Email)
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			GetByID(ctx, id).
			Return(dbgen.GetCustomerByIDRow{}, customer.ErrCustomerNotFound)

		_, err := svc.GetByID(ctx, id)

		assert.ErrorIs(t, err, customer.ErrCustomerNotFound)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			GetByID(ctx, id).
			Return(dbgen.GetCustomerByIDRow{}, errors.New("db error"))

		_, err := svc.GetByID(ctx, id)

		assert.Error(t, err)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		req := customer.UpdateCustomerRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		existing := dbgen.GetCustomerByIDRow{
			ID:        id,
			Name:      "Old Name",
			Email:     "old@example.com",
			CreatedAt: time.Now().Add(-time.Hour),
		}

		// Expect GetByID for existence check
		repo.EXPECT().
			GetByID(ctx, id).
			Return(existing, nil)

		// Expect Update call
		repo.EXPECT().
			Update(gomock.Any(), gomock.AssignableToTypeOf(dbgen.UpdateCustomerParams{})).
			DoAndReturn(func(_ context.Context, p dbgen.UpdateCustomerParams) error {
				assert.Equal(t, id, p.ID)
				assert.Equal(t, "Updated Name", p.Name)
				assert.Equal(t, "updated@example.com", p.Email)
				return nil
			})

		res, err := svc.Update(ctx, id, req)

		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
		assert.Equal(t, "Updated Name", res.Name)
		assert.Equal(t, "updated@example.com", res.Email)
		assert.Equal(t, existing.CreatedAt, res.CreatedAt)
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		req := customer.UpdateCustomerRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		repo.EXPECT().
			GetByID(ctx, id).
			Return(dbgen.GetCustomerByIDRow{}, customer.ErrCustomerNotFound)

		_, err := svc.Update(ctx, id, req)

		assert.ErrorIs(t, err, customer.ErrCustomerNotFound)
	})

	t.Run("update error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		req := customer.UpdateCustomerRequest{
			Name:  "Updated Name",
			Email: "updated@example.com",
		}

		existing := dbgen.GetCustomerByIDRow{
			ID:        id,
			Name:      "Old Name",
			Email:     "old@example.com",
			CreatedAt: time.Now().Add(-time.Hour),
		}

		repo.EXPECT().
			GetByID(ctx, id).
			Return(existing, nil)

		repo.EXPECT().
			Update(gomock.Any(), gomock.AssignableToTypeOf(dbgen.UpdateCustomerParams{})).
			Return(errors.New("update error"))

		_, err := svc.Update(ctx, id, req)

		assert.Error(t, err)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	id := uuid.NewString()

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			Delete(ctx, id).
			Return(nil)

		err := svc.Delete(ctx, id)

		assert.NoError(t, err)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			Delete(ctx, id).
			Return(errors.New("db error"))

		err := svc.Delete(ctx, id)

		assert.Error(t, err)
	})
}

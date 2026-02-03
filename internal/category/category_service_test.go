package category_test

import (
	"assignment-ptes-achmad-rifai/internal/category"
	"assignment-ptes-achmad-rifai/internal/shared/database/dbgen"
	"assignment-ptes-achmad-rifai/internal/shared/database/helper"
	"context"
	"errors"
	"testing"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"

	mockCategory "assignment-ptes-achmad-rifai/internal/category/mocks"
)

func setupServiceTest(t *testing.T) (category.Service, *mockCategory.MockRepository) {
	ctrl := gomock.NewController(t)
	t.Cleanup(ctrl.Finish)

	repo := mockCategory.NewMockRepository(ctrl)

	svc := category.NewService(repo)

	return svc, repo
}

func TestService_Create(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		desc := "Food Category"
		req := category.CreateCategoryRequest{
			Name:        "Food",
			Description: &desc,
		}

		// Only expect Create
		repo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(dbgen.CreateCategoryParams{})).
			DoAndReturn(func(_ context.Context, p dbgen.CreateCategoryParams) error {
				assert.NotEmpty(t, p.ID)
				assert.Equal(t, "Food", p.Name)
				assert.True(t, p.Description.Valid)
				assert.Equal(t, "Food Category", p.Description.String)
				return nil
			})

		res, err := svc.Create(ctx, req)

		assert.NoError(t, err)
		assert.Equal(t, "Food", res.Name)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		desc := "Food Category"
		req := category.CreateCategoryRequest{
			Name:        "Food",
			Description: &desc,
		}

		repo.EXPECT().
			Create(gomock.Any(), gomock.AssignableToTypeOf(dbgen.CreateCategoryParams{})).
			Return(errors.New("db error"))

		_, err := svc.Create(ctx, req)

		assert.Error(t, err)
	})
}

func TestService_List(t *testing.T) {
	ctx := context.Background()

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			GetCategories(ctx).
			Return([]dbgen.GetCategoriesRow{
				{ID: "1", Name: "Food"},
				{ID: "2", Name: "Drink"},
			}, nil)

		res, err := svc.List(ctx)

		assert.NoError(t, err)
		assert.Len(t, res, 2)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			GetCategories(ctx).
			Return(nil, errors.New("db error"))

		_, err := svc.List(ctx)

		assert.Error(t, err)
	})
}

func TestService_GetByID(t *testing.T) {
	ctx := context.Background()
	id := "uuid-1"

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			GetByID(ctx, id).
			Return(dbgen.GetCategoryByIDRow{
				ID:   id,
				Name: "Food",
			}, nil)

		res, err := svc.GetByID(ctx, id)

		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
	})

	t.Run("not found", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			GetByID(ctx, id).
			Return(dbgen.GetCategoryByIDRow{}, category.ErrCategoryNotFound)

		_, err := svc.GetByID(ctx, id)

		assert.ErrorIs(t, err, category.ErrCategoryNotFound)
	})

	t.Run("repo error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		repo.EXPECT().
			GetByID(ctx, id).
			Return(dbgen.GetCategoryByIDRow{}, errors.New("db error"))

		_, err := svc.GetByID(ctx, id)

		assert.Error(t, err)
	})
}

func TestService_Update(t *testing.T) {
	ctx := context.Background()
	id := "uuid-1"

	t.Run("success", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		desc := "Updated desc"
		req := category.UpdateCategoryRequest{
			Name:        "Updated",
			Description: &desc,
		}

		// Expect Update
		repo.EXPECT().
			Update(gomock.Any(), gomock.AssignableToTypeOf(dbgen.UpdateCategoryParams{})).
			DoAndReturn(func(_ context.Context, p dbgen.UpdateCategoryParams) error {
				assert.Equal(t, id, p.ID)
				assert.Equal(t, "Updated", p.Name)
				assert.True(t, p.Description.Valid)
				assert.Equal(t, "Updated desc", p.Description.String)
				return nil
			})

		// Expect GetByID (dipanggil setelah Update)
		repo.EXPECT().
			GetByID(gomock.Any(), id).
			Return(dbgen.GetCategoryByIDRow{
				ID:          id,
				Name:        "Updated",
				Description: helper.NewNullString(&desc),
			}, nil)

		res, err := svc.Update(ctx, id, req)

		assert.NoError(t, err)
		assert.Equal(t, id, res.ID)
		assert.Equal(t, "Updated", res.Name)
	})

	t.Run("update error", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		desc := "Updated desc"
		req := category.UpdateCategoryRequest{
			Name:        "Updated",
			Description: &desc,
		}

		repo.EXPECT().
			Update(gomock.Any(), gomock.AssignableToTypeOf(dbgen.UpdateCategoryParams{})).
			Return(errors.New("db error"))

		_, err := svc.Update(ctx, id, req)

		assert.Error(t, err)
	})

	t.Run("not found after update", func(t *testing.T) {
		svc, repo := setupServiceTest(t)

		desc := "Updated desc"
		req := category.UpdateCategoryRequest{
			Name:        "Updated",
			Description: &desc,
		}

		repo.EXPECT().
			Update(gomock.Any(), gomock.AssignableToTypeOf(dbgen.UpdateCategoryParams{})).
			Return(nil)

		repo.EXPECT().
			GetByID(gomock.Any(), id).
			Return(dbgen.GetCategoryByIDRow{}, category.ErrCategoryNotFound)

		_, err := svc.Update(ctx, id, req)

		assert.ErrorIs(t, err, category.ErrCategoryNotFound)
	})
}

func TestService_Delete(t *testing.T) {
	ctx := context.Background()
	id := "uuid-1"

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

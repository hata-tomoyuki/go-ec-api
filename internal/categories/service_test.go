package categories

import (
	"context"
	"errors"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

func TestCreateCategories_Success(t *testing.T) {
	mock := &mockQuerier{
		createCategoryFn: func(ctx context.Context, arg repo.CreateCategoryParams) (repo.Category, error) {
			return newTestCategory(1, arg.Name), nil
		},
	}

	svc := NewService(mock)
	category, err := svc.CreateCategories(context.Background(), "Electronics", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if category.Name != "Electronics" {
		t.Errorf("expected name 'Electronics', got '%s'", category.Name)
	}
	if category.ID != 1 {
		t.Errorf("expected id=1, got %d", category.ID)
	}
}

func TestFindCategoryById_Success(t *testing.T) {
	mock := &mockQuerier{
		findCategoryByIdFn: func(ctx context.Context, id int64) (repo.Category, error) {
			return newTestCategory(1, "Electronics"), nil
		},
	}

	svc := NewService(mock)
	category, err := svc.FindCategoryById(context.Background(), 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if category.ID != 1 {
		t.Errorf("expected id=1, got %d", category.ID)
	}
}

func TestFindCategoryById_NotFound(t *testing.T) {
	mock := &mockQuerier{
		findCategoryByIdFn: func(ctx context.Context, id int64) (repo.Category, error) {
			return repo.Category{}, pgx.ErrNoRows
		},
	}

	svc := NewService(mock)
	_, err := svc.FindCategoryById(context.Background(), 999)

	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestUpdateCategories_Success(t *testing.T) {
	mock := &mockQuerier{
		updateCategoryFn: func(ctx context.Context, arg repo.UpdateCategoryParams) (repo.Category, error) {
			return newTestCategory(arg.ID, arg.Name), nil
		},
	}

	svc := NewService(mock)
	category, err := svc.UpdateCategories(context.Background(), 1, "Updated Electronics", nil)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if category.Name != "Updated Electronics" {
		t.Errorf("expected name 'Updated Electronics', got '%s'", category.Name)
	}
}

func TestUpdateCategories_NotFound(t *testing.T) {
	mock := &mockQuerier{
		updateCategoryFn: func(ctx context.Context, arg repo.UpdateCategoryParams) (repo.Category, error) {
			return repo.Category{}, pgx.ErrNoRows
		},
	}

	svc := NewService(mock)
	_, err := svc.UpdateCategories(context.Background(), 999, "Updated Electronics", nil)

	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}

func TestDeleteCategory_Success(t *testing.T) {
	mock := &mockQuerier{
		deleteCategoryFn: func(ctx context.Context, id int64) (repo.Category, error) {
			return newTestCategory(id, "Electronics"), nil
		},
	}

	svc := NewService(mock)
	err := svc.DeleteCategory(context.Background(), 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteCategory_NotFound(t *testing.T) {
	mock := &mockQuerier{
		deleteCategoryFn: func(ctx context.Context, id int64) (repo.Category, error) {
			return repo.Category{}, pgx.ErrNoRows
		},
	}

	svc := NewService(mock)
	err := svc.DeleteCategory(context.Background(), 999)

	if !errors.Is(err, ErrCategoryNotFound) {
		t.Fatalf("expected ErrCategoryNotFound, got %v", err)
	}
}

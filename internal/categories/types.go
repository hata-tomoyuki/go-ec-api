package categories

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

var ErrCategoryNotFound = errors.New("category not found")

type Service interface {
	ListCategories(ctx context.Context) ([]repo.Category, error)
	FindCategoryById(ctx context.Context, id int64) (repo.Category, error)
	CreateCategories(ctx context.Context, name string, description *string) (repo.Category, error)
	UpdateCategories(ctx context.Context, id int64, name string, description *string) (repo.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	ListProductsByCategory(ctx context.Context, categoryId int64) ([]repo.Product, error)
	AddProductToCategory(ctx context.Context, categoryId int64, productId int64) error
	RemoveProductFromCategory(ctx context.Context, categoryId int64, productId int64) error
}

type createCategoryParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

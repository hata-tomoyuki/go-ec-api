package categories

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

var (
	ErrCategoryNotFound   = errors.New("category not found")
	ErrCategoryValidation = errors.New("name is required")
)

type Service interface {
	ListCategories(ctx context.Context) ([]repo.ListCategoriesRow, error)
	FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error)
	CreateCategories(ctx context.Context, name string, description *string, imageColor string) (repo.Category, error)
	UpdateCategories(ctx context.Context, id int64, name string, description *string, imageColor string) (repo.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	ListProductsByCategory(ctx context.Context, categoryId int64) ([]repo.ListProductsByCategoryRow, error)
	AddProductToCategory(ctx context.Context, categoryId int64, productId int64) error
	RemoveProductFromCategory(ctx context.Context, categoryId int64, productId int64) error
}

type createCategoryParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
	ImageColor  string  `json:"image_color"`
}

func (p createCategoryParams) validate() error {
	if p.Name == "" {
		return ErrCategoryValidation
	}
	return nil
}

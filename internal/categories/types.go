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
	ListCategoriesPaginated(ctx context.Context, params listCategoriesParams) (paginatedCategories, error)
	FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error)
	CreateCategories(ctx context.Context, name string, description *string, imageColor string) (repo.Category, error)
	UpdateCategories(ctx context.Context, id int64, name string, description *string, imageColor string) (repo.Category, error)
	DeleteCategory(ctx context.Context, id int64) error
	ListProductsByCategory(ctx context.Context, categoryId int64) ([]repo.ListProductsByCategoryRow, error)
	AddProductToCategory(ctx context.Context, categoryId int64, productId int64) error
	RemoveProductFromCategory(ctx context.Context, categoryId int64, productId int64) error
}

type listCategoriesParams struct {
	Page  int
	Limit int
}

func (p *listCategoriesParams) validate() error {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		return ErrCategoryValidation
	}
	return nil
}

type paginatedCategories struct {
	Data  []paginatedCategoryRow `json:"data"`
	Total int                    `json:"total"`
	Page  int                    `json:"page"`
	Limit int                    `json:"limit"`
}

type paginatedCategoryRow struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	Description  string `json:"description"`
	ImageColor   string `json:"image_color"`
	ProductCount int64  `json:"product_count"`
	TotalCount   int    `json:"-"`
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

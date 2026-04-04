package products

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

var (
	ErrProductNotFound  = errors.New("product not found")
	ErrProductValidation = errors.New("name is required and price_in_cents must be greater than 0")
)

type createProductParams struct {
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
	Description  string `json:"description"`
	ImageColor   string `json:"image_color"`
}

func (p createProductParams) validate() error {
	if p.Name == "" || p.PriceInCents <= 0 {
		return ErrProductValidation
	}
	return nil
}

type updateProductParams struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
	Description  string `json:"description"`
	ImageColor   string `json:"image_color"`
}

func (p updateProductParams) validate() error {
	if p.Name == "" || p.PriceInCents <= 0 {
		return ErrProductValidation
	}
	return nil
}

type Service interface {
	ListProducts(ctx context.Context) ([]repo.ListProductsRow, error)
	FindProductById(ctx context.Context, id int64) (repo.FindProductByIdRow, error)
	CreateProduct(ctx context.Context, tempProduct createProductParams) (repo.Product, error)
	UpdateProduct(ctx context.Context, tempProduct updateProductParams) (repo.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

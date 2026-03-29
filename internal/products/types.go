package products

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

var ErrProductNotFound = errors.New("product not found")

type createProductParams struct {
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
}

type updateProductParams struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
}

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductById(ctx context.Context, id int64) (repo.Product, error)
	CreateProduct(ctx context.Context, tempProduct createProductParams) (repo.Product, error)
	UpdateProduct(ctx context.Context, tempProduct updateProductParams) (repo.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

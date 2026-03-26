package products

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type createProductParams struct {
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
}

type Service interface {
	ListProducts(ctx context.Context) ([]repo.Product, error)
	FindProductById(ctx context.Context, id int64) (repo.Product, error)
	CreateProduct(ctx context.Context, tempProduct createProductParams) (repo.Product, error)
}

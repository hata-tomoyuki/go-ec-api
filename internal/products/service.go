package products

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}

}

func (s *svc) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return s.repo.ListProducts(ctx)
}

func (s *svc) FindProductById(ctx context.Context, id int64) (repo.Product, error) {
	return s.repo.FindProductById(ctx, id)
}

func (s *svc) CreateProduct(ctx context.Context, tempProduct createProductParams) (repo.Product, error) {
	return s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         tempProduct.Name,
		PriceInCents: tempProduct.PriceInCents,
	})
}

func (s *svc) UpdateProduct(ctx context.Context, tempProduct updateProductParams) (repo.Product, error) {
	return s.repo.UpdateProduct(ctx, repo.UpdateProductParams{
		ID:           tempProduct.ID,
		Name:         tempProduct.Name,
		PriceInCents: tempProduct.PriceInCents,
	})
}

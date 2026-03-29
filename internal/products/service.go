package products

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
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
	product, err := s.repo.FindProductById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Product{}, ErrProductNotFound
		}
		return repo.Product{}, err
	}
	return product, nil
}

func (s *svc) CreateProduct(ctx context.Context, tempProduct createProductParams) (repo.Product, error) {
	return s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         tempProduct.Name,
		PriceInCents: tempProduct.PriceInCents,
	})
}

func (s *svc) UpdateProduct(ctx context.Context, tempProduct updateProductParams) (repo.Product, error) {
	product, err := s.repo.UpdateProduct(ctx, repo.UpdateProductParams{
		ID:           tempProduct.ID,
		Name:         tempProduct.Name,
		PriceInCents: tempProduct.PriceInCents,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Product{}, ErrProductNotFound
		}
		return repo.Product{}, err
	}
	return product, nil
}

func (s *svc) DeleteProduct(ctx context.Context, id int64) error {
	_, err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

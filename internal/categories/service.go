package categories

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListCategories(ctx context.Context) ([]repo.Category, error) {
	return s.repo.ListCategories(ctx)
}

func (s *svc) CreateCategories(ctx context.Context, name string, description *string) (repo.Category, error) {
	desc := pgtype.Text{Valid: false}
	if description != nil {
		desc = pgtype.Text{
			String: *description,
			Valid:  true,
		}
	}

	return s.repo.CreateCategory(ctx, repo.CreateCategoryParams{
		Name:        name,
		Description: desc,
	})
}

func (s *svc) FindCategoryById(ctx context.Context, id int64) (repo.Category, error) {
	category, err := s.repo.FindCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Category{}, ErrCategoryNotFound
		}
		return repo.Category{}, err
	}
	return category, nil
}

func (s *svc) UpdateCategories(ctx context.Context, id int64, name string, description *string) (repo.Category, error) {
	desc := pgtype.Text{Valid: false}
	if description != nil {
		desc = pgtype.Text{
			String: *description,
			Valid:  true,
		}
	}

	category, err := s.repo.UpdateCategory(ctx, repo.UpdateCategoryParams{
		ID:          id,
		Name:        name,
		Description: desc,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Category{}, ErrCategoryNotFound
		}
		return repo.Category{}, err
	}
	return category, nil
}

func (s *svc) DeleteCategory(ctx context.Context, id int64) error {
	_, err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}
	return nil
}

func (s *svc) ListProductsByCategory(ctx context.Context, categoryId int64) ([]repo.Product, error) {
	return s.repo.ListProductsByCategory(ctx, categoryId)
}

func (s *svc) AddProductToCategory(ctx context.Context, categoryId int64, productId int64) error {
	return s.repo.AddProductToCategory(ctx, repo.AddProductToCategoryParams{
		ProductID:  productId,
		CategoryID: categoryId,
	})
}

func (s *svc) RemoveProductFromCategory(ctx context.Context, categoryId int64, productId int64) error {
	return s.repo.RemoveProductFromCategory(ctx, repo.RemoveProductFromCategoryParams{
		ProductID:  productId,
		CategoryID: categoryId,
	})
}

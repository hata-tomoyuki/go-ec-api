package categories

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
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
	return s.repo.FindCategoryById(ctx, id)
}

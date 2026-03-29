package categories

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type Service interface {
	CreateCategories(ctx context.Context, name string, description *string) (repo.Category, error)
}

type createCategoryParams struct {
	Name        string  `json:"name"`
	Description *string `json:"description"`
}

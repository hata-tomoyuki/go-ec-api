package auth

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type svc struct {
	repo *repo.Queries
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
	}
}

func (s *svc) RegisterUser(ctx context.Context, params registerParams) (repo.User, error) {
	hashedPassword, err := hashPassword(params.Password)
	if err != nil {
		return repo.User{}, err
	}

	return s.repo.CreateUser(ctx, repo.CreateUserParams{
		Email:        params.Email,
		PasswordHash: hashedPassword,
		Name:         params.Name,
	})
}

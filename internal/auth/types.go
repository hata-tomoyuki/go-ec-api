package auth

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type registerParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userResponse struct {
	ID    int64  `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

type Service interface {
	RegisterUser(ctx context.Context, params registerParams) (repo.User, error)
	Login(ctx context.Context, params loginParams) (string, error)
}

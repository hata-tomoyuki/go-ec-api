package auth

import (
	"context"
	"time"

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
	Role  string `json:"role"`
}

type updateUserParams struct {
	Name  *string `json:"name,omitempty"`
	Email *string `json:"email,omitempty"`
}

type Service interface {
	RegisterUser(ctx context.Context, params registerParams) (repo.User, error)
	Login(ctx context.Context, params loginParams) (string, error)
	Logout(ctx context.Context, jti string, expired_at time.Time) error
	UpdateUser(ctx context.Context, userID int64, params updateUserParams) (repo.User, error)
	UpdateUserPassword(ctx context.Context, userID int64, currentPassword, newPassword string) (repo.User, error)
}

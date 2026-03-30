package auth

import (
	"context"
	"errors"
	"time"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

var (
	ErrRegisterValidation   = errors.New("email, name, and password (min 8 chars) are required")
	ErrLoginValidation      = errors.New("email and password are required")
	ErrUpdateUserValidation = errors.New("at least one of name or email must be provided and not empty")
	ErrPasswordValidation   = errors.New("current password and new password (min 8 chars) are required")
)

type registerParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
	Name     string `json:"name"`
}

func (p registerParams) validate() error {
	if p.Email == "" || p.Name == "" || len(p.Password) < 8 {
		return ErrRegisterValidation
	}
	return nil
}

type loginParams struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (p loginParams) validate() error {
	if p.Email == "" || p.Password == "" {
		return ErrLoginValidation
	}
	return nil
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

func (p updateUserParams) validate() error {
	if p.Name == nil && p.Email == nil {
		return ErrUpdateUserValidation
	}
	if p.Name != nil && *p.Name == "" {
		return ErrUpdateUserValidation
	}
	if p.Email != nil && *p.Email == "" {
		return ErrUpdateUserValidation
	}
	return nil
}

type LoginTokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"`
}

type Service interface {
	RegisterUser(ctx context.Context, params registerParams) (repo.User, error)
	Login(ctx context.Context, params loginParams) (LoginTokens, error)
	Logout(ctx context.Context, jti string, expired_at time.Time, refreshTokenID int64) error
	Refresh(ctx context.Context, refreshTokenPlain string) (LoginTokens, error)
	GetProfile(ctx context.Context, userID int64) (repo.User, error)
	UpdateUser(ctx context.Context, userID int64, params updateUserParams) (repo.User, error)
	UpdateUserPassword(ctx context.Context, userID int64, currentPassword, newPassword string) (repo.User, error)
}

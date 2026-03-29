package auth

import (
	"context"
	"errors"
	"time"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

var ErrInvalidCredentials = errors.New("invalid email or password")

type svc struct {
	repo *repo.Queries
	ja   *jwtauth.JWTAuth
}

func NewService(repo *repo.Queries, ja *jwtauth.JWTAuth) Service {
	return &svc{
		repo: repo,
		ja:   ja,
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
		Role:         "user",
	})
}

func (s *svc) Login(ctx context.Context, params loginParams) (string, error) {
	user, err := s.repo.FindUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return "", ErrInvalidCredentials
		}
		return "", err
	}

	if err := checkPasswordHash(params.Password, user.PasswordHash); err != nil {
		return "", ErrInvalidCredentials
	}

	token, err := generateJWT(s.ja, user.ID, user.Name, user.Email, string(user.Role))
	if err != nil {
		return "", err
	}

	return token, nil
}

func (s *svc) Logout(ctx context.Context, jti string, expired_at time.Time) error {
	return s.repo.RevokeToken(ctx, repo.RevokeTokenParams{
		Jti: jti,
		ExpiredAt: pgtype.Timestamptz{
			Time:  expired_at,
			Valid: true,
		},
	})
}

func (s *svc) UpdateUser(ctx context.Context, userID int64, params updateUserParams) (repo.User, error) {
	current, err := s.repo.FindUserById(ctx, userID)
	if err != nil {
		return repo.User{}, err
	}

	name := current.Name
	if params.Name != nil {
		name = *params.Name
	}
	email := current.Email
	if params.Email != nil {
		email = *params.Email
	}

	return s.repo.UpdateUser(ctx, repo.UpdateUserParams{
		ID:    userID,
		Name:  name,
		Email: email,
	})
}

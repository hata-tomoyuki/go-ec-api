package auth

import (
	"context"
	"errors"
	"log/slog"
	"time"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type pgBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

var (
	ErrInvalidCredentials  = errors.New("invalid email or password")
	ErrInvalidRefreshToken = errors.New("invalid refresh token")
)

type svc struct {
	db     pgBeginner
	repo   repo.Querier
	ja     *jwtauth.JWTAuth
	newTxQ func(pgx.Tx) repo.Querier
}

func NewService(db pgBeginner, q repo.Querier, ja *jwtauth.JWTAuth) Service {
	return &svc{
		db:     db,
		repo:   q,
		ja:     ja,
		newTxQ: func(tx pgx.Tx) repo.Querier { return repo.New(tx) },
	}
}

func (s *svc) GetProfile(ctx context.Context, userID int64) (repo.User, error) {
	return s.repo.FindUserById(ctx, userID)
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

const refreshTokenTTL = 30 * 24 * time.Hour

func (s *svc) issueTokens(ctx context.Context, user repo.User) (LoginTokens, error) {
	if err := s.repo.DeleteExpiredRefreshTokens(ctx); err != nil {
		return LoginTokens{}, err
	}

	plain, err := newRefreshTokenPlaintext()
	if err != nil {
		return LoginTokens{}, err
	}
	row, err := s.repo.InsertRefreshToken(ctx, repo.InsertRefreshTokenParams{
		UserID:    user.ID,
		TokenHash: hashRefreshToken(plain),
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(refreshTokenTTL),
			Valid: true,
		},
	})
	if err != nil {
		return LoginTokens{}, err
	}
	access, err := generateJWT(s.ja, user.ID, user.Name, user.Email, string(user.Role), row.ID)
	if err != nil {
		if delErr := s.repo.DeleteRefreshToken(ctx, row.ID); delErr != nil {
			slog.ErrorContext(ctx, "failed to rollback refresh row after JWT encode error", "error", delErr)
		}
		return LoginTokens{}, err
	}
	return LoginTokens{
		AccessToken:  access,
		RefreshToken: plain,
		ExpiresIn:    int64(accessTokenTTL.Seconds()),
	}, nil
}

func (s *svc) Login(ctx context.Context, params loginParams) (LoginTokens, error) {
	user, err := s.repo.FindUserByEmail(ctx, params.Email)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LoginTokens{}, ErrInvalidCredentials
		}
		return LoginTokens{}, err
	}

	if err := checkPasswordHash(params.Password, user.PasswordHash); err != nil {
		return LoginTokens{}, ErrInvalidCredentials
	}

	return s.issueTokens(ctx, user)
}

func (s *svc) Refresh(ctx context.Context, refreshTokenPlain string) (LoginTokens, error) {
	if refreshTokenPlain == "" {
		return LoginTokens{}, ErrInvalidRefreshToken
	}
	row, err := s.repo.ConsumeRefreshToken(ctx, hashRefreshToken(refreshTokenPlain))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return LoginTokens{}, ErrInvalidRefreshToken
		}
		return LoginTokens{}, err
	}
	user, err := s.repo.FindUserById(ctx, row.UserID)
	if err != nil {
		return LoginTokens{}, err
	}
	return s.issueTokens(ctx, user)
}

func (s *svc) Logout(ctx context.Context, jti string, expired_at time.Time, refreshTokenID int64) error {
	tx, err := s.db.Begin(ctx)
	if err != nil {
		return err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.newTxQ(tx)

	if refreshTokenID > 0 {
		if err := qtx.DeleteRefreshToken(ctx, refreshTokenID); err != nil {
			return err
		}
	}
	if err := qtx.RevokeToken(ctx, repo.RevokeTokenParams{
		Jti: jti,
		ExpiredAt: pgtype.Timestamptz{
			Time:  expired_at,
			Valid: true,
		},
	}); err != nil {
		return err
	}
	return tx.Commit(ctx)
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

func (s *svc) UpdateUserPassword(ctx context.Context, userID int64, currentPassword, newPassword string) (repo.User, error) {
	user, err := s.repo.FindUserById(ctx, userID)
	if err != nil {
		return repo.User{}, err
	}

	if err := checkPasswordHash(currentPassword, user.PasswordHash); err != nil {
		return repo.User{}, ErrInvalidCredentials
	}

	hashedPassword, err := hashPassword(newPassword)
	if err != nil {
		return repo.User{}, err
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.User{}, err
	}
	defer func() { _ = tx.Rollback(ctx) }()

	qtx := s.newTxQ(tx)

	updated, err := qtx.UpdateUserPassword(ctx, repo.UpdateUserPasswordParams{
		ID:           userID,
		PasswordHash: hashedPassword,
	})
	if err != nil {
		return repo.User{}, err
	}

	if err := qtx.DeleteRefreshTokensByUserId(ctx, userID); err != nil {
		return repo.User{}, err
	}

	if err := tx.Commit(ctx); err != nil {
		return repo.User{}, err
	}

	return updated, nil
}

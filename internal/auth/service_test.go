package auth

import (
	"context"
	"errors"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
)

var testTokenAuth = jwtauth.New("HS256", []byte("test-secret"), nil)

func TestServiceGetProfile_Success(t *testing.T) {
	mock := &mockQuerier{
		findUserByIdFn: func(ctx context.Context, id int64) (repo.User, error) {
			return newTestUser(id, "Test User", "test@example.com"), nil
		},
	}
	svc := &svc{repo: mock, ja: testTokenAuth}

	user, err := svc.GetProfile(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.ID != 10 {
		t.Errorf("expected id=10, got %d", user.ID)
	}
}

func TestServiceGetProfile_NotFound(t *testing.T) {
	mock := &mockQuerier{
		findUserByIdFn: func(ctx context.Context, id int64) (repo.User, error) {
			return repo.User{}, pgx.ErrNoRows
		},
	}
	svc := &svc{repo: mock, ja: testTokenAuth}

	_, err := svc.GetProfile(context.Background(), 999)
	if !errors.Is(err, pgx.ErrNoRows) {
		t.Fatalf("expected pgx.ErrNoRows, got %v", err)
	}
}

func TestServiceRegisterUser_Success(t *testing.T) {
	mock := &mockQuerier{
		createUserFn: func(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
			return repo.User{
				ID:    1,
				Name:  arg.Name,
				Email: arg.Email,
				Role:  repo.UserRoleUser,
			}, nil
		},
	}
	svc := &svc{repo: mock, ja: testTokenAuth}

	user, err := svc.RegisterUser(context.Background(), registerParams{
		Email:    "test@example.com",
		Name:     "Test User",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Email != "test@example.com" {
		t.Errorf("expected email='test@example.com', got '%s'", user.Email)
	}
}

func TestServiceLogin_InvalidEmail(t *testing.T) {
	mock := &mockQuerier{
		findUserByEmailFn: func(ctx context.Context, email string) (repo.User, error) {
			return repo.User{}, pgx.ErrNoRows
		},
	}
	svc := &svc{repo: mock, ja: testTokenAuth}

	_, err := svc.Login(context.Background(), loginParams{
		Email:    "unknown@example.com",
		Password: "password123",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestServiceLogin_WrongPassword(t *testing.T) {
	// bcryptハッシュ化された "correctpassword" を生成
	hashed, _ := hashPassword("correctpassword")
	mock := &mockQuerier{
		findUserByEmailFn: func(ctx context.Context, email string) (repo.User, error) {
			return repo.User{
				ID:           1,
				Email:        email,
				PasswordHash: hashed,
				Role:         repo.UserRoleUser,
			}, nil
		},
	}
	svc := &svc{repo: mock, ja: testTokenAuth}

	_, err := svc.Login(context.Background(), loginParams{
		Email:    "test@example.com",
		Password: "wrongpassword",
	})
	if !errors.Is(err, ErrInvalidCredentials) {
		t.Fatalf("expected ErrInvalidCredentials, got %v", err)
	}
}

func TestServiceRefresh_EmptyToken(t *testing.T) {
	svc := &svc{repo: &mockQuerier{}, ja: testTokenAuth}

	_, err := svc.Refresh(context.Background(), "")
	if !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
}

func TestServiceRefresh_InvalidToken(t *testing.T) {
	mock := &mockQuerier{
		consumeRefreshTokenFn: func(ctx context.Context, tokenHash string) (repo.RefreshToken, error) {
			return repo.RefreshToken{}, pgx.ErrNoRows
		},
	}
	svc := &svc{repo: mock, ja: testTokenAuth}

	_, err := svc.Refresh(context.Background(), "invalid-refresh-token")
	if !errors.Is(err, ErrInvalidRefreshToken) {
		t.Fatalf("expected ErrInvalidRefreshToken, got %v", err)
	}
}

func TestServiceUpdateUser_Success(t *testing.T) {
	name := "Updated Name"
	mock := &mockQuerier{
		findUserByIdFn: func(ctx context.Context, id int64) (repo.User, error) {
			return newTestUser(id, "Old Name", "test@example.com"), nil
		},
		updateUserFn: func(ctx context.Context, arg repo.UpdateUserParams) (repo.User, error) {
			return repo.User{
				ID:    arg.ID,
				Name:  arg.Name,
				Email: arg.Email,
				Role:  repo.UserRoleUser,
			}, nil
		},
	}
	svc := &svc{repo: mock, ja: testTokenAuth}

	user, err := svc.UpdateUser(context.Background(), 10, updateUserParams{Name: &name})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if user.Name != "Updated Name" {
		t.Errorf("expected name='Updated Name', got '%s'", user.Name)
	}
}

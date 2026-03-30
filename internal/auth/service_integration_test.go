package auth

import (
	"context"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"example.com/ecommerce/internal/testhelper"
	"github.com/go-chi/jwtauth/v5"
)

var integrationTokenAuth = jwtauth.New("HS256", []byte("integration-test-secret"), nil)

func TestIntegrationLogin_Success(t *testing.T) {
	pool := testhelper.NewTestDB(t)
	ctx := context.Background()
	queries := repo.New(pool)
	svc := NewService(pool, queries, integrationTokenAuth)

	_, err := svc.RegisterUser(ctx, registerParams{
		Email:    "integration@example.com",
		Name:     "Integration User",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	tokens, err := svc.Login(ctx, loginParams{
		Email:    "integration@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}
	if tokens.AccessToken == "" {
		t.Error("expected non-empty access token")
	}
	if tokens.RefreshToken == "" {
		t.Error("expected non-empty refresh token")
	}
}

func TestIntegrationLogin_WrongPassword(t *testing.T) {
	pool := testhelper.NewTestDB(t)
	ctx := context.Background()
	queries := repo.New(pool)
	svc := NewService(pool, queries, integrationTokenAuth)

	_, err := svc.RegisterUser(ctx, registerParams{
		Email:    "wrong@example.com",
		Name:     "User",
		Password: "correctpassword",
	})
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	_, err = svc.Login(ctx, loginParams{
		Email:    "wrong@example.com",
		Password: "wrongpassword",
	})
	if err == nil {
		t.Fatal("expected error for wrong password, got nil")
	}
}

func TestIntegrationRefresh_Success(t *testing.T) {
	pool := testhelper.NewTestDB(t)
	ctx := context.Background()
	queries := repo.New(pool)
	svc := NewService(pool, queries, integrationTokenAuth)

	_, err := svc.RegisterUser(ctx, registerParams{
		Email:    "refresh@example.com",
		Name:     "Refresh User",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	tokens, err := svc.Login(ctx, loginParams{
		Email:    "refresh@example.com",
		Password: "password123",
	})
	if err != nil {
		t.Fatalf("Login failed: %v", err)
	}

	newTokens, err := svc.Refresh(ctx, tokens.RefreshToken)
	if err != nil {
		t.Fatalf("Refresh failed: %v", err)
	}
	if newTokens.AccessToken == "" {
		t.Error("expected non-empty new access token")
	}
}

func TestIntegrationUpdateUserPassword_Success(t *testing.T) {
	pool := testhelper.NewTestDB(t)
	ctx := context.Background()
	queries := repo.New(pool)
	svc := NewService(pool, queries, integrationTokenAuth)

	user, err := svc.RegisterUser(ctx, registerParams{
		Email:    "pwchange@example.com",
		Name:     "PW User",
		Password: "oldpassword",
	})
	if err != nil {
		t.Fatalf("RegisterUser failed: %v", err)
	}

	_, err = svc.UpdateUserPassword(ctx, user.ID, "oldpassword", "newpassword123")
	if err != nil {
		t.Fatalf("UpdateUserPassword failed: %v", err)
	}

	// 新パスワードでログインできること
	_, err = svc.Login(ctx, loginParams{
		Email:    "pwchange@example.com",
		Password: "newpassword123",
	})
	if err != nil {
		t.Fatalf("Login with new password failed: %v", err)
	}

	// 旧パスワードではログインできないこと
	_, err = svc.Login(ctx, loginParams{
		Email:    "pwchange@example.com",
		Password: "oldpassword",
	})
	if err == nil {
		t.Fatal("expected error with old password, got nil")
	}
}

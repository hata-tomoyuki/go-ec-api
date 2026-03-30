package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/jwtauth/v5"
	"github.com/jackc/pgx/v5"
)

// newRequestWithJWT は JWT トークンをコンテキストにセットしたリクエストを生成する。
func newRequestWithJWT(method, path, body string, userID string) *http.Request {
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
	} else {
		r = httptest.NewRequest(method, path, nil)
	}

	tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"sub": userID,
	})
	token, _ := jwtauth.VerifyToken(tokenAuth, tokenString)
	ctx := jwtauth.NewContext(r.Context(), token, nil)
	return r.WithContext(ctx)
}

// newRequestWithLogoutJWT は Logout に必要な jti/exp/rtid クレームを含む JWT をセットしたリクエストを生成する。
func newRequestWithLogoutJWT(method, path, body string, userID string) *http.Request {
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
	} else {
		r = httptest.NewRequest(method, path, nil)
	}

	tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"sub":  userID,
		"jti":  "test-jti-123",
		"exp":  float64(time.Now().Add(15 * time.Minute).Unix()),
		"rtid": "1",
	})
	token, _ := jwtauth.VerifyToken(tokenAuth, tokenString)
	ctx := jwtauth.NewContext(r.Context(), token, nil)
	return r.WithContext(ctx)
}

// ---------- mockService ----------

type mockService struct {
	registerUserFn       func(ctx context.Context, params registerParams) (repo.User, error)
	loginFn              func(ctx context.Context, params loginParams) (LoginTokens, error)
	logoutFn             func(ctx context.Context, jti string, expiredAt time.Time, refreshTokenID int64) error
	refreshFn            func(ctx context.Context, refreshTokenPlain string) (LoginTokens, error)
	getProfileFn         func(ctx context.Context, userID int64) (repo.User, error)
	updateUserFn         func(ctx context.Context, userID int64, params updateUserParams) (repo.User, error)
	updateUserPasswordFn func(ctx context.Context, userID int64, currentPassword, newPassword string) (repo.User, error)
}

func (m *mockService) RegisterUser(ctx context.Context, params registerParams) (repo.User, error) {
	return m.registerUserFn(ctx, params)
}
func (m *mockService) Login(ctx context.Context, params loginParams) (LoginTokens, error) {
	return m.loginFn(ctx, params)
}
func (m *mockService) Logout(ctx context.Context, jti string, expiredAt time.Time, refreshTokenID int64) error {
	return m.logoutFn(ctx, jti, expiredAt, refreshTokenID)
}
func (m *mockService) Refresh(ctx context.Context, refreshTokenPlain string) (LoginTokens, error) {
	return m.refreshFn(ctx, refreshTokenPlain)
}
func (m *mockService) GetProfile(ctx context.Context, userID int64) (repo.User, error) {
	return m.getProfileFn(ctx, userID)
}
func (m *mockService) UpdateUser(ctx context.Context, userID int64, params updateUserParams) (repo.User, error) {
	return m.updateUserFn(ctx, userID, params)
}
func (m *mockService) UpdateUserPassword(ctx context.Context, userID int64, currentPassword, newPassword string) (repo.User, error) {
	return m.updateUserPasswordFn(ctx, userID, currentPassword, newPassword)
}

// ---------- helpers ----------

func newTestUser(id int64, name, email string) repo.User {
	return repo.User{
		ID:    id,
		Name:  name,
		Email: email,
		Role:  repo.UserRoleUser,
	}
}

// ---------- Tests ----------

func TestHandlerRegisterUser_201(t *testing.T) {
	svc := &mockService{
		registerUserFn: func(ctx context.Context, params registerParams) (repo.User, error) {
			return newTestUser(1, params.Name, params.Email), nil
		},
	}
	h := NewHandler(svc)

	body := `{"email":"test@example.com","name":"Test User","password":"password123"}`
	r := httptest.NewRequest("POST", "/register", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RegisterUser(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var resp userResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Email != "test@example.com" {
		t.Errorf("expected email 'test@example.com', got '%s'", resp.Email)
	}
}

func TestHandlerRegisterUser_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	cases := []struct {
		name string
		body string
	}{
		{"empty email", `{"email":"","name":"Test","password":"password123"}`},
		{"empty name", `{"email":"test@example.com","name":"","password":"password123"}`},
		{"short password", `{"email":"test@example.com","name":"Test","password":"short"}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/register", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.RegisterUser(w, r)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestHandlerRegisterUser_400_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockService{})

	r := httptest.NewRequest("POST", "/register", strings.NewReader("{invalid json"))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.RegisterUser(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerLogin_200(t *testing.T) {
	svc := &mockService{
		loginFn: func(ctx context.Context, params loginParams) (LoginTokens, error) {
			return LoginTokens{
				AccessToken:  "access-token",
				RefreshToken: "refresh-token",
				ExpiresIn:    900,
			}, nil
		},
	}
	h := NewHandler(svc)

	body := `{"email":"test@example.com","password":"password123"}`
	r := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var tokens LoginTokens
	if err := json.NewDecoder(w.Body).Decode(&tokens); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if tokens.AccessToken != "access-token" {
		t.Errorf("expected access_token 'access-token', got '%s'", tokens.AccessToken)
	}
}

func TestHandlerLogin_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	cases := []struct {
		name string
		body string
	}{
		{"empty email", `{"email":"","password":"password123"}`},
		{"empty password", `{"email":"test@example.com","password":""}`},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := httptest.NewRequest("POST", "/login", strings.NewReader(tc.body))
			r.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			h.Login(w, r)

			if w.Code != http.StatusBadRequest {
				t.Fatalf("expected status 400, got %d", w.Code)
			}
		})
	}
}

func TestHandlerLogin_401_InvalidCredentials(t *testing.T) {
	svc := &mockService{
		loginFn: func(ctx context.Context, params loginParams) (LoginTokens, error) {
			return LoginTokens{}, ErrInvalidCredentials
		},
	}
	h := NewHandler(svc)

	body := `{"email":"wrong@example.com","password":"wrongpassword"}`
	r := httptest.NewRequest("POST", "/login", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Login(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestHandlerRefresh_200(t *testing.T) {
	svc := &mockService{
		refreshFn: func(ctx context.Context, refreshTokenPlain string) (LoginTokens, error) {
			return LoginTokens{
				AccessToken:  "new-access-token",
				RefreshToken: "new-refresh-token",
				ExpiresIn:    900,
			}, nil
		},
	}
	h := NewHandler(svc)

	body := `{"refresh_token":"valid-refresh-token"}`
	r := httptest.NewRequest("POST", "/refresh", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Refresh(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerRefresh_401(t *testing.T) {
	svc := &mockService{
		refreshFn: func(ctx context.Context, refreshTokenPlain string) (LoginTokens, error) {
			return LoginTokens{}, ErrInvalidRefreshToken
		},
	}
	h := NewHandler(svc)

	body := `{"refresh_token":"invalid-token"}`
	r := httptest.NewRequest("POST", "/refresh", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.Refresh(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestHandlerGetMe_200(t *testing.T) {
	svc := &mockService{
		getProfileFn: func(ctx context.Context, userID int64) (repo.User, error) {
			return newTestUser(10, "Test User", "test@example.com"), nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/me", "", "10")
	w := httptest.NewRecorder()

	h.GetMe(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp userResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.ID != 10 {
		t.Errorf("expected id=10, got %d", resp.ID)
	}
}

func TestHandlerUpdateMe_200(t *testing.T) {
	svc := &mockService{
		updateUserFn: func(ctx context.Context, userID int64, params updateUserParams) (repo.User, error) {
			name := "Updated Name"
			return newTestUser(userID, name, "test@example.com"), nil
		},
	}
	h := NewHandler(svc)

	body := `{"name":"Updated Name"}`
	r := newRequestWithJWT("PUT", "/me", body, "10")
	w := httptest.NewRecorder()

	h.UpdateMe(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var resp userResponse
	if err := json.NewDecoder(w.Body).Decode(&resp); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if resp.Name != "Updated Name" {
		t.Errorf("expected name 'Updated Name', got '%s'", resp.Name)
	}
}

func TestHandlerUpdateMe_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	// name も email も nil（JSON でフィールドを省略）になるような空のオブジェクトでは
	// validate() が ErrUpdateUserValidation を返す
	body := `{}`
	r := newRequestWithJWT("PUT", "/me", body, "10")
	w := httptest.NewRecorder()

	h.UpdateMe(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerUpdatePassword_204(t *testing.T) {
	svc := &mockService{
		updateUserPasswordFn: func(ctx context.Context, userID int64, currentPassword, newPassword string) (repo.User, error) {
			return newTestUser(userID, "Test User", "test@example.com"), nil
		},
	}
	h := NewHandler(svc)

	body := `{"current_password":"oldpassword","new_password":"newpassword123"}`
	r := newRequestWithJWT("PUT", "/me/password", body, "10")
	w := httptest.NewRecorder()

	h.UpdatePassword(w, r)

	// handlers.go の UpdatePassword は成功時に 200 で userResponse を返す
	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerUpdatePassword_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	// new_password が 8 文字未満
	body := `{"current_password":"oldpassword","new_password":"short"}`
	r := newRequestWithJWT("PUT", "/me/password", body, "10")
	w := httptest.NewRecorder()

	h.UpdatePassword(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerUpdatePassword_401_WrongPassword(t *testing.T) {
	svc := &mockService{
		updateUserPasswordFn: func(ctx context.Context, userID int64, currentPassword, newPassword string) (repo.User, error) {
			return repo.User{}, ErrInvalidCredentials
		},
	}
	h := NewHandler(svc)

	body := `{"current_password":"wrongpassword","new_password":"newpassword123"}`
	r := newRequestWithJWT("PUT", "/me/password", body, "10")
	w := httptest.NewRecorder()

	h.UpdatePassword(w, r)

	if w.Code != http.StatusUnauthorized {
		t.Fatalf("expected status 401, got %d", w.Code)
	}
}

func TestHandlerLogout_204(t *testing.T) {
	svc := &mockService{
		logoutFn: func(ctx context.Context, jti string, expiredAt time.Time, refreshTokenID int64) error {
			return nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithLogoutJWT("POST", "/auth/logout", "", "10")
	w := httptest.NewRecorder()

	h.Logout(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d\nbody: %s", w.Code, w.Body.String())
	}
}

func TestHandlerGetMe_404(t *testing.T) {
	svc := &mockService{
		getProfileFn: func(ctx context.Context, userID int64) (repo.User, error) {
			return repo.User{}, pgx.ErrNoRows
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/me", "", "10")
	w := httptest.NewRecorder()

	h.GetMe(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

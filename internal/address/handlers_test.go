package address

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

// テスト用の JWT トークンを生成し、リクエストのコンテキストに設定する
func newRequestWithJWT(method, path, body string, userID string) *http.Request {
	var r *http.Request
	if body != "" {
		r = httptest.NewRequest(method, path, strings.NewReader(body))
		r.Header.Set("Content-Type", "application/json")
	} else {
		r = httptest.NewRequest(method, path, nil)
	}

	// JWT コンテキストを構築
	tokenAuth := jwtauth.New("HS256", []byte("test-secret"), nil)
	_, tokenString, _ := tokenAuth.Encode(map[string]interface{}{
		"sub": userID,
	})
	token, _ := jwtauth.VerifyToken(tokenAuth, tokenString)
	ctx := jwtauth.NewContext(r.Context(), token, nil)
	return r.WithContext(ctx)
}

// chi URLParam をリクエストコンテキストに追加する
func withChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

// ---------- mockService ----------

// handler テストでは Service interface をモックする（repo ではなく）
type mockService struct {
	listAddressesByUserIdFn func(ctx context.Context, userId int64) ([]repo.Address, error)
	findAddressByIdFn       func(ctx context.Context, userId int64, addressId int32) (repo.Address, error)
	createAddressFn         func(ctx context.Context, userId int64, params createAddressParams) (repo.Address, error)
	updateAddressFn         func(ctx context.Context, userId int64, addressId int32, params createAddressParams) (repo.Address, error)
	deleteAddressFn         func(ctx context.Context, userId int64, addressId int32) error
}

func (m *mockService) ListAddressesByUserId(ctx context.Context, userId int64) ([]repo.Address, error) {
	return m.listAddressesByUserIdFn(ctx, userId)
}
func (m *mockService) FindAddressById(ctx context.Context, userId int64, addressId int32) (repo.Address, error) {
	return m.findAddressByIdFn(ctx, userId, addressId)
}
func (m *mockService) CreateAddress(ctx context.Context, userId int64, params createAddressParams) (repo.Address, error) {
	return m.createAddressFn(ctx, userId, params)
}
func (m *mockService) UpdateAddress(ctx context.Context, userId int64, addressId int32, params createAddressParams) (repo.Address, error) {
	return m.updateAddressFn(ctx, userId, addressId, params)
}
func (m *mockService) DeleteAddress(ctx context.Context, userId int64, addressId int32) error {
	return m.deleteAddressFn(ctx, userId, addressId)
}

// ---------- Tests ----------

func TestHandlerListAddresses_200(t *testing.T) {
	svc := &mockService{
		listAddressesByUserIdFn: func(ctx context.Context, userId int64) ([]repo.Address, error) {
			return []repo.Address{newTestAddress(1, userId), newTestAddress(2, userId)}, nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/addresses", "", "10")
	w := httptest.NewRecorder()

	h.ListAddresses(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var addresses []repo.Address
	if err := json.NewDecoder(w.Body).Decode(&addresses); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(addresses) != 2 {
		t.Errorf("expected 2 addresses, got %d", len(addresses))
	}
}

func TestHandlerFindAddressById_200(t *testing.T) {
	svc := &mockService{
		findAddressByIdFn: func(ctx context.Context, userId int64, addressId int32) (repo.Address, error) {
			return newTestAddress(1, userId), nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/addresses/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.FindAddressById(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerFindAddressById_404(t *testing.T) {
	svc := &mockService{
		findAddressByIdFn: func(ctx context.Context, userId int64, addressId int32) (repo.Address, error) {
			return repo.Address{}, ErrAddressNotFound
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/addresses/999", "", "10")
	r = withChiURLParam(r, "id", "999")
	w := httptest.NewRecorder()

	h.FindAddressById(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandlerFindAddressById_403(t *testing.T) {
	svc := &mockService{
		findAddressByIdFn: func(ctx context.Context, userId int64, addressId int32) (repo.Address, error) {
			return repo.Address{}, ErrForbidden
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/addresses/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.FindAddressById(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestHandlerFindAddressById_400_InvalidID(t *testing.T) {
	h := NewHandler(&mockService{})

	r := newRequestWithJWT("GET", "/addresses/abc", "", "10")
	r = withChiURLParam(r, "id", "abc")
	w := httptest.NewRecorder()

	h.FindAddressById(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerCreateAddress_201(t *testing.T) {
	svc := &mockService{
		createAddressFn: func(ctx context.Context, userId int64, params createAddressParams) (repo.Address, error) {
			return newTestAddress(1, userId), nil
		},
	}
	h := NewHandler(svc)

	body := `{"street":"1-1 Shibuya","city":"Shibuya","state":"Tokyo","zip_code":"150-0001","country":"Japan"}`
	r := newRequestWithJWT("POST", "/addresses", body, "10")
	w := httptest.NewRecorder()

	h.CreateAddress(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestHandlerCreateAddress_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	// street が空
	body := `{"street":"","city":"Shibuya","state":"Tokyo","zip_code":"150-0001","country":"Japan"}`
	r := newRequestWithJWT("POST", "/addresses", body, "10")
	w := httptest.NewRecorder()

	h.CreateAddress(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerCreateAddress_400_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockService{})

	r := newRequestWithJWT("POST", "/addresses", "{invalid json", "10")
	w := httptest.NewRecorder()

	h.CreateAddress(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerUpdateAddress_200(t *testing.T) {
	svc := &mockService{
		updateAddressFn: func(ctx context.Context, userId int64, addressId int32, params createAddressParams) (repo.Address, error) {
			return newTestAddress(1, userId), nil
		},
	}
	h := NewHandler(svc)

	body := `{"street":"2-2 Minato","city":"Minato","state":"Tokyo","zip_code":"105-0001","country":"Japan"}`
	r := newRequestWithJWT("PUT", "/addresses/1", body, "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateAddress(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerUpdateAddress_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	body := `{"street":"","city":"","state":"","zip_code":"","country":""}`
	r := newRequestWithJWT("PUT", "/addresses/1", body, "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateAddress(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerDeleteAddress_204(t *testing.T) {
	svc := &mockService{
		deleteAddressFn: func(ctx context.Context, userId int64, addressId int32) error {
			return nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/addresses/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.DeleteAddress(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}

func TestHandlerDeleteAddress_404(t *testing.T) {
	svc := &mockService{
		deleteAddressFn: func(ctx context.Context, userId int64, addressId int32) error {
			return ErrAddressNotFound
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/addresses/999", "", "10")
	r = withChiURLParam(r, "id", "999")
	w := httptest.NewRecorder()

	h.DeleteAddress(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

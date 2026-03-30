package carts

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
	"github.com/jackc/pgx/v5/pgtype"
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

// withChiURLParam は chi URLParam をリクエストコンテキストに追加する
func withChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

// ---------- mockService ----------

type mockService struct {
	createCartFn               func(ctx context.Context, userID int64) (repo.Cart, error)
	addItemToCartFn            func(ctx context.Context, userID int64, productID int64, quantity int) (repo.CartItem, error)
	listCartItemsByUserIdFn    func(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error)
	updateCartItemQuantityFn   func(ctx context.Context, userID int64, cartItemID int64, quantity int) (repo.CartItem, error)
	removeItemFromCartFn       func(ctx context.Context, userID int64, cartItemID int64) (repo.CartItem, error)
	clearCartFn                func(ctx context.Context, userID int64) error
}

func (m *mockService) CreateCart(ctx context.Context, userID int64) (repo.Cart, error) {
	return m.createCartFn(ctx, userID)
}
func (m *mockService) AddItemToCart(ctx context.Context, userID int64, productID int64, quantity int) (repo.CartItem, error) {
	return m.addItemToCartFn(ctx, userID, productID, quantity)
}
func (m *mockService) ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
	return m.listCartItemsByUserIdFn(ctx, userID)
}
func (m *mockService) UpdateCartItemQuantity(ctx context.Context, userID int64, cartItemID int64, quantity int) (repo.CartItem, error) {
	return m.updateCartItemQuantityFn(ctx, userID, cartItemID, quantity)
}
func (m *mockService) RemoveItemFromCart(ctx context.Context, userID int64, cartItemID int64) (repo.CartItem, error) {
	return m.removeItemFromCartFn(ctx, userID, cartItemID)
}
func (m *mockService) ClearCart(ctx context.Context, userID int64) error {
	return m.clearCartFn(ctx, userID)
}

// ---------- helpers ----------

func newTestCart(id int64, userID int64) repo.Cart {
	return repo.Cart{
		ID:        id,
		UserID:    userID,
		CreatedAt: pgtype.Timestamptz{Valid: true},
		UpdatedAt: pgtype.Timestamptz{Valid: true},
	}
}

func newTestCartItem(id int64, cartID int64, productID int64) repo.CartItem {
	return repo.CartItem{
		ID:        id,
		CartID:    cartID,
		ProductID: productID,
		Quantity:  1,
		CreatedAt: pgtype.Timestamptz{Valid: true},
		UpdatedAt: pgtype.Timestamptz{Valid: true},
	}
}

// ---------- Tests ----------

func TestHandlerCreateCart_201(t *testing.T) {
	svc := &mockService{
		createCartFn: func(ctx context.Context, userID int64) (repo.Cart, error) {
			return newTestCart(1, userID), nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("POST", "/cart", "", "10")
	w := httptest.NewRecorder()

	h.CreateCart(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var cart repo.Cart
	if err := json.NewDecoder(w.Body).Decode(&cart); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if cart.UserID != 10 {
		t.Errorf("expected user_id=10, got %d", cart.UserID)
	}
}

func TestHandlerAddItemToCart_201(t *testing.T) {
	svc := &mockService{
		addItemToCartFn: func(ctx context.Context, userID int64, productID int64, quantity int) (repo.CartItem, error) {
			return newTestCartItem(1, 1, productID), nil
		},
	}
	h := NewHandler(svc)

	body := `{"product_id":1,"quantity":2}`
	r := newRequestWithJWT("POST", "/cart/items", body, "10")
	w := httptest.NewRecorder()

	h.AddItemToCart(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}
}

func TestHandlerAddItemToCart_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	// quantity が 0
	body := `{"product_id":1,"quantity":0}`
	r := newRequestWithJWT("POST", "/cart/items", body, "10")
	w := httptest.NewRecorder()

	h.AddItemToCart(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerAddItemToCart_404_CartNotFound(t *testing.T) {
	svc := &mockService{
		addItemToCartFn: func(ctx context.Context, userID int64, productID int64, quantity int) (repo.CartItem, error) {
			return repo.CartItem{}, ErrCartNotFound
		},
	}
	h := NewHandler(svc)

	body := `{"product_id":1,"quantity":2}`
	r := newRequestWithJWT("POST", "/cart/items", body, "10")
	w := httptest.NewRecorder()

	h.AddItemToCart(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandlerShowCartItems_200(t *testing.T) {
	svc := &mockService{
		listCartItemsByUserIdFn: func(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
			return []repo.ListCartItemsByUserIdRow{
				{ID: 1, ProductID: 1, Quantity: 2, ProductPriceInCents: 1000},
				{ID: 2, ProductID: 2, Quantity: 1, ProductPriceInCents: 2000},
			}, nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/cart/items", "", "10")
	w := httptest.NewRecorder()

	h.ShowCartItems(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var items []repo.ListCartItemsByUserIdRow
	if err := json.NewDecoder(w.Body).Decode(&items); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestHandlerUpdateCartItemQuantity_200(t *testing.T) {
	svc := &mockService{
		updateCartItemQuantityFn: func(ctx context.Context, userID int64, cartItemID int64, quantity int) (repo.CartItem, error) {
			return newTestCartItem(cartItemID, 1, 1), nil
		},
	}
	h := NewHandler(svc)

	body := `{"cart_item_id":1,"quantity":3}`
	r := newRequestWithJWT("PUT", "/cart/items", body, "10")
	w := httptest.NewRecorder()

	h.UpdateCartItemQuantity(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerUpdateCartItemQuantity_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	// cart_item_id が 0
	body := `{"cart_item_id":0,"quantity":3}`
	r := newRequestWithJWT("PUT", "/cart/items", body, "10")
	w := httptest.NewRecorder()

	h.UpdateCartItemQuantity(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerUpdateCartItemQuantity_403(t *testing.T) {
	svc := &mockService{
		updateCartItemQuantityFn: func(ctx context.Context, userID int64, cartItemID int64, quantity int) (repo.CartItem, error) {
			return repo.CartItem{}, ErrCartForbidden
		},
	}
	h := NewHandler(svc)

	body := `{"cart_item_id":1,"quantity":3}`
	r := newRequestWithJWT("PUT", "/cart/items", body, "10")
	w := httptest.NewRecorder()

	h.UpdateCartItemQuantity(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestHandlerRemoveItemFromCart_204(t *testing.T) {
	svc := &mockService{
		removeItemFromCartFn: func(ctx context.Context, userID int64, cartItemID int64) (repo.CartItem, error) {
			return newTestCartItem(cartItemID, 1, 1), nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/cart/items/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.RemoveItemFromCart(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}

func TestHandlerRemoveItemFromCart_403(t *testing.T) {
	svc := &mockService{
		removeItemFromCartFn: func(ctx context.Context, userID int64, cartItemID int64) (repo.CartItem, error) {
			return repo.CartItem{}, ErrCartForbidden
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/cart/items/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.RemoveItemFromCart(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestHandlerClearCart_204(t *testing.T) {
	svc := &mockService{
		clearCartFn: func(ctx context.Context, userID int64) error {
			return nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/cart", "", "10")
	w := httptest.NewRecorder()

	h.ClearCart(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}

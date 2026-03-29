package products

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/chi/v5"
)

// chi URLParam をリクエストコンテキストに追加する
func withChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

// ---------- mockService ----------

type mockService struct {
	listProductsFn    func(ctx context.Context) ([]repo.Product, error)
	findProductByIdFn func(ctx context.Context, id int64) (repo.Product, error)
	createProductFn   func(ctx context.Context, p createProductParams) (repo.Product, error)
	updateProductFn   func(ctx context.Context, p updateProductParams) (repo.Product, error)
	deleteProductFn   func(ctx context.Context, id int64) error
}

func (m *mockService) ListProducts(ctx context.Context) ([]repo.Product, error) {
	return m.listProductsFn(ctx)
}
func (m *mockService) FindProductById(ctx context.Context, id int64) (repo.Product, error) {
	return m.findProductByIdFn(ctx, id)
}
func (m *mockService) CreateProduct(ctx context.Context, p createProductParams) (repo.Product, error) {
	return m.createProductFn(ctx, p)
}
func (m *mockService) UpdateProduct(ctx context.Context, p updateProductParams) (repo.Product, error) {
	return m.updateProductFn(ctx, p)
}
func (m *mockService) DeleteProduct(ctx context.Context, id int64) error {
	return m.deleteProductFn(ctx, id)
}

// ---------- Tests ----------

func TestHandlerListProduct_200(t *testing.T) {
	svc := &mockService{
		listProductsFn: func(ctx context.Context) ([]repo.Product, error) {
			return []repo.Product{
				newTestProduct(1, "T-shirt", 2000),
				newTestProduct(2, "Hoodie", 5000),
			}, nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	h.ListProduct(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var products []repo.Product
	if err := json.NewDecoder(w.Body).Decode(&products); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestHandlerListProduct_500(t *testing.T) {
	svc := &mockService{
		listProductsFn: func(ctx context.Context) ([]repo.Product, error) {
			return nil, errors.New("db error")
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("GET", "/products", nil)
	w := httptest.NewRecorder()

	h.ListProduct(w, r)

	if w.Code != http.StatusInternalServerError {
		t.Fatalf("expected status 500, got %d", w.Code)
	}
}

func TestHandlerFindProductById_200(t *testing.T) {
	svc := &mockService{
		findProductByIdFn: func(ctx context.Context, id int64) (repo.Product, error) {
			return newTestProduct(1, "T-shirt", 2000), nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("GET", "/products/1", nil)
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.FindProductById(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerFindProductById_400_InvalidID(t *testing.T) {
	h := NewHandler(&mockService{})

	r := httptest.NewRequest("GET", "/products/abc", nil)
	r = withChiURLParam(r, "id", "abc")
	w := httptest.NewRecorder()

	h.FindProductById(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerCreateProduct_201(t *testing.T) {
	svc := &mockService{
		createProductFn: func(ctx context.Context, p createProductParams) (repo.Product, error) {
			return newTestProduct(1, p.Name, p.PriceInCents), nil
		},
	}
	h := NewHandler(svc)

	body := `{"name":"New Hat","price_in_cents":1500}`
	r := httptest.NewRequest("POST", "/products", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateProduct(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var product repo.Product
	if err := json.NewDecoder(w.Body).Decode(&product); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if product.Name != "New Hat" {
		t.Errorf("expected name 'New Hat', got '%s'", product.Name)
	}
}

func TestHandlerCreateProduct_400_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockService{})

	r := httptest.NewRequest("POST", "/products", strings.NewReader("{bad json"))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateProduct(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerDeleteProduct_204(t *testing.T) {
	svc := &mockService{
		deleteProductFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("DELETE", "/products/1", nil)
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.DeleteProduct(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}

func TestHandlerDeleteProduct_400_InvalidID(t *testing.T) {
	h := NewHandler(&mockService{})

	r := httptest.NewRequest("DELETE", "/products/abc", nil)
	r = withChiURLParam(r, "id", "abc")
	w := httptest.NewRecorder()

	h.DeleteProduct(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

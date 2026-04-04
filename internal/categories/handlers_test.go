package categories

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/go-chi/chi/v5"
)

// withChiURLParam は chi URLParam をリクエストコンテキストに追加する
func withChiURLParam(r *http.Request, key, value string) *http.Request {
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add(key, value)
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	return r.WithContext(ctx)
}

// ---------- mockService ----------

type mockService struct {
	listCategoriesFn            func(ctx context.Context) ([]repo.ListCategoriesRow, error)
	findCategoryByIdFn          func(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error)
	createCategoriesFn          func(ctx context.Context, name string, description *string, imageColor string) (repo.Category, error)
	updateCategoriesFn          func(ctx context.Context, id int64, name string, description *string, imageColor string) (repo.Category, error)
	deleteCategoryFn            func(ctx context.Context, id int64) error
	listProductsByCategoryFn    func(ctx context.Context, categoryId int64) ([]repo.ListProductsByCategoryRow, error)
	addProductToCategoryFn      func(ctx context.Context, categoryId int64, productId int64) error
	removeProductFromCategoryFn func(ctx context.Context, categoryId int64, productId int64) error
}

func (m *mockService) ListCategories(ctx context.Context) ([]repo.ListCategoriesRow, error) {
	return m.listCategoriesFn(ctx)
}
func (m *mockService) FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
	return m.findCategoryByIdFn(ctx, id)
}
func (m *mockService) CreateCategories(ctx context.Context, name string, description *string, imageColor string) (repo.Category, error) {
	return m.createCategoriesFn(ctx, name, description, imageColor)
}
func (m *mockService) UpdateCategories(ctx context.Context, id int64, name string, description *string, imageColor string) (repo.Category, error) {
	return m.updateCategoriesFn(ctx, id, name, description, imageColor)
}
func (m *mockService) DeleteCategory(ctx context.Context, id int64) error {
	return m.deleteCategoryFn(ctx, id)
}
func (m *mockService) ListProductsByCategory(ctx context.Context, categoryId int64) ([]repo.ListProductsByCategoryRow, error) {
	return m.listProductsByCategoryFn(ctx, categoryId)
}
func (m *mockService) AddProductToCategory(ctx context.Context, categoryId int64, productId int64) error {
	return m.addProductToCategoryFn(ctx, categoryId, productId)
}
func (m *mockService) RemoveProductFromCategory(ctx context.Context, categoryId int64, productId int64) error {
	return m.removeProductFromCategoryFn(ctx, categoryId, productId)
}

// ---------- Tests ----------

func TestHandlerCreateCategories_201(t *testing.T) {
	svc := &mockService{
		createCategoriesFn: func(ctx context.Context, name string, description *string, imageColor string) (repo.Category, error) {
			return newTestCategory(1, name), nil
		},
	}
	h := NewHandler(svc)

	body := `{"name":"Electronics"}`
	r := httptest.NewRequest("POST", "/categories", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCategories(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var category repo.Category
	if err := json.NewDecoder(w.Body).Decode(&category); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if category.Name != "Electronics" {
		t.Errorf("expected name 'Electronics', got '%s'", category.Name)
	}
}

func TestHandlerCreateCategories_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	// name が空
	body := `{"name":""}`
	r := httptest.NewRequest("POST", "/categories", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	h.CreateCategories(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerListCategories_200(t *testing.T) {
	svc := &mockService{
		listCategoriesFn: func(ctx context.Context) ([]repo.ListCategoriesRow, error) {
			return []repo.ListCategoriesRow{
				newTestListCategoriesRow(1, "Electronics"),
				newTestListCategoriesRow(2, "Clothing"),
			}, nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("GET", "/categories", nil)
	w := httptest.NewRecorder()

	h.ListCategories(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var categories []repo.ListCategoriesRow
	if err := json.NewDecoder(w.Body).Decode(&categories); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(categories) != 2 {
		t.Errorf("expected 2 categories, got %d", len(categories))
	}
}

func TestHandlerFindCategoryById_200(t *testing.T) {
	svc := &mockService{
		findCategoryByIdFn: func(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
			return newTestFindCategoryByIdRow(1, "Electronics"), nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("GET", "/categories/1", nil)
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.FindCategoryById(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerFindCategoryById_404(t *testing.T) {
	svc := &mockService{
		findCategoryByIdFn: func(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
			return repo.FindCategoryByIdRow{}, ErrCategoryNotFound
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("GET", "/categories/999", nil)
	r = withChiURLParam(r, "id", "999")
	w := httptest.NewRecorder()

	h.FindCategoryById(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandlerFindCategoryById_400_InvalidID(t *testing.T) {
	h := NewHandler(&mockService{})

	r := httptest.NewRequest("GET", "/categories/abc", nil)
	r = withChiURLParam(r, "id", "abc")
	w := httptest.NewRecorder()

	h.FindCategoryById(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerUpdateCategory_200(t *testing.T) {
	svc := &mockService{
		updateCategoriesFn: func(ctx context.Context, id int64, name string, description *string, imageColor string) (repo.Category, error) {
			return newTestCategory(id, name), nil
		},
	}
	h := NewHandler(svc)

	body := `{"name":"Updated Electronics"}`
	r := httptest.NewRequest("PUT", "/categories/1", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateCategory(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerUpdateCategory_400_Validation(t *testing.T) {
	h := NewHandler(&mockService{})

	body := `{"name":""}`
	r := httptest.NewRequest("PUT", "/categories/1", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateCategory(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerDeleteCategory_204(t *testing.T) {
	svc := &mockService{
		deleteCategoryFn: func(ctx context.Context, id int64) error {
			return nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("DELETE", "/categories/1", nil)
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.DeleteCategory(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}

func TestHandlerDeleteCategory_404(t *testing.T) {
	svc := &mockService{
		deleteCategoryFn: func(ctx context.Context, id int64) error {
			return ErrCategoryNotFound
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("DELETE", "/categories/999", nil)
	r = withChiURLParam(r, "id", "999")
	w := httptest.NewRecorder()

	h.DeleteCategory(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandlerListProductsByCategory_200(t *testing.T) {
	svc := &mockService{
		listProductsByCategoryFn: func(ctx context.Context, categoryId int64) ([]repo.ListProductsByCategoryRow, error) {
			return []repo.ListProductsByCategoryRow{
				{ID: 1, Name: "T-shirt", PriceInCents: 2000},
				{ID: 2, Name: "Hoodie", PriceInCents: 5000},
			}, nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("GET", "/categories/1/products", nil)
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.ListProductsByCategory(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var products []repo.ListProductsByCategoryRow
	if err := json.NewDecoder(w.Body).Decode(&products); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(products) != 2 {
		t.Errorf("expected 2 products, got %d", len(products))
	}
}

func TestHandlerAddProductToCategory_200(t *testing.T) {
	svc := &mockService{
		addProductToCategoryFn: func(ctx context.Context, categoryId int64, productId int64) error {
			return nil
		},
	}
	h := NewHandler(svc)

	body := `{"product_id":1}`
	r := httptest.NewRequest("POST", "/categories/1/products", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.AddProductToCategory(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerRemoveProductFromCategory_204(t *testing.T) {
	svc := &mockService{
		removeProductFromCategoryFn: func(ctx context.Context, categoryId int64, productId int64) error {
			return nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("DELETE", "/categories/1/products/2", nil)
	rctx := chi.NewRouteContext()
	rctx.URLParams.Add("id", "1")
	rctx.URLParams.Add("productId", "2")
	ctx := context.WithValue(r.Context(), chi.RouteCtxKey, rctx)
	r = r.WithContext(ctx)
	w := httptest.NewRecorder()

	h.RemoveProductFromCategory(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}

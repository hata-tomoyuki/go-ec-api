package orders

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
	listAllOrdersFn            func(ctx context.Context) ([]repo.ListAllOrdersRow, error)
	listAllOrdersPaginatedFn   func(ctx context.Context, params listOrdersParams) (paginatedOrders, error)
	listOrdersByCustomerIDFn   func(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error)
	listOrdersPaginatedFn      func(ctx context.Context, customerID int64, params listOrdersParams) (paginatedOrders, error)
	findOrderByIdFn            func(ctx context.Context, orderID int64) ([]repo.FindOrderByIdRow, error)
	placeOrderFn               func(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
	cancelOrderFn              func(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error)
	updateOrderStatusFn        func(ctx context.Context, orderID int64, status string) (repo.FindOrderByIdRow, error)
}

func (m *mockService) ListAllOrders(ctx context.Context) ([]repo.ListAllOrdersRow, error) {
	return m.listAllOrdersFn(ctx)
}
func (m *mockService) ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error) {
	return m.listOrdersByCustomerIDFn(ctx, customerID)
}
func (m *mockService) FindOrderById(ctx context.Context, orderID int64) ([]repo.FindOrderByIdRow, error) {
	return m.findOrderByIdFn(ctx, orderID)
}
func (m *mockService) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	return m.placeOrderFn(ctx, tempOrder)
}
func (m *mockService) CancelOrder(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error) {
	return m.cancelOrderFn(ctx, orderID, customerID)
}
func (m *mockService) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (repo.FindOrderByIdRow, error) {
	return m.updateOrderStatusFn(ctx, orderID, status)
}
func (m *mockService) ListOrdersPaginated(ctx context.Context, customerID int64, params listOrdersParams) (paginatedOrders, error) {
	return m.listOrdersPaginatedFn(ctx, customerID, params)
}
func (m *mockService) ListAllOrdersPaginated(ctx context.Context, params listOrdersParams) (paginatedOrders, error) {
	return m.listAllOrdersPaginatedFn(ctx, params)
}

// ---------- helpers ----------

func newTestOrderRow(id int64, customerID int64, status repo.Status) []repo.FindOrderByIdRow {
	return []repo.FindOrderByIdRow{
		{
			ID:           id,
			CustomerID:   customerID,
			Status:       status,
			CreatedAt:    pgtype.Timestamptz{Valid: true},
			UpdatedAt:    pgtype.Timestamptz{Valid: true},
			ProductID:    1,
			Quantity:     2,
			PriceInCents: 1000,
		},
	}
}

func newTestOrder(id int64, customerID int64) repo.Order {
	return repo.Order{
		ID:         id,
		CustomerID: customerID,
		Status:     repo.StatusPending,
		CreatedAt:  pgtype.Timestamptz{Valid: true},
		UpdatedAt:  pgtype.Timestamptz{Valid: true},
	}
}

// ---------- Tests ----------

func TestHandlerListOrdersByCustomerID_200(t *testing.T) {
	svc := &mockService{
		listOrdersPaginatedFn: func(ctx context.Context, customerID int64, params listOrdersParams) (paginatedOrders, error) {
			return paginatedOrders{
				Data: []paginatedOrderRow{
					{ID: 1, CustomerID: customerID, Status: "pending"},
					{ID: 2, CustomerID: customerID, Status: "completed"},
				},
				Total: 2,
				Page:  1,
				Limit: 20,
			}, nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/orders", "", "10")
	w := httptest.NewRecorder()

	h.ListOrdersByCustomerID(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result paginatedOrders
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result.Data) != 2 {
		t.Errorf("expected 2 orders, got %d", len(result.Data))
	}
}

func TestHandlerFindOrderById_200(t *testing.T) {
	svc := &mockService{
		findOrderByIdFn: func(ctx context.Context, orderID int64) ([]repo.FindOrderByIdRow, error) {
			return newTestOrderRow(orderID, 10, repo.StatusPending), nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/orders/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.FindOrderById(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerFindOrderById_404(t *testing.T) {
	svc := &mockService{
		findOrderByIdFn: func(ctx context.Context, orderID int64) ([]repo.FindOrderByIdRow, error) {
			return nil, ErrOrderNotFound
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("GET", "/orders/999", "", "10")
	r = withChiURLParam(r, "id", "999")
	w := httptest.NewRecorder()

	h.FindOrderById(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandlerFindOrderById_403(t *testing.T) {
	svc := &mockService{
		findOrderByIdFn: func(ctx context.Context, orderID int64) ([]repo.FindOrderByIdRow, error) {
			// 別のユーザー(99)の注文を返す
			return newTestOrderRow(orderID, 99, repo.StatusPending), nil
		},
	}
	h := NewHandler(svc)

	// userID=10 でアクセスするが、注文の customerID=99
	r := newRequestWithJWT("GET", "/orders/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.FindOrderById(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestHandlerFindOrderById_400_InvalidID(t *testing.T) {
	h := NewHandler(&mockService{})

	r := newRequestWithJWT("GET", "/orders/abc", "", "10")
	r = withChiURLParam(r, "id", "abc")
	w := httptest.NewRecorder()

	h.FindOrderById(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerPlaceOrder_201(t *testing.T) {
	svc := &mockService{
		placeOrderFn: func(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
			return newTestOrder(1, tempOrder.CustomerID), nil
		},
	}
	h := NewHandler(svc)

	body := `{"items":[{"product_id":1,"quantity":2}]}`
	r := newRequestWithJWT("POST", "/orders", body, "10")
	w := httptest.NewRecorder()

	h.PlaceOrder(w, r)

	if w.Code != http.StatusCreated {
		t.Fatalf("expected status 201, got %d", w.Code)
	}

	var order repo.Order
	if err := json.NewDecoder(w.Body).Decode(&order); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if order.ID != 1 {
		t.Errorf("expected order id=1, got %d", order.ID)
	}
}

func TestHandlerPlaceOrder_400_EmptyItems(t *testing.T) {
	svc := &mockService{
		placeOrderFn: func(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
			return repo.Order{}, ErrOrderEmptyItems
		},
	}
	h := NewHandler(svc)

	body := `{"items":[]}`
	r := newRequestWithJWT("POST", "/orders", body, "10")
	w := httptest.NewRecorder()

	h.PlaceOrder(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerPlaceOrder_400_InvalidItemQuantity(t *testing.T) {
	svc := &mockService{
		placeOrderFn: func(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
			return repo.Order{}, ErrOrderItemValidation
		},
	}
	h := NewHandler(svc)

	body := `{"items":[{"product_id":1,"quantity":0}]}`
	r := newRequestWithJWT("POST", "/orders", body, "10")
	w := httptest.NewRecorder()

	h.PlaceOrder(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerPlaceOrder_400_InvalidJSON(t *testing.T) {
	h := NewHandler(&mockService{})

	r := newRequestWithJWT("POST", "/orders", "{invalid json", "10")
	w := httptest.NewRecorder()

	h.PlaceOrder(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerCancelOrder_204(t *testing.T) {
	svc := &mockService{
		cancelOrderFn: func(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error) {
			return newTestOrderRow(orderID, customerID, repo.StatusPending)[0], nil
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/orders/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.CancelOrder(w, r)

	if w.Code != http.StatusNoContent {
		t.Fatalf("expected status 204, got %d", w.Code)
	}
}

func TestHandlerCancelOrder_404(t *testing.T) {
	svc := &mockService{
		cancelOrderFn: func(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error) {
			return repo.FindOrderByIdRow{}, ErrOrderNotFound
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/orders/999", "", "10")
	r = withChiURLParam(r, "id", "999")
	w := httptest.NewRecorder()

	h.CancelOrder(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandlerCancelOrder_403(t *testing.T) {
	svc := &mockService{
		cancelOrderFn: func(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error) {
			return repo.FindOrderByIdRow{}, ErrOrderForbidden
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/orders/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.CancelOrder(w, r)

	if w.Code != http.StatusForbidden {
		t.Fatalf("expected status 403, got %d", w.Code)
	}
}

func TestHandlerCancelOrder_400_NotPending(t *testing.T) {
	svc := &mockService{
		cancelOrderFn: func(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error) {
			return repo.FindOrderByIdRow{}, ErrOrderNotPending
		},
	}
	h := NewHandler(svc)

	r := newRequestWithJWT("DELETE", "/orders/1", "", "10")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.CancelOrder(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerUpdateOrderStatus_200(t *testing.T) {
	svc := &mockService{
		updateOrderStatusFn: func(ctx context.Context, orderID int64, status string) (repo.FindOrderByIdRow, error) {
			return newTestOrderRow(orderID, 10, repo.Status(status))[0], nil
		},
	}
	h := NewHandler(svc)

	body := `{"status":"completed"}`
	r := httptest.NewRequest("PUT", "/orders/1/status", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateOrderStatus(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}
}

func TestHandlerUpdateOrderStatus_400_InvalidStatus(t *testing.T) {
	svc := &mockService{
		updateOrderStatusFn: func(ctx context.Context, orderID int64, status string) (repo.FindOrderByIdRow, error) {
			return repo.FindOrderByIdRow{}, ErrInvalidStatus
		},
	}
	h := NewHandler(svc)

	body := `{"status":"invalid-status"}`
	r := httptest.NewRequest("PUT", "/orders/1/status", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = withChiURLParam(r, "id", "1")
	w := httptest.NewRecorder()

	h.UpdateOrderStatus(w, r)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("expected status 400, got %d", w.Code)
	}
}

func TestHandlerUpdateOrderStatus_404(t *testing.T) {
	svc := &mockService{
		updateOrderStatusFn: func(ctx context.Context, orderID int64, status string) (repo.FindOrderByIdRow, error) {
			return repo.FindOrderByIdRow{}, ErrOrderNotFound
		},
	}
	h := NewHandler(svc)

	body := `{"status":"completed"}`
	r := httptest.NewRequest("PUT", "/orders/999/status", strings.NewReader(body))
	r.Header.Set("Content-Type", "application/json")
	r = withChiURLParam(r, "id", "999")
	w := httptest.NewRecorder()

	h.UpdateOrderStatus(w, r)

	if w.Code != http.StatusNotFound {
		t.Fatalf("expected status 404, got %d", w.Code)
	}
}

func TestHandlerListAllOrders_200(t *testing.T) {
	svc := &mockService{
		listAllOrdersPaginatedFn: func(ctx context.Context, params listOrdersParams) (paginatedOrders, error) {
			return paginatedOrders{
				Data: []paginatedOrderRow{
					{ID: 1, CustomerID: 10, Status: "pending"},
					{ID: 2, CustomerID: 20, Status: "completed"},
				},
				Total: 2,
				Page:  1,
				Limit: 20,
			}, nil
		},
	}
	h := NewHandler(svc)

	r := httptest.NewRequest("GET", "/admin/orders", nil)
	w := httptest.NewRecorder()

	h.ListAllOrders(w, r)

	if w.Code != http.StatusOK {
		t.Fatalf("expected status 200, got %d", w.Code)
	}

	var result paginatedOrders
	if err := json.NewDecoder(w.Body).Decode(&result); err != nil {
		t.Fatalf("failed to decode response: %v", err)
	}
	if len(result.Data) != 2 {
		t.Errorf("expected 2 orders, got %d", len(result.Data))
	}
}

package orders

import (
	"context"
	"errors"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

// mockDB は pgBeginner interface のモック（トランザクション不要のテスト用）
// 実際のトランザクションを開始しない nil Tx を返すだけ
// PlaceOrder はトランザクションを使うため、このモックではテストしない

func TestServiceListOrdersByCustomerID(t *testing.T) {
	mock := &mockQuerier{
		listOrdersByCustomerIDFn: func(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error) {
			return []repo.ListOrdersByCustomerIDRow{
				{ID: 1, CustomerID: customerID, Status: repo.StatusPending},
				{ID: 2, CustomerID: customerID, Status: repo.StatusCompleted},
			}, nil
		},
	}
	svc := &svc{q: mock}

	orders, err := svc.ListOrdersByCustomerID(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 2 {
		t.Errorf("expected 2 orders, got %d", len(orders))
	}
}

func TestServiceListAllOrders(t *testing.T) {
	mock := &mockQuerier{
		listAllOrdersFn: func(ctx context.Context) ([]repo.ListAllOrdersRow, error) {
			return []repo.ListAllOrdersRow{
				{ID: 1},
				{ID: 2},
				{ID: 3},
			}, nil
		},
	}
	svc := &svc{q: mock}

	orders, err := svc.ListAllOrders(context.Background())
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(orders) != 3 {
		t.Errorf("expected 3 orders, got %d", len(orders))
	}
}

func TestServiceFindOrderById_Success(t *testing.T) {
	mock := &mockQuerier{
		findOrderByIdFn: func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
			return newTestOrderHelper(id, 10, repo.StatusPending), nil
		},
	}
	svc := &svc{q: mock}

	rows, err := svc.FindOrderById(context.Background(), 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if rows[0].ID != 1 {
		t.Errorf("expected id=1, got %d", rows[0].ID)
	}
}

func TestServiceFindOrderById_NotFound(t *testing.T) {
	mock := &mockQuerier{
		findOrderByIdFn: func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
			return nil, nil
		},
	}
	svc := &svc{q: mock}

	_, err := svc.FindOrderById(context.Background(), 999)
	if !errors.Is(err, ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestServiceCancelOrder_Success(t *testing.T) {
	mock := &mockQuerier{
		findOrderByIdFn: func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
			return newTestOrderHelper(id, 10, repo.StatusPending), nil
		},
		cancelOrderFn: func(ctx context.Context, id int64) (repo.Order, error) {
			return repo.Order{ID: id, CustomerID: 10, Status: repo.StatusCancelled}, nil
		},
	}

	// txMock 用の Querier も同じ mock を使う
	tx := &mockTx{}
	db := &mockDB{tx: tx}
	svc := &svc{
		q:  mock,
		db: db,
		newTxQ: func(_ pgx.Tx) repo.Querier {
			return mock
		},
	}

	order, err := svc.CancelOrder(context.Background(), 1, 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.CustomerID != 10 {
		t.Errorf("expected customer_id=10, got %d", order.CustomerID)
	}
	if !tx.committed {
		t.Error("expected transaction to be committed")
	}
}

func TestServiceCancelOrder_NotFound(t *testing.T) {
	mock := &mockQuerier{
		findOrderByIdFn: func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
			return nil, nil
		},
	}
	svc := &svc{q: mock}

	_, err := svc.CancelOrder(context.Background(), 999, 10)
	if !errors.Is(err, ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestServiceCancelOrder_Forbidden(t *testing.T) {
	mock := &mockQuerier{
		findOrderByIdFn: func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
			return newTestOrderHelper(id, 99, repo.StatusPending), nil // customerID=99
		},
	}
	svc := &svc{q: mock}

	_, err := svc.CancelOrder(context.Background(), 1, 10) // caller is 10, not 99
	if !errors.Is(err, ErrOrderForbidden) {
		t.Fatalf("expected ErrOrderForbidden, got %v", err)
	}
}

func TestServiceCancelOrder_NotPending(t *testing.T) {
	mock := &mockQuerier{
		findOrderByIdFn: func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
			return newTestOrderHelper(id, 10, repo.StatusCompleted), nil // already completed
		},
	}
	svc := &svc{q: mock}

	_, err := svc.CancelOrder(context.Background(), 1, 10)
	if !errors.Is(err, ErrOrderNotPending) {
		t.Fatalf("expected ErrOrderNotPending, got %v", err)
	}
}

func TestServiceUpdateOrderStatus_InvalidStatus(t *testing.T) {
	svc := &svc{q: &mockQuerier{}}

	_, err := svc.UpdateOrderStatus(context.Background(), 1, "invalid-status")
	if !errors.Is(err, ErrInvalidStatus) {
		t.Fatalf("expected ErrInvalidStatus, got %v", err)
	}
}

func TestServiceUpdateOrderStatus_NotFound(t *testing.T) {
	mock := &mockQuerier{
		findOrderByIdFn: func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
			return nil, nil
		},
	}
	svc := &svc{q: mock}

	_, err := svc.UpdateOrderStatus(context.Background(), 999, "completed")
	if !errors.Is(err, ErrOrderNotFound) {
		t.Fatalf("expected ErrOrderNotFound, got %v", err)
	}
}

func TestServiceUpdateOrderStatus_Success(t *testing.T) {
	callCount := 0
	mock := &mockQuerier{
		findOrderByIdFn: func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
			callCount++
			return newTestOrderHelper(id, 10, repo.StatusCompleted), nil
		},
		updateOrderStatusFn: func(ctx context.Context, arg repo.UpdateOrderStatusParams) (repo.Order, error) {
			return repo.Order{ID: arg.ID, Status: arg.Status}, nil
		},
	}
	svc := &svc{q: mock}

	order, err := svc.UpdateOrderStatus(context.Background(), 1, "completed")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if order.ID != 1 {
		t.Errorf("expected id=1, got %d", order.ID)
	}
}

func TestOrderItemValidate_Valid(t *testing.T) {
	item := orderItem{ProductID: 1, Quantity: 2}
	if err := item.validate(); err != nil {
		t.Errorf("expected no error, got %v", err)
	}
}

func TestOrderItemValidate_ZeroProductID(t *testing.T) {
	item := orderItem{ProductID: 0, Quantity: 2}
	if err := item.validate(); !errors.Is(err, ErrOrderItemValidation) {
		t.Errorf("expected ErrOrderItemValidation, got %v", err)
	}
}

func TestOrderItemValidate_NegativeQuantity(t *testing.T) {
	item := orderItem{ProductID: 1, Quantity: -1}
	if err := item.validate(); !errors.Is(err, ErrOrderItemValidation) {
		t.Errorf("expected ErrOrderItemValidation, got %v", err)
	}
}

func TestOrderItemValidate_ZeroQuantity(t *testing.T) {
	item := orderItem{ProductID: 1, Quantity: 0}
	if err := item.validate(); !errors.Is(err, ErrOrderItemValidation) {
		t.Errorf("expected ErrOrderItemValidation, got %v", err)
	}
}

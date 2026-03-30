package orders

import (
	"context"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"example.com/ecommerce/internal/testhelper"
)

func TestIntegrationPlaceOrder_Success(t *testing.T) {
	pool := testhelper.NewTestDB(t)
	ctx := context.Background()
	queries := repo.New(pool)
	svc := NewService(queries, pool)

	// テスト用ユーザー作成
	user, err := queries.CreateUser(ctx, repo.CreateUserParams{
		Name:         "Order User",
		Email:        "order@example.com",
		PasswordHash: "hashed",
		Role:         repo.UserRoleUser,
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// テスト用商品作成（量はSQLで直接設定）
	product, err := queries.CreateProduct(ctx, repo.CreateProductParams{
		Name:         "Test Product",
		PriceInCents: 1000,
	})
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	// quantity を直接 UPDATE して在庫を設定
	_, err = pool.Exec(ctx, "UPDATE products SET quantity = 10 WHERE id = $1", product.ID)
	if err != nil {
		t.Fatalf("failed to set product quantity: %v", err)
	}

	// 注文作成
	order, err := svc.PlaceOrder(ctx, createOrderParams{
		CustomerID: user.ID,
		Items: []orderItem{
			{ProductID: product.ID, Quantity: 2},
		},
	})
	if err != nil {
		t.Fatalf("PlaceOrder failed: %v", err)
	}
	if order.CustomerID != user.ID {
		t.Errorf("expected customer_id=%d, got %d", user.ID, order.CustomerID)
	}
	if order.Status != repo.StatusPending {
		t.Errorf("expected status=pending, got %s", order.Status)
	}
}

func TestIntegrationPlaceOrder_InsufficientStock(t *testing.T) {
	pool := testhelper.NewTestDB(t)
	ctx := context.Background()
	queries := repo.New(pool)
	svc := NewService(queries, pool)

	user, err := queries.CreateUser(ctx, repo.CreateUserParams{
		Name:         "Stock User",
		Email:        "stock@example.com",
		PasswordHash: "hashed",
		Role:         repo.UserRoleUser,
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	// 在庫 1 の商品
	product, err := queries.CreateProduct(ctx, repo.CreateProductParams{
		Name:         "Low Stock Product",
		PriceInCents: 500,
	})
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	// quantity を 1 に設定
	_, err = pool.Exec(ctx, "UPDATE products SET quantity = 1 WHERE id = $1", product.ID)
	if err != nil {
		t.Fatalf("failed to set product quantity: %v", err)
	}

	// 在庫を超える数量で注文
	_, err = svc.PlaceOrder(ctx, createOrderParams{
		CustomerID: user.ID,
		Items: []orderItem{
			{ProductID: product.ID, Quantity: 5},
		},
	})
	if err == nil {
		t.Fatal("expected error for insufficient stock, got nil")
	}
}

func TestIntegrationPlaceOrder_EmptyItems(t *testing.T) {
	pool := testhelper.NewTestDB(t)
	ctx := context.Background()
	queries := repo.New(pool)
	svc := NewService(queries, pool)

	_, err := svc.PlaceOrder(ctx, createOrderParams{
		CustomerID: 1,
		Items:      []orderItem{},
	})
	if err == nil {
		t.Fatal("expected error for empty items, got nil")
	}
}

func TestIntegrationCancelOrder_Success(t *testing.T) {
	pool := testhelper.NewTestDB(t)
	ctx := context.Background()
	queries := repo.New(pool)
	svc := NewService(queries, pool)

	user, err := queries.CreateUser(ctx, repo.CreateUserParams{
		Name:         "Cancel User",
		Email:        "cancel@example.com",
		PasswordHash: "hashed",
		Role:         repo.UserRoleUser,
	})
	if err != nil {
		t.Fatalf("CreateUser failed: %v", err)
	}

	product, err := queries.CreateProduct(ctx, repo.CreateProductParams{
		Name:         "Cancel Product",
		PriceInCents: 1000,
	})
	if err != nil {
		t.Fatalf("CreateProduct failed: %v", err)
	}

	// quantity を 10 に設定
	_, err = pool.Exec(ctx, "UPDATE products SET quantity = 10 WHERE id = $1", product.ID)
	if err != nil {
		t.Fatalf("failed to set product quantity: %v", err)
	}

	order, err := svc.PlaceOrder(ctx, createOrderParams{
		CustomerID: user.ID,
		Items:      []orderItem{{ProductID: product.ID, Quantity: 1}},
	})
	if err != nil {
		t.Fatalf("PlaceOrder failed: %v", err)
	}

	_, err = svc.CancelOrder(ctx, order.ID, user.ID)
	if err != nil {
		t.Fatalf("CancelOrder failed: %v", err)
	}
}

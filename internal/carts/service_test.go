package carts

import (
	"context"
	"errors"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

func TestServiceCreateCart(t *testing.T) {
	mock := &mockQuerier{
		createCartFn: func(ctx context.Context, userID int64) (repo.Cart, error) {
			return newTestCartHelper(1, userID), nil
		},
	}
	svc := NewService(mock)

	cart, err := svc.CreateCart(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cart.UserID != 10 {
		t.Errorf("expected user_id=10, got %d", cart.UserID)
	}
}

func TestServiceAddItemToCart_Success(t *testing.T) {
	mock := &mockQuerier{
		findCartByUserIdFn: func(ctx context.Context, userID int64) (repo.Cart, error) {
			return newTestCartHelper(1, userID), nil
		},
		addItemToCartFn: func(ctx context.Context, arg repo.AddItemToCartParams) (repo.CartItem, error) {
			return newTestCartItemHelper(1, arg.CartID, arg.ProductID), nil
		},
	}
	svc := NewService(mock)

	item, err := svc.AddItemToCart(context.Background(), 10, 5, 3)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ProductID != 5 {
		t.Errorf("expected product_id=5, got %d", item.ProductID)
	}
}

func TestServiceAddItemToCart_CartNotFound(t *testing.T) {
	mock := &mockQuerier{
		findCartByUserIdFn: func(ctx context.Context, userID int64) (repo.Cart, error) {
			return repo.Cart{}, pgx.ErrNoRows
		},
	}
	svc := NewService(mock)

	_, err := svc.AddItemToCart(context.Background(), 10, 5, 3)
	if !errors.Is(err, ErrCartNotFound) {
		t.Fatalf("expected ErrCartNotFound, got %v", err)
	}
}

func TestServiceListCartItemsByUserId(t *testing.T) {
	mock := &mockQuerier{
		listCartItemsByUserIdFn: func(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
			return []repo.ListCartItemsByUserIdRow{
				{ID: 1, ProductID: 1, Quantity: 2, ProductPriceInCents: 1000},
				{ID: 2, ProductID: 2, Quantity: 1, ProductPriceInCents: 2000},
			}, nil
		},
	}
	svc := NewService(mock)

	items, err := svc.ListCartItemsByUserId(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(items) != 2 {
		t.Errorf("expected 2 items, got %d", len(items))
	}
}

func TestServiceUpdateCartItemQuantity_Success(t *testing.T) {
	mock := &mockQuerier{
		findCartItemByIdFn: func(ctx context.Context, id int64) (repo.CartItem, error) {
			return newTestCartItemHelper(id, 1, 5), nil // cart_id=1
		},
		findCartByUserIdFn: func(ctx context.Context, userID int64) (repo.Cart, error) {
			return newTestCartHelper(1, userID), nil // id=1 matches
		},
		updateCartItemQuantityFn: func(ctx context.Context, arg repo.UpdateCartItemQuantityParams) (repo.CartItem, error) {
			item := newTestCartItemHelper(arg.ID, 1, 5)
			item.Quantity = arg.Quantity
			return item, nil
		},
	}
	svc := NewService(mock)

	item, err := svc.UpdateCartItemQuantity(context.Background(), 10, 1, 5)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.Quantity != 5 {
		t.Errorf("expected quantity=5, got %d", item.Quantity)
	}
}

func TestServiceUpdateCartItemQuantity_Forbidden(t *testing.T) {
	mock := &mockQuerier{
		findCartItemByIdFn: func(ctx context.Context, id int64) (repo.CartItem, error) {
			return newTestCartItemHelper(id, 99, 5), nil // cart_id=99 (他人のカート)
		},
		findCartByUserIdFn: func(ctx context.Context, userID int64) (repo.Cart, error) {
			return newTestCartHelper(1, userID), nil // id=1, mismatch with 99
		},
	}
	svc := NewService(mock)

	_, err := svc.UpdateCartItemQuantity(context.Background(), 10, 1, 5)
	if !errors.Is(err, ErrCartForbidden) {
		t.Fatalf("expected ErrCartForbidden, got %v", err)
	}
}

func TestServiceRemoveItemFromCart_Success(t *testing.T) {
	mock := &mockQuerier{
		findCartItemByIdFn: func(ctx context.Context, id int64) (repo.CartItem, error) {
			return newTestCartItemHelper(id, 1, 5), nil // cart_id=1
		},
		findCartByUserIdFn: func(ctx context.Context, userID int64) (repo.Cart, error) {
			return newTestCartHelper(1, userID), nil // id=1 matches
		},
		removeItemFromCartFn: func(ctx context.Context, id int64) (repo.CartItem, error) {
			return newTestCartItemHelper(id, 1, 5), nil
		},
	}
	svc := NewService(mock)

	item, err := svc.RemoveItemFromCart(context.Background(), 10, 1)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if item.ID != 1 {
		t.Errorf("expected id=1, got %d", item.ID)
	}
}

func TestServiceRemoveItemFromCart_Forbidden(t *testing.T) {
	mock := &mockQuerier{
		findCartItemByIdFn: func(ctx context.Context, id int64) (repo.CartItem, error) {
			return newTestCartItemHelper(id, 99, 5), nil // 他人のカート
		},
		findCartByUserIdFn: func(ctx context.Context, userID int64) (repo.Cart, error) {
			return newTestCartHelper(1, userID), nil
		},
	}
	svc := NewService(mock)

	_, err := svc.RemoveItemFromCart(context.Background(), 10, 1)
	if !errors.Is(err, ErrCartForbidden) {
		t.Fatalf("expected ErrCartForbidden, got %v", err)
	}
}

func TestServiceClearCart(t *testing.T) {
	mock := &mockQuerier{
		clearCartFn: func(ctx context.Context, userID int64) error {
			return nil
		},
	}
	svc := NewService(mock)

	err := svc.ClearCart(context.Background(), 10)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestServiceClearCart_Error(t *testing.T) {
	mock := &mockQuerier{
		clearCartFn: func(ctx context.Context, userID int64) error {
			return errors.New("db error")
		},
	}
	svc := NewService(mock)

	err := svc.ClearCart(context.Background(), 10)
	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

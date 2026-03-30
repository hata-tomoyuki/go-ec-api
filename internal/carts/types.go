package carts

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

var (
	ErrCartNotFound  = errors.New("cart not found")
	ErrCartForbidden = errors.New("you do not have permission to access this cart")
)

type Service interface {
	CreateCart(ctx context.Context, userID int64) (repo.Cart, error)
	AddItemToCart(ctx context.Context, userID int64, productID int64, quantity int) (repo.CartItem, error)
	ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error)
	UpdateCartItemQuantity(ctx context.Context, userID int64, cartItemID int64, quantity int) (repo.CartItem, error)
	RemoveItemFromCart(ctx context.Context, userID int64, cartItemID int64) (repo.CartItem, error)
	ClearCart(ctx context.Context, userID int64) error
}

type addItemToCartParams struct {
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

type updateCartItemParams struct {
	CartItemID int64 `json:"cart_item_id"`
	Quantity   int   `json:"quantity"`
}

package carts

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type Service interface {
	CreateCart(ctx context.Context, userID int64) (repo.Cart, error)
	AddItemToCart(ctx context.Context, cartID int64, productID int64, quantity int) (repo.CartItem, error)
	ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error)
}

type addItemToCartParams struct {
	CartID    int64 `json:"cart_id"`
	ProductID int64 `json:"product_id"`
	Quantity  int   `json:"quantity"`
}

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

func (p addItemToCartParams) validate() error {
	if p.ProductID <= 0 || p.Quantity <= 0 {
		return errors.New("product_id must be provided and quantity must be greater than 0")
	}
	return nil
}

type updateCartItemParams struct {
	Quantity int `json:"quantity"`
}

func (p updateCartItemParams) validate() error {
	if p.Quantity <= 0 {
		return errors.New("quantity must be greater than 0")
	}
	return nil
}

package carts

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) CreateCart(ctx context.Context, userID int64) (repo.Cart, error) {
	return s.repo.CreateCart(ctx, userID)
}

func (s *svc) AddItemToCart(ctx context.Context, cartID int64, productID int64, quantity int) (repo.CartItem, error) {
	item, err := s.repo.AddItemToCart(ctx, repo.AddItemToCartParams{
		CartID:    cartID,
		ProductID: productID,
		Quantity:  int32(quantity),
	})

	if err != nil {
		return repo.CartItem{}, err
	}

	return item, nil
}

func (s *svc) ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
	return s.repo.ListCartItemsByUserId(ctx, userID)
}

func (s *svc) UpdateCartItemQuantity(ctx context.Context, productID int64, quantity int) (repo.CartItem, error) {
	return s.repo.UpdateCartItemQuantity(ctx, repo.UpdateCartItemQuantityParams{
		ID:       productID,
		Quantity: int32(quantity),
	})
}

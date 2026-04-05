package carts

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
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

func (s *svc) AddItemToCart(ctx context.Context, userID int64, productID int64, quantity int) (repo.CartItem, error) {
	cart, err := s.repo.FindCartByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CartItem{}, ErrCartNotFound
		}
		return repo.CartItem{}, err
	}

	product, err := s.repo.FindProductById(ctx, productID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CartItem{}, ErrProductNotFound
		}
		return repo.CartItem{}, err
	}

	if product.Quantity < int32(quantity) {
		return repo.CartItem{}, ErrInsufficientStock
	}

	return s.repo.AddItemToCart(ctx, repo.AddItemToCartParams{
		CartID:    cart.ID,
		ProductID: productID,
		Quantity:  int32(quantity),
	})
}

func (s *svc) ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
	return s.repo.ListCartItemsByUserId(ctx, userID)
}

func (s *svc) UpdateCartItemQuantity(ctx context.Context, userID int64, cartItemID int64, quantity int) (repo.CartItem, error) {
	cartItem, err := s.repo.FindCartItemById(ctx, cartItemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CartItem{}, ErrCartNotFound
		}
		return repo.CartItem{}, err
	}

	cart, err := s.repo.FindCartByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CartItem{}, ErrCartNotFound
		}
		return repo.CartItem{}, err
	}

	if cartItem.CartID != cart.ID {
		return repo.CartItem{}, ErrCartForbidden
	}

	return s.repo.UpdateCartItemQuantity(ctx, repo.UpdateCartItemQuantityParams{
		ID:       cartItemID,
		Quantity: int32(quantity),
	})
}

func (s *svc) RemoveItemFromCart(ctx context.Context, userID int64, cartItemID int64) (repo.CartItem, error) {
	cartItem, err := s.repo.FindCartItemById(ctx, cartItemID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CartItem{}, ErrCartNotFound
		}
		return repo.CartItem{}, err
	}

	cart, err := s.repo.FindCartByUserId(ctx, userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.CartItem{}, ErrCartNotFound
		}
		return repo.CartItem{}, err
	}

	if cartItem.CartID != cart.ID {
		return repo.CartItem{}, ErrCartForbidden
	}

	return s.repo.RemoveItemFromCart(ctx, cartItemID)
}

func (s *svc) ClearCart(ctx context.Context, userID int64) error {
	return s.repo.ClearCart(ctx, userID)
}

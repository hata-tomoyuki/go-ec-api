package orders

import (
	"context"
	"errors"
	"fmt"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

var (
	ErrorProductNotFound = errors.New("product not found")
	ErrorProductNoStock  = errors.New("product out of stock")
)

type svc struct {
	repo *repo.Queries
	db   *pgx.Conn
}

func NewService(repo *repo.Queries, db *pgx.Conn) Service {
	return &svc{
		repo: repo,
		db:   db,
	}
}

func (s *svc) ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error) {
	return s.repo.ListOrdersByCustomerID(ctx, customerID)
}

func (s *svc) ListAllOrders(ctx context.Context) ([]repo.ListAllOrdersRow, error) {
	return s.repo.ListAllOrders(ctx)
}

func (s *svc) FindOrderById(ctx context.Context, orderID int64) (repo.FindOrderByIdRow, error) {
	return s.repo.FindOrderById(ctx, orderID)
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, fmt.Errorf("at least one order item is required")
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		return repo.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.repo.WithTx(tx)

	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		return repo.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range tempOrder.Items {
		product, err := qtx.FindProductById(ctx, item.ProductID)
		if err != nil {
			return repo.Order{}, ErrorProductNotFound
		}

		if product.Quantity < item.Quantity {
			return repo.Order{}, ErrorProductNoStock
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:      order.ID,
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PriceInCents: product.PriceInCents,
		})
		if err != nil {
			return repo.Order{}, fmt.Errorf("failed to create order item: %w", err)
		}
	}
	tx.Commit(ctx)

	return order, nil
}

func (s *svc) CancelOrder(ctx context.Context, orderID int64) (repo.FindOrderByIdRow, error) {
	order, err := s.repo.FindOrderById(ctx, orderID)
	if err != nil {
		return repo.FindOrderByIdRow{}, fmt.Errorf("failed to find order: %w", err)
	}

	if order.Status != "pending" {
		return repo.FindOrderByIdRow{}, fmt.Errorf("only pending orders can be cancelled")
	}

	_, err = s.repo.CancelOrder(ctx, orderID)
	if err != nil {
		return repo.FindOrderByIdRow{}, fmt.Errorf("failed to cancel order: %w", err)
	}

	return order, nil
}

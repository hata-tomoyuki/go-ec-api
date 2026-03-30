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

// pgBeginner はトランザクションを開始できる任意の接続型を抽象化する。
// *pgx.Conn と *pgxpool.Pool の両方がこの interface を満たす。
type pgBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type svc struct {
	repo *repo.Queries
	db   pgBeginner
}

func NewService(repo *repo.Queries, db pgBeginner) Service {
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
	order, err := s.repo.FindOrderById(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.FindOrderByIdRow{}, ErrOrderNotFound
		}
		return repo.FindOrderByIdRow{}, err
	}
	return order, nil
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
	if err := tx.Commit(ctx); err != nil {
		return repo.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}

	return order, nil
}

func (s *svc) CancelOrder(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error) {
	order, err := s.repo.FindOrderById(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.FindOrderByIdRow{}, ErrOrderNotFound
		}
		return repo.FindOrderByIdRow{}, err
	}

	if order.CustomerID != customerID {
		return repo.FindOrderByIdRow{}, ErrOrderForbidden
	}

	if order.Status != "pending" {
		return repo.FindOrderByIdRow{}, ErrOrderNotPending
	}

	_, err = s.repo.CancelOrder(ctx, orderID)
	if err != nil {
		return repo.FindOrderByIdRow{}, fmt.Errorf("failed to cancel order: %w", err)
	}

	return order, nil
}

func isValidStatus(status string) bool {
	switch repo.Status(status) {
	case repo.StatusPending, repo.StatusCompleted, repo.StatusCancelled:
		return true
	}
	return false
}

func (s *svc) UpdateOrderStatus(ctx context.Context, orderID int64, status string) (repo.FindOrderByIdRow, error) {
	if !isValidStatus(status) {
		return repo.FindOrderByIdRow{}, ErrInvalidStatus
	}

	_, err := s.repo.FindOrderById(ctx, orderID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.FindOrderByIdRow{}, ErrOrderNotFound
		}
		return repo.FindOrderByIdRow{}, err
	}

	_, err = s.repo.UpdateOrderStatus(ctx, repo.UpdateOrderStatusParams{
		ID:     orderID,
		Status: repo.Status(status),
	})
	if err != nil {
		return repo.FindOrderByIdRow{}, fmt.Errorf("failed to update order status: %w", err)
	}

	return s.repo.FindOrderById(ctx, orderID)
}

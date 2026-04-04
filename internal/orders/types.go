package orders

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

var (
	ErrOrderNotFound      = errors.New("order not found")
	ErrOrderNotPending    = errors.New("only pending orders can be cancelled")
	ErrOrderForbidden     = errors.New("you do not have permission to access this order")
	ErrInvalidStatus      = errors.New("invalid status: must be one of pending, completed, cancelled")
	ErrOrderItemValidation = errors.New("each item must have a valid product_id and quantity greater than 0")
	ErrOrderEmptyItems    = errors.New("at least one order item is required")
)

type orderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

func (i orderItem) validate() error {
	if i.ProductID <= 0 || i.Quantity <= 0 {
		return ErrOrderItemValidation
	}
	return nil
}

type createOrderParams struct {
	CustomerID int64       `json:"customer_id"`
	Items      []orderItem `json:"items"`
}

type Service interface {
	ListAllOrders(ctx context.Context) ([]repo.ListAllOrdersRow, error)
	ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error)
	FindOrderById(ctx context.Context, orderID int64) ([]repo.FindOrderByIdRow, error)
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
	CancelOrder(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) (repo.FindOrderByIdRow, error)
}

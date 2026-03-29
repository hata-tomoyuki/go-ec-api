package orders

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type orderItem struct {
	ProductID int64 `json:"product_id"`
	Quantity  int32 `json:"quantity"`
}

type createOrderParams struct {
	CustomerID int64       `json:"customer_id"`
	Items      []orderItem `json:"items"`
}

type Service interface {
	ListAllOrders(ctx context.Context) ([]repo.ListAllOrdersRow, error)
	ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error)
	FindOrderById(ctx context.Context, orderID int64) (repo.FindOrderByIdRow, error)
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
	CancelOrder(ctx context.Context, orderID int64) (repo.FindOrderByIdRow, error)
}

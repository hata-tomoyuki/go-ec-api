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
	ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error)
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
}

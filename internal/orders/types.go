package orders

import (
	"context"
	"errors"
	"time"

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

var allowedOrderSorts = map[string]string{
	"created_at_desc": "o.created_at DESC",
	"created_at_asc":  "o.created_at ASC",
	"status_asc":      "o.status ASC",
	"status_desc":     "o.status DESC",
}

type listOrdersParams struct {
	Page   int
	Limit  int
	Sort   string
	Status string // フィルタ: pending, completed, cancelled
}

func (p *listOrdersParams) validate() error {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		return ErrOrderItemValidation
	}
	if p.Sort == "" {
		p.Sort = "created_at_desc"
	}
	if _, ok := allowedOrderSorts[p.Sort]; !ok {
		return ErrInvalidStatus
	}
	if p.Status != "" && !isValidStatusFilter(p.Status) {
		return ErrInvalidStatus
	}
	return nil
}

func isValidStatusFilter(status string) bool {
	switch status {
	case "pending", "completed", "cancelled":
		return true
	}
	return false
}

type paginatedOrders struct {
	Data  []paginatedOrderRow `json:"data"`
	Total int                 `json:"total"`
	Page  int                 `json:"page"`
	Limit int                 `json:"limit"`
}

type paginatedOrderRow struct {
	ID           int64  `json:"id"`
	CustomerID   int64  `json:"customer_id"`
	Status       string `json:"status"`
	ItemCount    int    `json:"item_count"`
	TotalInCents int64  `json:"total_in_cents"`
	CreatedAt    time.Time `json:"created_at"`
	UpdatedAt    time.Time `json:"updated_at"`
	TotalCount   int    `json:"-"`
}

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
	ListAllOrdersPaginated(ctx context.Context, params listOrdersParams) (paginatedOrders, error)
	ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error)
	ListOrdersPaginated(ctx context.Context, customerID int64, params listOrdersParams) (paginatedOrders, error)
	FindOrderById(ctx context.Context, orderID int64) ([]repo.FindOrderByIdRow, error)
	PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error)
	CancelOrder(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error)
	UpdateOrderStatus(ctx context.Context, orderID int64, status string) (repo.FindOrderByIdRow, error)
}

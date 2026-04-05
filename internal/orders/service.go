package orders

import (
	"context"
	"errors"
	"fmt"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
)

var tracer = otel.Tracer("orders")


var (
	ErrorProductNotFound = errors.New("product not found")
	ErrorProductNoStock  = errors.New("product out of stock")
)

type pgBeginner interface {
	Begin(ctx context.Context) (pgx.Tx, error)
}

type svc struct {
	q      repo.Querier
	db     pgBeginner
	newTxQ func(pgx.Tx) repo.Querier
}

func NewService(q repo.Querier, db pgBeginner) Service {
	return &svc{
		q:      q,
		db:     db,
		newTxQ: func(tx pgx.Tx) repo.Querier { return repo.New(tx) },
	}
}

func (s *svc) ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error) {
	return s.q.ListOrdersByCustomerID(ctx, customerID)
}

func (s *svc) ListAllOrders(ctx context.Context) ([]repo.ListAllOrdersRow, error) {
	return s.q.ListAllOrders(ctx)
}

func (s *svc) listOrdersPaginated(ctx context.Context, customerID *int64, params listOrdersParams) (paginatedOrders, error) {
	whereSQL := ""
	args := []any{}
	argN := 1

	if customerID != nil {
		whereSQL = fmt.Sprintf(" WHERE o.customer_id = $%d", argN)
		args = append(args, *customerID)
		argN++
	}

	if params.Status != "" {
		if whereSQL == "" {
			whereSQL = fmt.Sprintf(" WHERE o.status = $%d", argN)
		} else {
			whereSQL += fmt.Sprintf(" AND o.status = $%d", argN)
		}
		args = append(args, params.Status)
		argN++
	}

	sql := `
	SELECT o.id, o.customer_id, o.status,
	       COUNT(oi.id)::int AS item_count,
	       COALESCE(SUM(oi.price_in_cents * oi.quantity), 0)::bigint AS total_in_cents,
	       o.created_at, o.updated_at,
	       COUNT(*) OVER() AS total_count
	FROM orders o
	LEFT JOIN order_items oi ON o.id = oi.order_id
	` + whereSQL + `
	GROUP BY o.id`

	sortSQL := allowedOrderSorts[params.Sort]
	offset := (params.Page - 1) * params.Limit

	sql += fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", sortSQL, argN, argN+1)
	args = append(args, params.Limit, offset)

	rows, err := s.db.(repo.DBTX).Query(ctx, sql, args...)
	if err != nil {
		return paginatedOrders{}, err
	}
	defer rows.Close()

	orders := make([]paginatedOrderRow, 0)
	for rows.Next() {
		var o paginatedOrderRow
		if err := rows.Scan(&o.ID, &o.CustomerID, &o.Status, &o.ItemCount, &o.TotalInCents, &o.CreatedAt, &o.UpdatedAt, &o.TotalCount); err != nil {
			return paginatedOrders{}, err
		}
		orders = append(orders, o)
	}

	total := 0
	if len(orders) > 0 {
		total = orders[0].TotalCount
	}

	return paginatedOrders{
		Data:  orders,
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
	}, nil
}

func (s *svc) ListOrdersPaginated(ctx context.Context, customerID int64, params listOrdersParams) (paginatedOrders, error) {
	return s.listOrdersPaginated(ctx, &customerID, params)
}

func (s *svc) ListAllOrdersPaginated(ctx context.Context, params listOrdersParams) (paginatedOrders, error) {
	return s.listOrdersPaginated(ctx, nil, params)
}

func (s *svc) FindOrderById(ctx context.Context, orderID int64) ([]repo.FindOrderByIdRow, error) {
	rows, err := s.q.FindOrderById(ctx, orderID)
	if err != nil {
		return nil, err
	}
	if len(rows) == 0 {
		return nil, ErrOrderNotFound
	}
	return rows, nil
}

func (s *svc) PlaceOrder(ctx context.Context, tempOrder createOrderParams) (repo.Order, error) {
	ctx, span := tracer.Start(ctx, "orders.PlaceOrder")

	defer span.End()

	span.SetAttributes(
		attribute.Int64("customer_id", tempOrder.CustomerID),
		attribute.Int("item_count", len(tempOrder.Items)),
	)

	if tempOrder.CustomerID == 0 {
		return repo.Order{}, fmt.Errorf("customer ID is required")
	}
	if len(tempOrder.Items) == 0 {
		return repo.Order{}, ErrOrderEmptyItems
	}

	for _, item := range tempOrder.Items {
		if err := item.validate(); err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return repo.Order{}, err
		}
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return repo.Order{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.newTxQ(tx)

	order, err := qtx.CreateOrder(ctx, tempOrder.CustomerID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return repo.Order{}, fmt.Errorf("failed to create order: %w", err)
	}

	for _, item := range tempOrder.Items {
		product, err := qtx.DecrementProductQuantity(ctx, repo.DecrementProductQuantityParams{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		})
		if err != nil {
			if errors.Is(err, pgx.ErrNoRows) {
				span.RecordError(err)
				span.SetStatus(codes.Error, err.Error())
				return repo.Order{}, ErrorProductNoStock
			}
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return repo.Order{}, fmt.Errorf("failed to decrement stock: %w", err)
		}

		_, err = qtx.CreateOrderItem(ctx, repo.CreateOrderItemParams{
			OrderID:      order.ID,
			ProductID:    item.ProductID,
			Quantity:     item.Quantity,
			PriceInCents: product.PriceInCents,
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return repo.Order{}, fmt.Errorf("failed to create order item: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return repo.Order{}, fmt.Errorf("failed to commit transaction: %w", err)
	}
	span.SetAttributes(attribute.Int64("order_id", order.ID))
	return order, nil
}

func (s *svc) CancelOrder(ctx context.Context, orderID int64, customerID int64) (repo.FindOrderByIdRow, error) {
	ctx, span := tracer.Start(ctx, "orders.CancelOrder")
	defer span.End()

	span.SetAttributes(
		attribute.Int64("order_id", orderID),
		attribute.Int64("customer_id", customerID),
	)

	rows, err := s.q.FindOrderById(ctx, orderID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return repo.FindOrderByIdRow{}, err
	}
	if len(rows) == 0 {
		span.RecordError(ErrOrderNotFound)
		span.SetStatus(codes.Error, ErrOrderNotFound.Error())
		return repo.FindOrderByIdRow{}, ErrOrderNotFound
	}

	order := rows[0]
	if order.CustomerID != customerID {
		return repo.FindOrderByIdRow{}, ErrOrderForbidden
	}

	if order.Status != "pending" {
		return repo.FindOrderByIdRow{}, ErrOrderNotPending
	}

	tx, err := s.db.Begin(ctx)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return repo.FindOrderByIdRow{}, fmt.Errorf("failed to begin transaction: %w", err)
	}
	defer tx.Rollback(ctx)
	qtx := s.newTxQ(tx)

	_, err = qtx.CancelOrder(ctx, orderID)
	if err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return repo.FindOrderByIdRow{}, fmt.Errorf("failed to cancel order: %w", err)
	}

	for _, item := range rows {
		_, err = qtx.IncrementProductQuantity(ctx, repo.IncrementProductQuantityParams{
			ID:       item.ProductID,
			Quantity: item.Quantity,
		})
		if err != nil {
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
			return repo.FindOrderByIdRow{}, fmt.Errorf("failed to restore stock: %w", err)
		}
	}
	if err := tx.Commit(ctx); err != nil {
		span.RecordError(err)
		span.SetStatus(codes.Error, err.Error())
		return repo.FindOrderByIdRow{}, fmt.Errorf("failed to commit transaction: %w", err)
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

	rows, err := s.q.FindOrderById(ctx, orderID)
	if err != nil {
		return repo.FindOrderByIdRow{}, err
	}
	if len(rows) == 0 {
		return repo.FindOrderByIdRow{}, ErrOrderNotFound
	}

	_, err = s.q.UpdateOrderStatus(ctx, repo.UpdateOrderStatusParams{
		ID:     orderID,
		Status: repo.Status(status),
	})
	if err != nil {
		return repo.FindOrderByIdRow{}, fmt.Errorf("failed to update order status: %w", err)
	}

	updatedRows, err := s.q.FindOrderById(ctx, orderID)
	if err != nil {
		return repo.FindOrderByIdRow{}, err
	}
	if len(updatedRows) == 0 {
		return repo.FindOrderByIdRow{}, ErrOrderNotFound
	}
	return updatedRows[0], nil
}

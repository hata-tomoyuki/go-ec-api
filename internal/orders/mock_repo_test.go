package orders

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// mockQuerier は Querier interface の手動モック（orders テスト用）
type mockQuerier struct {
	listAllOrdersFn          func(ctx context.Context) ([]repo.ListAllOrdersRow, error)
	listOrdersByCustomerIDFn func(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error)
	findOrderByIdFn          func(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error)
	cancelOrderFn            func(ctx context.Context, id int64) (repo.Order, error)
	updateOrderStatusFn      func(ctx context.Context, arg repo.UpdateOrderStatusParams) (repo.Order, error)
}

func (m *mockQuerier) ListAllOrders(ctx context.Context) ([]repo.ListAllOrdersRow, error) {
	return m.listAllOrdersFn(ctx)
}
func (m *mockQuerier) ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error) {
	return m.listOrdersByCustomerIDFn(ctx, customerID)
}
func (m *mockQuerier) FindOrderById(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
	return m.findOrderByIdFn(ctx, id)
}
func (m *mockQuerier) CancelOrder(ctx context.Context, id int64) (repo.Order, error) {
	return m.cancelOrderFn(ctx, id)
}
func (m *mockQuerier) UpdateOrderStatus(ctx context.Context, arg repo.UpdateOrderStatusParams) (repo.Order, error) {
	return m.updateOrderStatusFn(ctx, arg)
}

// --- Querier interface を満たすためのスタブ ---

func (m *mockQuerier) AddItemToCart(ctx context.Context, arg repo.AddItemToCartParams) (repo.CartItem, error) {
	panic("not implemented")
}
func (m *mockQuerier) AddProductToCategory(ctx context.Context, arg repo.AddProductToCategoryParams) error {
	panic("not implemented")
}
func (m *mockQuerier) ClearCart(ctx context.Context, userID int64) error {
	panic("not implemented")
}
func (m *mockQuerier) ConsumeRefreshToken(ctx context.Context, tokenHash string) (repo.RefreshToken, error) {
	panic("not implemented")
}
func (m *mockQuerier) CreateAddress(ctx context.Context, arg repo.CreateAddressParams) (repo.Address, error) {
	panic("not implemented")
}
func (m *mockQuerier) CreateCart(ctx context.Context, userID int64) (repo.Cart, error) {
	panic("not implemented")
}
func (m *mockQuerier) CreateCategory(ctx context.Context, arg repo.CreateCategoryParams) (repo.Category, error) {
	panic("not implemented")
}
func (m *mockQuerier) CreateOrder(ctx context.Context, customerID int64) (repo.Order, error) {
	panic("not implemented")
}
func (m *mockQuerier) CreateOrderItem(ctx context.Context, arg repo.CreateOrderItemParams) (repo.OrderItem, error) {
	panic("not implemented")
}
func (m *mockQuerier) CreateProduct(ctx context.Context, arg repo.CreateProductParams) (repo.Product, error) {
	panic("not implemented")
}
func (m *mockQuerier) CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
	panic("not implemented")
}
func (m *mockQuerier) DeleteAddress(ctx context.Context, id int32) error {
	panic("not implemented")
}
func (m *mockQuerier) DeleteCategory(ctx context.Context, id int64) (repo.Category, error) {
	panic("not implemented")
}
func (m *mockQuerier) DeleteExpiredRefreshTokens(ctx context.Context) error {
	panic("not implemented")
}
func (m *mockQuerier) DeleteProduct(ctx context.Context, id int64) (repo.Product, error) {
	panic("not implemented")
}
func (m *mockQuerier) DeleteRefreshToken(ctx context.Context, id int64) error {
	panic("not implemented")
}
func (m *mockQuerier) DeleteRefreshTokensByUserId(ctx context.Context, userID int64) error {
	panic("not implemented")
}
func (m *mockQuerier) FindAddressById(ctx context.Context, id int32) (repo.Address, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindCartByUserId(ctx context.Context, userID int64) (repo.Cart, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindCartItemById(ctx context.Context, id int64) (repo.CartItem, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindProductById(ctx context.Context, id int64) (repo.FindProductByIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindUserByEmail(ctx context.Context, email string) (repo.User, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindUserById(ctx context.Context, id int64) (repo.User, error) {
	panic("not implemented")
}
func (m *mockQuerier) InsertRefreshToken(ctx context.Context, arg repo.InsertRefreshTokenParams) (repo.RefreshToken, error) {
	panic("not implemented")
}
func (m *mockQuerier) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListAddressesByUserId(ctx context.Context, userID int64) ([]repo.Address, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListCategories(ctx context.Context) ([]repo.ListCategoriesRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListProducts(ctx context.Context) ([]repo.ListProductsRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListProductsByCategory(ctx context.Context, categoryID int64) ([]repo.ListProductsByCategoryRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) RemoveItemFromCart(ctx context.Context, id int64) (repo.CartItem, error) {
	panic("not implemented")
}
func (m *mockQuerier) RemoveProductFromCategory(ctx context.Context, arg repo.RemoveProductFromCategoryParams) error {
	panic("not implemented")
}
func (m *mockQuerier) RevokeToken(ctx context.Context, arg repo.RevokeTokenParams) error {
	panic("not implemented")
}
func (m *mockQuerier) UpdateAddress(ctx context.Context, arg repo.UpdateAddressParams) (repo.Address, error) {
	panic("not implemented")
}
func (m *mockQuerier) UpdateCartItemQuantity(ctx context.Context, arg repo.UpdateCartItemQuantityParams) (repo.CartItem, error) {
	panic("not implemented")
}
func (m *mockQuerier) UpdateCategory(ctx context.Context, arg repo.UpdateCategoryParams) (repo.Category, error) {
	panic("not implemented")
}
func (m *mockQuerier) UpdateProduct(ctx context.Context, arg repo.UpdateProductParams) (repo.Product, error) {
	panic("not implemented")
}
func (m *mockQuerier) UpdateUser(ctx context.Context, arg repo.UpdateUserParams) (repo.User, error) {
	panic("not implemented")
}
func (m *mockQuerier) UpdateUserPassword(ctx context.Context, arg repo.UpdateUserPasswordParams) (repo.User, error) {
	panic("not implemented")
}

// テスト用ヘルパー
func newTestOrderHelper(id int64, customerID int64, status repo.Status) []repo.FindOrderByIdRow {
	return []repo.FindOrderByIdRow{
		{
			ID:           id,
			CustomerID:   customerID,
			Status:       status,
			CreatedAt:    pgtype.Timestamptz{Valid: true},
			UpdatedAt:    pgtype.Timestamptz{Valid: true},
			ProductID:    1,
			Quantity:     2,
			PriceInCents: 1000,
		},
	}
}

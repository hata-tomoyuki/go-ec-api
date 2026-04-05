package products

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgconn"
	"github.com/jackc/pgx/v5/pgtype"
)

// mockQuerier は Querier interface の手動モック（products テスト用）
type mockQuerier struct {
	listProductsFn    func(ctx context.Context) ([]repo.ListProductsRow, error)
	findProductByIdFn func(ctx context.Context, id int64) (repo.FindProductByIdRow, error)
	createProductFn   func(ctx context.Context, arg repo.CreateProductParams) (repo.Product, error)
	updateProductFn   func(ctx context.Context, arg repo.UpdateProductParams) (repo.Product, error)
	deleteProductFn   func(ctx context.Context, id int64) (repo.Product, error)
}

func (m *mockQuerier) ListProducts(ctx context.Context) ([]repo.ListProductsRow, error) {
	return m.listProductsFn(ctx)
}
func (m *mockQuerier) FindProductById(ctx context.Context, id int64) (repo.FindProductByIdRow, error) {
	return m.findProductByIdFn(ctx, id)
}
func (m *mockQuerier) CreateProduct(ctx context.Context, arg repo.CreateProductParams) (repo.Product, error) {
	return m.createProductFn(ctx, arg)
}
func (m *mockQuerier) UpdateProduct(ctx context.Context, arg repo.UpdateProductParams) (repo.Product, error) {
	return m.updateProductFn(ctx, arg)
}
func (m *mockQuerier) DeleteProduct(ctx context.Context, id int64) (repo.Product, error) {
	return m.deleteProductFn(ctx, id)
}

// --- Querier interface を満たすためのスタブ ---

func (m *mockQuerier) AddItemToCart(ctx context.Context, arg repo.AddItemToCartParams) (repo.CartItem, error) {
	panic("not implemented")
}
func (m *mockQuerier) AddProductToCategory(ctx context.Context, arg repo.AddProductToCategoryParams) error {
	panic("not implemented")
}
func (m *mockQuerier) CancelOrder(ctx context.Context, id int64) (repo.Order, error) {
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
func (m *mockQuerier) CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
	panic("not implemented")
}
func (m *mockQuerier) DecrementProductQuantity(ctx context.Context, arg repo.DecrementProductQuantityParams) (repo.Product, error) {
	panic("not implemented")
}
func (m *mockQuerier) DeleteAddress(ctx context.Context, arg repo.DeleteAddressParams) (repo.Address, error) {
	panic("not implemented")
}
func (m *mockQuerier) DeleteCategory(ctx context.Context, id int64) (repo.Category, error) {
	panic("not implemented")
}
func (m *mockQuerier) DeleteExpiredRefreshTokens(ctx context.Context) error {
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
func (m *mockQuerier) FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindOrderById(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindUserByEmail(ctx context.Context, email string) (repo.User, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindUserById(ctx context.Context, id int64) (repo.User, error) {
	panic("not implemented")
}
func (m *mockQuerier) IncrementProductQuantity(ctx context.Context, arg repo.IncrementProductQuantityParams) (repo.Product, error) {
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
func (m *mockQuerier) ListAllOrders(ctx context.Context) ([]repo.ListAllOrdersRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListCategories(ctx context.Context) ([]repo.ListCategoriesRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListProductsByCategory(ctx context.Context, categoryID int64) ([]repo.ListProductsByCategoryRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindCartByUserId(ctx context.Context, userID int64) (repo.Cart, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindCartItemById(ctx context.Context, id int64) (repo.CartItem, error) {
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
func (m *mockQuerier) UpdateOrderStatus(ctx context.Context, arg repo.UpdateOrderStatusParams) (repo.Order, error) {
	panic("not implemented")
}
func (m *mockQuerier) UpdateUser(ctx context.Context, arg repo.UpdateUserParams) (repo.User, error) {
	panic("not implemented")
}
func (m *mockQuerier) UpdateUserPassword(ctx context.Context, arg repo.UpdateUserPasswordParams) (repo.User, error) {
	panic("not implemented")
}

// テスト用ヘルパー
func newTestProduct(id int64, name string, price int32) repo.Product {
	return repo.Product{
		ID:           id,
		Name:         name,
		PriceInCents: price,
		Quantity:     10,
		CreatedAt:    pgtype.Timestamptz{Valid: true},
	}
}

func newTestListProductsRow(id int64, name string, price int32) repo.ListProductsRow {
	return repo.ListProductsRow{
		ID:           id,
		Name:         name,
		PriceInCents: price,
		Quantity:     10,
		ImageColor:   "from-gray-400 to-gray-600",
		CreatedAt:    pgtype.Timestamptz{Valid: true},
	}
}

func newTestFindProductByIdRow(id int64, name string, price int32) repo.FindProductByIdRow {
	return repo.FindProductByIdRow{
		ID:           id,
		Name:         name,
		PriceInCents: price,
		Quantity:     10,
		ImageColor:   "from-gray-400 to-gray-600",
		CreatedAt:    pgtype.Timestamptz{Valid: true},
	}
}

func newTestPaginatedProductRow(id int64, name string, price int32) paginatedProductRow {
	return paginatedProductRow{
		ID:           id,
		Name:         name,
		PriceInCents: price,
		Quantity:     10,
		ImageColor:   "from-gray-400 to-gray-600",
		CreatedAt:    pgtype.Timestamptz{Valid: true},
	}
}

// mockDBTX は repo.DBTX interface のモック（テスト用）
type mockDBTX struct {
	execFn     func(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error)
	queryFn    func(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error)
	queryRowFn func(ctx context.Context, sql string, args ...interface{}) pgx.Row
}

func (m *mockDBTX) Exec(ctx context.Context, sql string, args ...interface{}) (pgconn.CommandTag, error) {
	if m.execFn != nil {
		return m.execFn(ctx, sql, args...)
	}
	return pgconn.CommandTag{}, nil
}

func (m *mockDBTX) Query(ctx context.Context, sql string, args ...interface{}) (pgx.Rows, error) {
	if m.queryFn != nil {
		return m.queryFn(ctx, sql, args...)
	}
	return nil, nil
}

func (m *mockDBTX) QueryRow(ctx context.Context, sql string, args ...interface{}) pgx.Row {
	if m.queryRowFn != nil {
		return m.queryRowFn(ctx, sql, args...)
	}
	return nil
}

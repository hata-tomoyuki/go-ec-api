package carts

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// mockQuerier は Querier interface の手動モック（carts テスト用）
type mockQuerier struct {
	findCartByUserIdFn       func(ctx context.Context, userID int64) (repo.Cart, error)
	findCartItemByIdFn       func(ctx context.Context, id int64) (repo.CartItem, error)
	createCartFn             func(ctx context.Context, userID int64) (repo.Cart, error)
	addItemToCartFn          func(ctx context.Context, arg repo.AddItemToCartParams) (repo.CartItem, error)
	listCartItemsByUserIdFn  func(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error)
	updateCartItemQuantityFn func(ctx context.Context, arg repo.UpdateCartItemQuantityParams) (repo.CartItem, error)
	removeItemFromCartFn     func(ctx context.Context, id int64) (repo.CartItem, error)
	clearCartFn              func(ctx context.Context, userID int64) error
}

func (m *mockQuerier) FindCartByUserId(ctx context.Context, userID int64) (repo.Cart, error) {
	return m.findCartByUserIdFn(ctx, userID)
}
func (m *mockQuerier) FindCartItemById(ctx context.Context, id int64) (repo.CartItem, error) {
	return m.findCartItemByIdFn(ctx, id)
}
func (m *mockQuerier) CreateCart(ctx context.Context, userID int64) (repo.Cart, error) {
	return m.createCartFn(ctx, userID)
}
func (m *mockQuerier) AddItemToCart(ctx context.Context, arg repo.AddItemToCartParams) (repo.CartItem, error) {
	return m.addItemToCartFn(ctx, arg)
}
func (m *mockQuerier) ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
	return m.listCartItemsByUserIdFn(ctx, userID)
}
func (m *mockQuerier) UpdateCartItemQuantity(ctx context.Context, arg repo.UpdateCartItemQuantityParams) (repo.CartItem, error) {
	return m.updateCartItemQuantityFn(ctx, arg)
}
func (m *mockQuerier) RemoveItemFromCart(ctx context.Context, id int64) (repo.CartItem, error) {
	return m.removeItemFromCartFn(ctx, id)
}
func (m *mockQuerier) ClearCart(ctx context.Context, userID int64) error {
	return m.clearCartFn(ctx, userID)
}

// --- Querier interface を満たすためのスタブ ---

func (m *mockQuerier) AddProductToCategory(ctx context.Context, arg repo.AddProductToCategoryParams) error {
	panic("not implemented")
}
func (m *mockQuerier) CancelOrder(ctx context.Context, id int64) (repo.Order, error) {
	panic("not implemented")
}
func (m *mockQuerier) ConsumeRefreshToken(ctx context.Context, tokenHash string) (repo.RefreshToken, error) {
	panic("not implemented")
}
func (m *mockQuerier) CreateAddress(ctx context.Context, arg repo.CreateAddressParams) (repo.Address, error) {
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
func (m *mockQuerier) FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindOrderById(ctx context.Context, id int64) (repo.FindOrderByIdRow, error) {
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
func (m *mockQuerier) ListAllOrders(ctx context.Context) ([]repo.ListAllOrdersRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListCategories(ctx context.Context) ([]repo.ListCategoriesRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListProducts(ctx context.Context) ([]repo.ListProductsRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListProductsByCategory(ctx context.Context, categoryID int64) ([]repo.ListProductsByCategoryRow, error) {
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
func (m *mockQuerier) UpdateCategory(ctx context.Context, arg repo.UpdateCategoryParams) (repo.Category, error) {
	panic("not implemented")
}
func (m *mockQuerier) UpdateOrderStatus(ctx context.Context, arg repo.UpdateOrderStatusParams) (repo.Order, error) {
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
func newTestCartHelper(id int64, userID int64) repo.Cart {
	return repo.Cart{
		ID:        id,
		UserID:    userID,
		CreatedAt: pgtype.Timestamptz{Valid: true},
		UpdatedAt: pgtype.Timestamptz{Valid: true},
	}
}

func newTestCartItemHelper(id int64, cartID int64, productID int64) repo.CartItem {
	return repo.CartItem{
		ID:        id,
		CartID:    cartID,
		ProductID: productID,
		Quantity:  2,
		CreatedAt: pgtype.Timestamptz{Valid: true},
		UpdatedAt: pgtype.Timestamptz{Valid: true},
	}
}

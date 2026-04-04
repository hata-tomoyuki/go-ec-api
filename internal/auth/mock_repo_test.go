package auth

import (
	"context"
	"time"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// mockQuerier は Querier interface の手動モック（auth テスト用）
type mockQuerier struct {
	findUserByIdFn                func(ctx context.Context, id int64) (repo.User, error)
	findUserByEmailFn             func(ctx context.Context, email string) (repo.User, error)
	createUserFn                  func(ctx context.Context, arg repo.CreateUserParams) (repo.User, error)
	updateUserFn                  func(ctx context.Context, arg repo.UpdateUserParams) (repo.User, error)
	updateUserPasswordFn          func(ctx context.Context, arg repo.UpdateUserPasswordParams) (repo.User, error)
	deleteExpiredRefreshTokensFn  func(ctx context.Context) error
	insertRefreshTokenFn          func(ctx context.Context, arg repo.InsertRefreshTokenParams) (repo.RefreshToken, error)
	deleteRefreshTokenFn          func(ctx context.Context, id int64) error
	deleteRefreshTokensByUserIdFn func(ctx context.Context, userID int64) error
	consumeRefreshTokenFn         func(ctx context.Context, tokenHash string) (repo.RefreshToken, error)
	revokeTokenFn                 func(ctx context.Context, arg repo.RevokeTokenParams) error
	isTokenRevokedFn              func(ctx context.Context, jti string) (bool, error)
}

func (m *mockQuerier) FindUserById(ctx context.Context, id int64) (repo.User, error) {
	return m.findUserByIdFn(ctx, id)
}
func (m *mockQuerier) FindUserByEmail(ctx context.Context, email string) (repo.User, error) {
	return m.findUserByEmailFn(ctx, email)
}
func (m *mockQuerier) CreateUser(ctx context.Context, arg repo.CreateUserParams) (repo.User, error) {
	return m.createUserFn(ctx, arg)
}
func (m *mockQuerier) UpdateUser(ctx context.Context, arg repo.UpdateUserParams) (repo.User, error) {
	return m.updateUserFn(ctx, arg)
}
func (m *mockQuerier) UpdateUserPassword(ctx context.Context, arg repo.UpdateUserPasswordParams) (repo.User, error) {
	return m.updateUserPasswordFn(ctx, arg)
}
func (m *mockQuerier) DeleteExpiredRefreshTokens(ctx context.Context) error {
	return m.deleteExpiredRefreshTokensFn(ctx)
}
func (m *mockQuerier) InsertRefreshToken(ctx context.Context, arg repo.InsertRefreshTokenParams) (repo.RefreshToken, error) {
	return m.insertRefreshTokenFn(ctx, arg)
}
func (m *mockQuerier) DeleteRefreshToken(ctx context.Context, id int64) error {
	return m.deleteRefreshTokenFn(ctx, id)
}
func (m *mockQuerier) DeleteRefreshTokensByUserId(ctx context.Context, userID int64) error {
	return m.deleteRefreshTokensByUserIdFn(ctx, userID)
}
func (m *mockQuerier) ConsumeRefreshToken(ctx context.Context, tokenHash string) (repo.RefreshToken, error) {
	return m.consumeRefreshTokenFn(ctx, tokenHash)
}
func (m *mockQuerier) RevokeToken(ctx context.Context, arg repo.RevokeTokenParams) error {
	return m.revokeTokenFn(ctx, arg)
}
func (m *mockQuerier) IsTokenRevoked(ctx context.Context, jti string) (bool, error) {
	return m.isTokenRevokedFn(ctx, jti)
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
func (m *mockQuerier) DeleteAddress(ctx context.Context, id int32) error {
	panic("not implemented")
}
func (m *mockQuerier) DeleteCategory(ctx context.Context, id int64) (repo.Category, error) {
	panic("not implemented")
}
func (m *mockQuerier) DeleteProduct(ctx context.Context, id int64) (repo.Product, error) {
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
func (m *mockQuerier) FindOrderById(ctx context.Context, id int64) (repo.FindOrderByIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) FindProductById(ctx context.Context, id int64) (repo.FindProductByIdRow, error) {
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
func (m *mockQuerier) UpdateProduct(ctx context.Context, arg repo.UpdateProductParams) (repo.Product, error) {
	panic("not implemented")
}

// テスト用ヘルパー: ダミーの RefreshToken を生成
func newTestRefreshToken(id int64, userID int64) repo.RefreshToken {
	return repo.RefreshToken{
		ID:        id,
		UserID:    userID,
		TokenHash: "hashed-token",
		ExpiresAt: pgtype.Timestamptz{
			Time:  time.Now().Add(30 * 24 * time.Hour),
			Valid: true,
		},
	}
}

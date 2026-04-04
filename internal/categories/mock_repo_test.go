package categories

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

// mockQuerier は Querier interface の手動モック（categories テスト用）
type mockQuerier struct {
	findCategoryByIdFn          func(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error)
	listCategoriesFn            func(ctx context.Context) ([]repo.ListCategoriesRow, error)
	createCategoryFn            func(ctx context.Context, arg repo.CreateCategoryParams) (repo.Category, error)
	updateCategoryFn            func(ctx context.Context, arg repo.UpdateCategoryParams) (repo.Category, error)
	deleteCategoryFn            func(ctx context.Context, id int64) (repo.Category, error)
	listProductsByCategoryFn    func(ctx context.Context, categoryID int64) ([]repo.ListProductsByCategoryRow, error)
	addProductToCategoryFn      func(ctx context.Context, arg repo.AddProductToCategoryParams) error
	removeProductFromCategoryFn func(ctx context.Context, arg repo.RemoveProductFromCategoryParams) error
}

func (m *mockQuerier) FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
	return m.findCategoryByIdFn(ctx, id)
}
func (m *mockQuerier) ListCategories(ctx context.Context) ([]repo.ListCategoriesRow, error) {
	return m.listCategoriesFn(ctx)
}
func (m *mockQuerier) CreateCategory(ctx context.Context, arg repo.CreateCategoryParams) (repo.Category, error) {
	return m.createCategoryFn(ctx, arg)
}
func (m *mockQuerier) UpdateCategory(ctx context.Context, arg repo.UpdateCategoryParams) (repo.Category, error) {
	return m.updateCategoryFn(ctx, arg)
}
func (m *mockQuerier) DeleteCategory(ctx context.Context, id int64) (repo.Category, error) {
	return m.deleteCategoryFn(ctx, id)
}
func (m *mockQuerier) ListProductsByCategory(ctx context.Context, categoryID int64) ([]repo.ListProductsByCategoryRow, error) {
	return m.listProductsByCategoryFn(ctx, categoryID)
}
func (m *mockQuerier) AddProductToCategory(ctx context.Context, arg repo.AddProductToCategoryParams) error {
	return m.addProductToCategoryFn(ctx, arg)
}
func (m *mockQuerier) RemoveProductFromCategory(ctx context.Context, arg repo.RemoveProductFromCategoryParams) error {
	return m.removeProductFromCategoryFn(ctx, arg)
}

// --- 以下は Querier interface を満たすためのスタブ ---

func (m *mockQuerier) AddItemToCart(ctx context.Context, arg repo.AddItemToCartParams) (repo.CartItem, error) {
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
func (m *mockQuerier) FindOrderById(ctx context.Context, id int64) ([]repo.FindOrderByIdRow, error) {
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
func (m *mockQuerier) ListCartItemsByUserId(ctx context.Context, userID int64) ([]repo.ListCartItemsByUserIdRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListOrdersByCustomerID(ctx context.Context, customerID int64) ([]repo.ListOrdersByCustomerIDRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) ListProducts(ctx context.Context) ([]repo.ListProductsRow, error) {
	panic("not implemented")
}
func (m *mockQuerier) RemoveItemFromCart(ctx context.Context, id int64) (repo.CartItem, error) {
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

// テスト用のヘルパー: ダミーの Category を生成
func newTestCategory(id int64, name string) repo.Category {
	return repo.Category{
		ID:          id,
		Name:        name,
		Description: pgtype.Text{Valid: false},
		CreatedAt:   pgtype.Timestamptz{Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Valid: true},
		ImageColor:  "from-gray-400 to-gray-600",
	}
}

func newTestListCategoriesRow(id int64, name string) repo.ListCategoriesRow {
	return repo.ListCategoriesRow{
		ID:          id,
		Name:        name,
		Description: pgtype.Text{Valid: false},
		CreatedAt:   pgtype.Timestamptz{Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Valid: true},
		ImageColor:  "from-gray-400 to-gray-600",
	}
}

func newTestFindCategoryByIdRow(id int64, name string) repo.FindCategoryByIdRow {
	return repo.FindCategoryByIdRow{
		ID:          id,
		Name:        name,
		Description: pgtype.Text{Valid: false},
		CreatedAt:   pgtype.Timestamptz{Valid: true},
		UpdatedAt:   pgtype.Timestamptz{Valid: true},
		ImageColor:  "from-gray-400 to-gray-600",
	}
}

package products

import (
	"context"
	"errors"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

// ---------- validate() テスト ----------

func TestListProductsParams_Validate_Defaults(t *testing.T) {
	p := listProductsParams{}
	if err := p.validate(); err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if p.Page != 1 {
		t.Errorf("expected Page=1, got %d", p.Page)
	}
	if p.Limit != 20 {
		t.Errorf("expected Limit=20, got %d", p.Limit)
	}
	if p.Sort != "created_at_desc" {
		t.Errorf("expected Sort='created_at_desc', got '%s'", p.Sort)
	}
}

func TestListProductsParams_Validate_InvalidSort(t *testing.T) {
	p := listProductsParams{Sort: "invalid_sort"}
	err := p.validate()
	if !errors.Is(err, ErrInvalidSort) {
		t.Errorf("expected ErrInvalidSort, got %v", err)
	}
}

func TestListProductsParams_Validate_LimitTooHigh(t *testing.T) {
	p := listProductsParams{Limit: 101}
	err := p.validate()
	if !errors.Is(err, ErrInvalidLimit) {
		t.Errorf("expected ErrInvalidLimit, got %v", err)
	}
}

func TestListProductsParams_Validate_LimitBoundary(t *testing.T) {
	p := listProductsParams{Limit: 100}
	if err := p.validate(); err != nil {
		t.Fatalf("limit=100 should be valid, got: %v", err)
	}
}

func TestListProductsParams_Validate_AllSortOptions(t *testing.T) {
	validSorts := []string{"created_at_desc", "created_at_asc", "price_desc", "price_asc", "name_asc", "name_desc"}
	for _, sort := range validSorts {
		p := listProductsParams{Sort: sort}
		if err := p.validate(); err != nil {
			t.Errorf("sort=%q should be valid, got: %v", sort, err)
		}
	}
}

func TestListProducts(t *testing.T) {
	mock := &mockQuerier{
		listProductsFn: func(ctx context.Context) ([]repo.ListProductsRow, error) {
			return []repo.ListProductsRow{
				newTestListProductsRow(1, "T-shirt", 2000),
				newTestListProductsRow(2, "Hoodie", 5000),
			}, nil
		},
	}

	svc := NewService(mock, &mockDBTX{})
	products, err := svc.ListProducts(context.Background())

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(products) != 2 {
		t.Fatalf("expected 2 products, got %d", len(products))
	}
	if products[0].Name != "T-shirt" {
		t.Errorf("expected name 'T-shirt', got '%s'", products[0].Name)
	}
}

func TestListProducts_Error(t *testing.T) {
	mock := &mockQuerier{
		listProductsFn: func(ctx context.Context) ([]repo.ListProductsRow, error) {
			return nil, errors.New("db connection failed")
		},
	}

	svc := NewService(mock, &mockDBTX{})
	_, err := svc.ListProducts(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFindProductById(t *testing.T) {
	mock := &mockQuerier{
		findProductByIdFn: func(ctx context.Context, id int64) (repo.FindProductByIdRow, error) {
			return newTestFindProductByIdRow(1, "T-shirt", 2000), nil
		},
	}

	svc := NewService(mock, &mockDBTX{})
	product, err := svc.FindProductById(context.Background(), 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.ID != 1 {
		t.Errorf("expected id=1, got %d", product.ID)
	}
	if product.PriceInCents != 2000 {
		t.Errorf("expected price=2000, got %d", product.PriceInCents)
	}
}

func TestCreateProduct(t *testing.T) {
	mock := &mockQuerier{
		createProductFn: func(ctx context.Context, arg repo.CreateProductParams) (repo.Product, error) {
			return repo.Product{
				ID:           1,
				Name:         arg.Name,
				PriceInCents: arg.PriceInCents,
			}, nil
		},
	}

	svc := NewService(mock, &mockDBTX{})
	product, err := svc.CreateProduct(context.Background(), createProductParams{
		Name:         "New Jacket",
		PriceInCents: 8000,
		Quantity:     5,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.Name != "New Jacket" {
		t.Errorf("expected name 'New Jacket', got '%s'", product.Name)
	}
	if product.PriceInCents != 8000 {
		t.Errorf("expected price=8000, got %d", product.PriceInCents)
	}
}

func TestUpdateProduct(t *testing.T) {
	mock := &mockQuerier{
		updateProductFn: func(ctx context.Context, arg repo.UpdateProductParams) (repo.Product, error) {
			return repo.Product{
				ID:           arg.ID,
				Name:         arg.Name,
				PriceInCents: arg.PriceInCents,
			}, nil
		},
	}

	svc := NewService(mock, &mockDBTX{})
	product, err := svc.UpdateProduct(context.Background(), updateProductParams{
		ID:           1,
		Name:         "Updated Jacket",
		PriceInCents: 9000,
		Quantity:     3,
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if product.Name != "Updated Jacket" {
		t.Errorf("expected name 'Updated Jacket', got '%s'", product.Name)
	}
}

func TestDeleteProduct(t *testing.T) {
	mock := &mockQuerier{
		deleteProductFn: func(ctx context.Context, id int64) (repo.Product, error) {
			return repo.Product{}, nil
		},
	}

	svc := NewService(mock, &mockDBTX{})
	err := svc.DeleteProduct(context.Background(), 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteProduct_Error(t *testing.T) {
	mock := &mockQuerier{
		deleteProductFn: func(ctx context.Context, id int64) (repo.Product, error) {
			return repo.Product{}, errors.New("product not found")
		},
	}

	svc := NewService(mock, &mockDBTX{})
	err := svc.DeleteProduct(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

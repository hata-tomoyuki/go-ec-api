package products

import (
	"context"
	"errors"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

func TestListProducts(t *testing.T) {
	mock := &mockQuerier{
		listProductsFn: func(ctx context.Context) ([]repo.Product, error) {
			return []repo.Product{
				newTestProduct(1, "T-shirt", 2000),
				newTestProduct(2, "Hoodie", 5000),
			}, nil
		},
	}

	svc := NewService(mock)
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
		listProductsFn: func(ctx context.Context) ([]repo.Product, error) {
			return nil, errors.New("db connection failed")
		},
	}

	svc := NewService(mock)
	_, err := svc.ListProducts(context.Background())

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

func TestFindProductById(t *testing.T) {
	mock := &mockQuerier{
		findProductByIdFn: func(ctx context.Context, id int64) (repo.Product, error) {
			return newTestProduct(1, "T-shirt", 2000), nil
		},
	}

	svc := NewService(mock)
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

	svc := NewService(mock)
	product, err := svc.CreateProduct(context.Background(), createProductParams{
		Name:         "New Jacket",
		PriceInCents: 8000,
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

	svc := NewService(mock)
	product, err := svc.UpdateProduct(context.Background(), updateProductParams{
		ID:           1,
		Name:         "Updated Jacket",
		PriceInCents: 9000,
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

	svc := NewService(mock)
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

	svc := NewService(mock)
	err := svc.DeleteProduct(context.Background(), 999)

	if err == nil {
		t.Fatal("expected error, got nil")
	}
}

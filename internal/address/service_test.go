package address

import (
	"context"
	"errors"
	"testing"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

func TestListAddressesByUserId(t *testing.T) {
	mock := &mockQuerier{
		listAddressesByUserIdFn: func(ctx context.Context, userID int64) ([]repo.Address, error) {
			return []repo.Address{
				newTestAddress(1, userID),
				newTestAddress(2, userID),
			}, nil
		},
	}

	svc := NewService(mock)
	addresses, err := svc.ListAddressesByUserId(context.Background(), 10)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(addresses) != 2 {
		t.Fatalf("expected 2 addresses, got %d", len(addresses))
	}
	if addresses[0].UserID != 10 {
		t.Errorf("expected user_id=10, got %d", addresses[0].UserID)
	}
}

func TestFindAddressById_Success(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			return newTestAddress(1, 10), nil
		},
	}

	svc := NewService(mock)
	addr, err := svc.FindAddressById(context.Background(), 10, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.ID != 1 {
		t.Errorf("expected id=1, got %d", addr.ID)
	}
}

func TestFindAddressById_NotFound(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			return repo.Address{}, pgx.ErrNoRows
		},
	}

	svc := NewService(mock)
	_, err := svc.FindAddressById(context.Background(), 10, 999)

	if !errors.Is(err, ErrAddressNotFound) {
		t.Fatalf("expected ErrAddressNotFound, got %v", err)
	}
}

func TestFindAddressById_Forbidden(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			// 他人(user_id=99)の住所を返す
			return newTestAddress(1, 99), nil
		},
	}

	svc := NewService(mock)
	_, err := svc.FindAddressById(context.Background(), 10, 1)

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestCreateAddress(t *testing.T) {
	mock := &mockQuerier{
		createAddressFn: func(ctx context.Context, arg repo.CreateAddressParams) (repo.Address, error) {
			return repo.Address{
				ID:      1,
				UserID:  arg.UserID,
				Street:  arg.Street,
				City:    arg.City,
				State:   arg.State,
				ZipCode: arg.ZipCode,
				Country: arg.Country,
			}, nil
		},
	}

	svc := NewService(mock)
	addr, err := svc.CreateAddress(context.Background(), 10, createAddressParams{
		Street:  "1-1 Shibuya",
		City:    "Shibuya",
		State:   "Tokyo",
		ZipCode: "150-0001",
		Country: "Japan",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.Street != "1-1 Shibuya" {
		t.Errorf("expected street '1-1 Shibuya', got '%s'", addr.Street)
	}
	if addr.UserID != 10 {
		t.Errorf("expected user_id=10, got %d", addr.UserID)
	}
}

func TestUpdateAddress_Success(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			return newTestAddress(1, 10), nil // 自分の住所
		},
		updateAddressFn: func(ctx context.Context, arg repo.UpdateAddressParams) (repo.Address, error) {
			return repo.Address{
				ID:      arg.ID,
				UserID:  10,
				Street:  arg.Street,
				City:    arg.City,
				State:   arg.State,
				ZipCode: arg.ZipCode,
				Country: arg.Country,
			}, nil
		},
	}

	svc := NewService(mock)
	addr, err := svc.UpdateAddress(context.Background(), 10, 1, createAddressParams{
		Street:  "2-2 Minato",
		City:    "Minato",
		State:   "Tokyo",
		ZipCode: "105-0001",
		Country: "Japan",
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if addr.Street != "2-2 Minato" {
		t.Errorf("expected street '2-2 Minato', got '%s'", addr.Street)
	}
}

func TestUpdateAddress_NotFound(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			return repo.Address{}, pgx.ErrNoRows
		},
	}

	svc := NewService(mock)
	_, err := svc.UpdateAddress(context.Background(), 10, 999, createAddressParams{
		Street: "x", City: "x", State: "x", ZipCode: "x", Country: "x",
	})

	if !errors.Is(err, ErrAddressNotFound) {
		t.Fatalf("expected ErrAddressNotFound, got %v", err)
	}
}

func TestUpdateAddress_Forbidden(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			return newTestAddress(1, 99), nil // 他人の住所
		},
	}

	svc := NewService(mock)
	_, err := svc.UpdateAddress(context.Background(), 10, 1, createAddressParams{
		Street: "x", City: "x", State: "x", ZipCode: "x", Country: "x",
	})

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

func TestDeleteAddress_Success(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			return newTestAddress(1, 10), nil
		},
		deleteAddressFn: func(ctx context.Context, id int32) error {
			return nil
		},
	}

	svc := NewService(mock)
	err := svc.DeleteAddress(context.Background(), 10, 1)

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestDeleteAddress_NotFound(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			return repo.Address{}, pgx.ErrNoRows
		},
	}

	svc := NewService(mock)
	err := svc.DeleteAddress(context.Background(), 10, 999)

	if !errors.Is(err, ErrAddressNotFound) {
		t.Fatalf("expected ErrAddressNotFound, got %v", err)
	}
}

func TestDeleteAddress_Forbidden(t *testing.T) {
	mock := &mockQuerier{
		findAddressByIdFn: func(ctx context.Context, id int32) (repo.Address, error) {
			return newTestAddress(1, 99), nil // 他人の住所
		},
	}

	svc := NewService(mock)
	err := svc.DeleteAddress(context.Background(), 10, 1)

	if !errors.Is(err, ErrForbidden) {
		t.Fatalf("expected ErrForbidden, got %v", err)
	}
}

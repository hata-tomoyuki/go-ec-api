package address

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) ListAddressesByUserId(ctx context.Context, userId int64) ([]repo.Address, error) {
	return s.repo.ListAddressesByUserId(ctx, userId)
}

func (s *svc) FindAddressById(ctx context.Context, userId int64, addressId int32) (repo.Address, error) {
	addr, err := s.repo.FindAddressById(ctx, addressId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Address{}, ErrAddressNotFound
		}
		return repo.Address{}, err
	}

	if addr.UserID != userId {
		return repo.Address{}, ErrForbidden
	}

	return addr, nil
}

func (s *svc) CreateAddress(ctx context.Context, userId int64, params createAddressParams) (repo.Address, error) {
	return s.repo.CreateAddress(ctx, repo.CreateAddressParams{
		UserID:  userId,
		Street:  params.Street,
		City:    params.City,
		State:   params.State,
		ZipCode: params.ZipCode,
		Country: params.Country,
	})
}

func (s *svc) UpdateAddress(ctx context.Context, userId int64, addressId int32, params createAddressParams) (repo.Address, error) {
	// 所有者チェック
	addr, err := s.repo.FindAddressById(ctx, addressId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Address{}, ErrAddressNotFound
		}
		return repo.Address{}, err
	}

	if addr.UserID != userId {
		return repo.Address{}, ErrForbidden
	}

	return s.repo.UpdateAddress(ctx, repo.UpdateAddressParams{
		ID:      addressId,
		Street:  params.Street,
		City:    params.City,
		State:   params.State,
		ZipCode: params.ZipCode,
		Country: params.Country,
	})
}

func (s *svc) DeleteAddress(ctx context.Context, userId int64, addressId int32) error {
	// 所有者チェック
	addr, err := s.repo.FindAddressById(ctx, addressId)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrAddressNotFound
		}
		return err
	}

	if addr.UserID != userId {
		return ErrForbidden
	}

	return s.repo.DeleteAddress(ctx, addressId)
}

package address

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type svc struct {
	repo repo.Querier
}

func NewService(repo repo.Querier) Service {
	return &svc{repo: repo}
}

func (s *svc) FindAddressByUserId(ctx context.Context, userId int64) (repo.Address, error) {
	return s.repo.FindAddressByUserId(ctx, userId)
}

func (s *svc) CreateAddress(ctx context.Context, userId int64, tempAddress createAddressParams) (repo.Address, error) {
	return s.repo.CreateAddress(ctx, repo.CreateAddressParams{
		UserID:  userId,
		Street:  tempAddress.Street,
		City:    tempAddress.City,
		State:   tempAddress.State,
		ZipCode: tempAddress.ZipCode,
		Country: tempAddress.Country,
	})
}

func (s *svc) UpdateAddress(ctx context.Context, userId int64, tempAddress createAddressParams) (repo.Address, error) {
	return s.repo.UpdateAddress(ctx, repo.UpdateAddressParams{
		UserID:  userId,
		Street:  tempAddress.Street,
		City:    tempAddress.City,
		State:   tempAddress.State,
		ZipCode: tempAddress.ZipCode,
		Country: tempAddress.Country,
	})
}

func (s *svc) DeleteAddress(ctx context.Context, userId int64) (repo.Address, error) {
	return s.repo.DeleteAddress(ctx, userId)
}

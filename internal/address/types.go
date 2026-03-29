package address

import (
	"context"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

type createAddressParams struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

type Service interface {
	FindAddressByUserId(ctx context.Context, userId int64) (repo.Address, error)
	CreateAddress(ctx context.Context, userId int64, tempAddress createAddressParams) (repo.Address, error)
	UpdateAddress(ctx context.Context, userId int64, tempAddress createAddressParams) (repo.Address, error)
}

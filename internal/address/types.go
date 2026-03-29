package address

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
)

var (
	ErrForbidden        = errors.New("you do not have permission to access this address")
	ErrAddressNotFound  = errors.New("address not found")
	ErrValidation       = errors.New("all address fields are required")
)

type createAddressParams struct {
	Street  string `json:"street"`
	City    string `json:"city"`
	State   string `json:"state"`
	ZipCode string `json:"zip_code"`
	Country string `json:"country"`
}

func (p createAddressParams) validate() error {
	if p.Street == "" || p.City == "" || p.State == "" || p.ZipCode == "" || p.Country == "" {
		return ErrValidation
	}
	return nil
}

type Service interface {
	ListAddressesByUserId(ctx context.Context, userId int64) ([]repo.Address, error)
	FindAddressById(ctx context.Context, userId int64, addressId int32) (repo.Address, error)
	CreateAddress(ctx context.Context, userId int64, params createAddressParams) (repo.Address, error)
	UpdateAddress(ctx context.Context, userId int64, addressId int32, params createAddressParams) (repo.Address, error)
	DeleteAddress(ctx context.Context, userId int64, addressId int32) error
}

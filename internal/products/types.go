package products

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5/pgtype"
)

var (
	ErrProductNotFound   = errors.New("product not found")
	ErrProductValidation = errors.New("name is required and price_in_cents must be greater than 0")
	ErrQuantityNegative  = errors.New("quantity must not be negative")
	ErrInvalidSort       = errors.New("invalid sort parameter")
	ErrInvalidLimit      = errors.New("limit must be between 1 and 100")
)

var allowedSorts = map[string]string{
	"created_at_desc": "p.created_at DESC",
	"created_at_asc":  "p.created_at ASC",
	"price_desc":      "p.price_in_cents DESC",
	"price_asc":       "p.price_in_cents ASC",
	"name_asc":        "p.name ASC",
	"name_desc":       "p.name DESC",
}

type listProductsParams struct {
	Page   int
	Limit  int
	Sort   string
	Search string
}

func (p *listProductsParams) validate() error {
	if p.Page <= 0 {
		p.Page = 1
	}
	if p.Limit <= 0 {
		p.Limit = 20
	}
	if p.Limit > 100 {
		return ErrInvalidLimit
	}
	if p.Sort == "" {
		p.Sort = "created_at_desc"
	}
	if _, ok := allowedSorts[p.Sort]; !ok {
		return ErrInvalidSort
	}
	return nil
}

type paginatedProducts struct {
	Data  []paginatedProductRow `json:"data"`
	Total int                   `json:"total"`
	Page  int                   `json:"page"`
	Limit int                   `json:"limit"`
}

type paginatedProductRow struct {
	ID           int64              `json:"id"`
	Name         string             `json:"name"`
	Description  string             `json:"description"`
	PriceInCents int32              `json:"price_in_cents"`
	Quantity     int32              `json:"quantity"`
	ImageColor   string             `json:"image_color"`
	CreatedAt    pgtype.Timestamptz `json:"created_at"`
	CategoryID   int64              `json:"category_id"`
	CategoryName string             `json:"category_name"`
}

type createProductParams struct {
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
	Description  string `json:"description"`
	ImageColor   string `json:"image_color"`
	Quantity     int32  `json:"quantity"`
}

func (p createProductParams) validate() error {
	if p.Name == "" || p.PriceInCents <= 0 {
		return ErrProductValidation
	}
	if p.Quantity < 0 {
		return ErrQuantityNegative
	}
	return nil
}

type updateProductParams struct {
	ID           int64  `json:"id"`
	Name         string `json:"name"`
	PriceInCents int32  `json:"price_in_cents"`
	Description  string `json:"description"`
	ImageColor   string `json:"image_color"`
	Quantity     int32  `json:"quantity"`
}

func (p updateProductParams) validate() error {
	if p.Name == "" || p.PriceInCents <= 0 {
		return ErrProductValidation
	}
	if p.Quantity < 0 {
		return ErrQuantityNegative
	}
	return nil
}

type Service interface {
	ListProducts(ctx context.Context) ([]repo.ListProductsRow, error)
	ListProductsPaginated(ctx context.Context, params listProductsParams) (paginatedProducts, error)
	FindProductById(ctx context.Context, id int64) (repo.FindProductByIdRow, error)
	CreateProduct(ctx context.Context, tempProduct createProductParams) (repo.Product, error)
	UpdateProduct(ctx context.Context, tempProduct updateProductParams) (repo.Product, error)
	DeleteProduct(ctx context.Context, id int64) error
}

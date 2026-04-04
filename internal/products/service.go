package products

import (
	"context"
	"errors"
	"fmt"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
)

type svc struct {
	repo repo.Querier
	db   repo.DBTX
}

func NewService(repo repo.Querier, db repo.DBTX) Service {
	return &svc{repo: repo, db: db}
}

func (s *svc) ListProducts(ctx context.Context) ([]repo.ListProductsRow, error) {
	return s.repo.ListProducts(ctx)
}

func buildWhereClause(params listProductsParams) (string, []interface{}, int) {
	whereSQL := ""
	args := []interface{}{}
	argN := 1

	if params.Search != "" {
		whereSQL = fmt.Sprintf(" WHERE p.name ILIKE $%d", argN)
		args = append(args, "%"+params.Search+"%")
		argN++
	}

	if params.CategoryID > 0 {
		if whereSQL == "" {
			whereSQL = fmt.Sprintf(" WHERE pc.category_id = $%d", argN)
		} else {
			whereSQL += fmt.Sprintf(" AND pc.category_id = $%d", argN)
		}
		args = append(args, params.CategoryID)
		argN++
	}

	return whereSQL, args, argN
}

func (s *svc) ListProductsPaginated(ctx context.Context, params listProductsParams) (paginatedProducts, error) {
	whereSQL, whereArgs, argN := buildWhereClause(params)

	sql := `
	SELECT p.id, p.name, p.description, p.price_in_cents, p.quantity,
	       p.image_color, p.created_at,
		   COALESCE(pc.category_id, 0)::bigint AS category_id,
		   COALESCE(c.name, '')::text AS category_name
	FROM products p
	LEFT JOIN product_categories pc ON p.id = pc.product_id
	LEFT JOIN categories c ON pc.category_id = c.id
	` + whereSQL

	sortSQL, _ := allowedSorts[params.Sort]
	offset := (params.Page - 1) * params.Limit

	sql += fmt.Sprintf(" ORDER BY %s LIMIT $%d OFFSET $%d", sortSQL, argN, argN+1)
	args := append(whereArgs, params.Limit, offset)

	rows, err := s.db.Query(ctx, sql, args...)
	if err != nil {
		return paginatedProducts{}, err
	}
	defer rows.Close()

	products := make([]paginatedProductRow, 0)
	for rows.Next() {
		var p paginatedProductRow
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PriceInCents, &p.Quantity, &p.ImageColor, &p.CreatedAt, &p.CategoryID, &p.CategoryName); err != nil {
			return paginatedProducts{}, err
		}
		products = append(products, p)
	}

	countSQL := "SELECT COUNT(*) FROM products p"
	if params.CategoryID > 0 {
		countSQL += " LEFT JOIN product_categories pc ON p.id = pc.product_id"
	}
	countSQL += whereSQL
	var total int
	if err := s.db.QueryRow(ctx, countSQL, whereArgs...).Scan(&total); err != nil {
		return paginatedProducts{}, err
	}

	return paginatedProducts{
		Data:  products,
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
	}, nil
}

func (s *svc) FindProductById(ctx context.Context, id int64) (repo.FindProductByIdRow, error) {
	product, err := s.repo.FindProductById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.FindProductByIdRow{}, ErrProductNotFound
		}
		return repo.FindProductByIdRow{}, err
	}
	return product, nil
}

func (s *svc) CreateProduct(ctx context.Context, tempProduct createProductParams) (repo.Product, error) {
	return s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         tempProduct.Name,
		PriceInCents: tempProduct.PriceInCents,
		Description:  tempProduct.Description,
		ImageColor:   tempProduct.ImageColor,
		Quantity:     tempProduct.Quantity,
	})
}

func (s *svc) UpdateProduct(ctx context.Context, tempProduct updateProductParams) (repo.Product, error) {
	product, err := s.repo.UpdateProduct(ctx, repo.UpdateProductParams{
		ID:           tempProduct.ID,
		Name:         tempProduct.Name,
		PriceInCents: tempProduct.PriceInCents,
		Description:  tempProduct.Description,
		ImageColor:   tempProduct.ImageColor,
		Quantity:     tempProduct.Quantity,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Product{}, ErrProductNotFound
		}
		return repo.Product{}, err
	}
	return product, nil
}

func (s *svc) DeleteProduct(ctx context.Context, id int64) error {
	_, err := s.repo.DeleteProduct(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrProductNotFound
		}
		return err
	}
	return nil
}

package products

import (
	"context"
	"errors"
	"fmt"
	"log/slog"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"example.com/ecommerce/internal/cache"
	"github.com/jackc/pgx/v5"
)

type svc struct {
	repo  repo.Querier
	db    repo.DBTX
	cache *cache.Store
}

func NewService(repo repo.Querier, db repo.DBTX, cacheStore ...*cache.Store) Service {
	s := &svc{repo: repo, db: db}
	if len(cacheStore) > 0 {
		s.cache = cacheStore[0]
	}
	return s
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
	// Cache check
	key := cache.ProductListKey(params.Page, params.Limit, params.Sort, params.Search, params.CategoryID)
	if key != "" {
		if cached, ok := cache.Get[paginatedProducts](ctx, s.cache, key); ok {
			slog.Debug("products list cache hit", "key", key)
			return cached, nil
		}
	}

	whereSQL, whereArgs, argN := buildWhereClause(params)

	sql := `
	SELECT p.id, p.name, p.description, p.price_in_cents, p.quantity,
	       p.image_color, p.created_at, COUNT(*) OVER() AS total_count,
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
		if err := rows.Scan(&p.ID, &p.Name, &p.Description, &p.PriceInCents, &p.Quantity, &p.ImageColor, &p.CreatedAt, &p.TotalCount, &p.CategoryID, &p.CategoryName); err != nil {
			return paginatedProducts{}, err
		}
		products = append(products, p)
	}

	total := 0
	if len(products) > 0 {
		total = products[0].TotalCount
	}

	result := paginatedProducts{
		Data:  products,
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
	}

	// Cache set
	if key != "" {
		cache.Set(ctx, s.cache, key, result, cache.ProductListTTL)
	}

	return result, nil
}

func (s *svc) FindProductById(ctx context.Context, id int64) (repo.FindProductByIdRow, error) {
	// Cache check
	key := cache.ProductDetailKey(id)
	if cached, ok := cache.Get[repo.FindProductByIdRow](ctx, s.cache, key); ok {
		slog.Debug("product detail cache hit", "id", id)
		return cached, nil
	}

	product, err := s.repo.FindProductById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.FindProductByIdRow{}, ErrProductNotFound
		}
		return repo.FindProductByIdRow{}, err
	}

	// Cache set
	cache.Set(ctx, s.cache, key, product, cache.ProductDetailTTL)

	return product, nil
}

func (s *svc) CreateProduct(ctx context.Context, tempProduct createProductParams) (repo.Product, error) {
	product, err := s.repo.CreateProduct(ctx, repo.CreateProductParams{
		Name:         tempProduct.Name,
		PriceInCents: tempProduct.PriceInCents,
		Description:  tempProduct.Description,
		ImageColor:   tempProduct.ImageColor,
		Quantity:     tempProduct.Quantity,
	})
	if err != nil {
		return repo.Product{}, err
	}

	// Invalidate list caches
	s.cache.InvalidateByPrefix(ctx, cache.ProductListPrefix)

	return product, nil
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

	// Invalidate detail + list + category products caches
	s.cache.Delete(ctx, cache.ProductDetailKey(tempProduct.ID))
	s.cache.InvalidateByPrefix(ctx, cache.ProductListPrefix, cache.CategoryProductsPrefix)

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

	// Invalidate detail + list + category products + category list caches
	s.cache.Delete(ctx, cache.ProductDetailKey(id))
	s.cache.InvalidateByPrefix(ctx, cache.ProductListPrefix, cache.CategoryProductsPrefix, cache.CategoryListPrefix)

	return nil
}

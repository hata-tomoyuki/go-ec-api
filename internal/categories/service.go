package categories

import (
	"context"
	"errors"
	"log/slog"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"example.com/ecommerce/internal/cache"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
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

func (s *svc) ListCategories(ctx context.Context) ([]repo.ListCategoriesRow, error) {
	return s.repo.ListCategories(ctx)
}

func (s *svc) ListCategoriesPaginated(ctx context.Context, params listCategoriesParams) (paginatedCategories, error) {
	// Cache check
	key := cache.CategoryListKey(params.Page, params.Limit)
	if cached, ok := cache.Get[paginatedCategories](ctx, s.cache, key); ok {
		slog.Debug("categories list cache hit", "key", key)
		return cached, nil
	}

	offset := (params.Page - 1) * params.Limit

	sql := `
	SELECT c.id, c.name, COALESCE(c.description, '')::text AS description,
	       c.image_color,
	       (SELECT COUNT(*) FROM product_categories pc WHERE pc.category_id = c.id)::bigint AS product_count,
	       COUNT(*) OVER() AS total_count
	FROM categories c
	ORDER BY c.created_at DESC
	LIMIT $1 OFFSET $2`

	rows, err := s.db.Query(ctx, sql, params.Limit, offset)
	if err != nil {
		return paginatedCategories{}, err
	}
	defer rows.Close()

	cats := make([]paginatedCategoryRow, 0)
	for rows.Next() {
		var c paginatedCategoryRow
		if err := rows.Scan(&c.ID, &c.Name, &c.Description, &c.ImageColor, &c.ProductCount, &c.TotalCount); err != nil {
			return paginatedCategories{}, err
		}
		cats = append(cats, c)
	}

	total := 0
	if len(cats) > 0 {
		total = cats[0].TotalCount
	}

	result := paginatedCategories{
		Data:  cats,
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
	}

	// Cache set
	cache.Set(ctx, s.cache, key, result, cache.CategoryListTTL)

	return result, nil
}

func (s *svc) CreateCategories(ctx context.Context, name string, description *string, imageColor string) (repo.Category, error) {
	desc := pgtype.Text{Valid: false}
	if description != nil {
		desc = pgtype.Text{
			String: *description,
			Valid:  true,
		}
	}

	category, err := s.repo.CreateCategory(ctx, repo.CreateCategoryParams{
		Name:        name,
		Description: desc,
		ImageColor:  imageColor,
	})
	if err != nil {
		return repo.Category{}, err
	}

	// Invalidate list caches
	s.cache.InvalidateByPrefix(ctx, cache.CategoryListPrefix)

	return category, nil
}

func (s *svc) FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
	// Cache check
	key := cache.CategoryDetailKey(id)
	if cached, ok := cache.Get[repo.FindCategoryByIdRow](ctx, s.cache, key); ok {
		slog.Debug("category detail cache hit", "id", id)
		return cached, nil
	}

	category, err := s.repo.FindCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.FindCategoryByIdRow{}, ErrCategoryNotFound
		}
		return repo.FindCategoryByIdRow{}, err
	}

	// Cache set
	cache.Set(ctx, s.cache, key, category, cache.CategoryDetailTTL)

	return category, nil
}

func (s *svc) UpdateCategories(ctx context.Context, id int64, name string, description *string, imageColor string) (repo.Category, error) {
	desc := pgtype.Text{Valid: false}
	if description != nil {
		desc = pgtype.Text{
			String: *description,
			Valid:  true,
		}
	}

	category, err := s.repo.UpdateCategory(ctx, repo.UpdateCategoryParams{
		ID:          id,
		Name:        name,
		Description: desc,
		ImageColor:  imageColor,
	})
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.Category{}, ErrCategoryNotFound
		}
		return repo.Category{}, err
	}

	// Invalidate detail + list caches
	s.cache.Delete(ctx, cache.CategoryDetailKey(id))
	s.cache.InvalidateByPrefix(ctx, cache.CategoryListPrefix)

	return category, nil
}

func (s *svc) DeleteCategory(ctx context.Context, id int64) error {
	_, err := s.repo.DeleteCategory(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return ErrCategoryNotFound
		}
		return err
	}

	// Invalidate detail + list + category products caches
	s.cache.Delete(ctx, cache.CategoryDetailKey(id))
	s.cache.InvalidateByPrefix(ctx, cache.CategoryListPrefix, cache.CategoryProductsPrefix)

	return nil
}

func (s *svc) ListProductsByCategory(ctx context.Context, categoryId int64) ([]repo.ListProductsByCategoryRow, error) {
	// Cache check
	key := cache.CategoryProductsKey(categoryId)
	if cached, ok := cache.Get[[]repo.ListProductsByCategoryRow](ctx, s.cache, key); ok {
		slog.Debug("category products cache hit", "categoryId", categoryId)
		return cached, nil
	}

	products, err := s.repo.ListProductsByCategory(ctx, categoryId)
	if err != nil {
		return nil, err
	}

	// Cache set
	cache.Set(ctx, s.cache, key, products, cache.CategoryProductsTTL)

	return products, nil
}

func (s *svc) AddProductToCategory(ctx context.Context, categoryId int64, productId int64) error {
	err := s.repo.AddProductToCategory(ctx, repo.AddProductToCategoryParams{
		ProductID:  productId,
		CategoryID: categoryId,
	})
	if err != nil {
		return err
	}

	// Invalidate related caches
	s.cache.Delete(ctx, cache.CategoryProductsKey(categoryId))
	s.cache.Delete(ctx, cache.CategoryDetailKey(categoryId))
	s.cache.InvalidateByPrefix(ctx, cache.CategoryListPrefix, cache.ProductListPrefix)

	return nil
}

func (s *svc) RemoveProductFromCategory(ctx context.Context, categoryId int64, productId int64) error {
	err := s.repo.RemoveProductFromCategory(ctx, repo.RemoveProductFromCategoryParams{
		ProductID:  productId,
		CategoryID: categoryId,
	})
	if err != nil {
		return err
	}

	// Invalidate related caches
	s.cache.Delete(ctx, cache.CategoryProductsKey(categoryId))
	s.cache.Delete(ctx, cache.CategoryDetailKey(categoryId))
	s.cache.InvalidateByPrefix(ctx, cache.CategoryListPrefix, cache.ProductListPrefix)

	return nil
}

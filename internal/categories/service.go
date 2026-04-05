package categories

import (
	"context"
	"errors"

	repo "example.com/ecommerce/internal/adapters/postgresql/sqlc"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgtype"
)

type svc struct {
	repo repo.Querier
	db   repo.DBTX
}

func NewService(repo repo.Querier, db ...repo.DBTX) Service {
	s := &svc{repo: repo}
	if len(db) > 0 {
		s.db = db[0]
	}
	return s
}

func (s *svc) ListCategories(ctx context.Context) ([]repo.ListCategoriesRow, error) {
	return s.repo.ListCategories(ctx)
}

func (s *svc) ListCategoriesPaginated(ctx context.Context, params listCategoriesParams) (paginatedCategories, error) {
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

	return paginatedCategories{
		Data:  cats,
		Total: total,
		Page:  params.Page,
		Limit: params.Limit,
	}, nil
}

func (s *svc) CreateCategories(ctx context.Context, name string, description *string, imageColor string) (repo.Category, error) {
	desc := pgtype.Text{Valid: false}
	if description != nil {
		desc = pgtype.Text{
			String: *description,
			Valid:  true,
		}
	}

	return s.repo.CreateCategory(ctx, repo.CreateCategoryParams{
		Name:        name,
		Description: desc,
		ImageColor:  imageColor,
	})
}

func (s *svc) FindCategoryById(ctx context.Context, id int64) (repo.FindCategoryByIdRow, error) {
	category, err := s.repo.FindCategoryById(ctx, id)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return repo.FindCategoryByIdRow{}, ErrCategoryNotFound
		}
		return repo.FindCategoryByIdRow{}, err
	}
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
	return nil
}

func (s *svc) ListProductsByCategory(ctx context.Context, categoryId int64) ([]repo.ListProductsByCategoryRow, error) {
	return s.repo.ListProductsByCategory(ctx, categoryId)
}

func (s *svc) AddProductToCategory(ctx context.Context, categoryId int64, productId int64) error {
	return s.repo.AddProductToCategory(ctx, repo.AddProductToCategoryParams{
		ProductID:  productId,
		CategoryID: categoryId,
	})
}

func (s *svc) RemoveProductFromCategory(ctx context.Context, categoryId int64, productId int64) error {
	return s.repo.RemoveProductFromCategory(ctx, repo.RemoveProductFromCategoryParams{
		ProductID:  productId,
		CategoryID: categoryId,
	})
}

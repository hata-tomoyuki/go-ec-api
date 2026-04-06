package cache

import (
	"fmt"
	"time"
)

// --- Key Prefixes ---

const (
	ProductListPrefix      = "products:list:"
	ProductDetailPrefix    = "products:detail:"
	CategoryListPrefix     = "categories:list:"
	CategoryDetailPrefix   = "categories:detail:"
	CategoryProductsPrefix = "categories:products:"
)

// --- TTLs ---

const (
	ProductListTTL      = 2 * time.Minute
	ProductDetailTTL    = 5 * time.Minute
	CategoryListTTL     = 5 * time.Minute
	CategoryDetailTTL   = 5 * time.Minute
	CategoryProductsTTL = 3 * time.Minute
)

// ProductListKey は商品一覧のキャッシュキーを生成する。
// search が指定されている場合はカーディナリティが高くヒット率が低いため
// キャッシュをスキップする（空文字を返す）。
func ProductListKey(page, limit int, sort, search string, categoryID int) string {
	if search != "" {
		return ""
	}
	return fmt.Sprintf("%spage=%d:limit=%d:sort=%s:cat=%d", ProductListPrefix, page, limit, sort, categoryID)
}

func ProductDetailKey(id int64) string {
	return fmt.Sprintf("%s%d", ProductDetailPrefix, id)
}

func CategoryListKey(page, limit int) string {
	return fmt.Sprintf("%spage=%d:limit=%d", CategoryListPrefix, page, limit)
}

func CategoryDetailKey(id int64) string {
	return fmt.Sprintf("%s%d", CategoryDetailPrefix, id)
}

func CategoryProductsKey(categoryID int64) string {
	return fmt.Sprintf("%s%d", CategoryProductsPrefix, categoryID)
}

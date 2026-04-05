package categories

import (
	"errors"
	"log/slog"
	"net/http"
	"strconv"

	"example.com/ecommerce/internal/json"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateCategories(w http.ResponseWriter, r *http.Request) {
	var tempCategory createCategoryParams
	if err := json.Read(r, &tempCategory); err != nil {
		slog.ErrorContext(r.Context(), "failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := tempCategory.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdCategory, err := h.service.CreateCategories(r.Context(), tempCategory.Name, tempCategory.Description, tempCategory.ImageColor)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create category", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusCreated, createdCategory)
}

func (h *handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	page, _ := strconv.Atoi(r.URL.Query().Get("page"))
	limit, _ := strconv.Atoi(r.URL.Query().Get("limit"))

	params := listCategoriesParams{Page: page, Limit: limit}
	if err := params.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.ListCategoriesPaginated(r.Context(), params)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list categories", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *handler) FindCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.service.FindCategoryById(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		slog.ErrorContext(r.Context(), "failed to find category", "error", err, "id", id)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, category)
}

func (h *handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var tempCategory createCategoryParams
	if err := json.Read(r, &tempCategory); err != nil {
		slog.ErrorContext(r.Context(), "failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := tempCategory.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedCategory, err := h.service.UpdateCategories(r.Context(), id, tempCategory.Name, tempCategory.Description, tempCategory.ImageColor)
	if err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		slog.ErrorContext(r.Context(), "failed to update category", "error", err, "id", id)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, updatedCategory)
}

func (h *handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.service.DeleteCategory(r.Context(), id); err != nil {
		if errors.Is(err, ErrCategoryNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		slog.ErrorContext(r.Context(), "failed to delete category", "error", err, "id", id)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) ListProductsByCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	products, err := h.service.ListProductsByCategory(r.Context(), id)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list products by category", "error", err, "category_id", id)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) AddProductToCategory(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var temp struct {
		ProductID int64 `json:"product_id"`
	}

	if err := json.Read(r, &temp); err != nil {
		slog.ErrorContext(r.Context(), "failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.AddProductToCategory(r.Context(), categoryId, temp.ProductID); err != nil {
		slog.ErrorContext(r.Context(), "failed to add product to category", "error", err, "category_id", categoryId, "product_id", temp.ProductID)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Product added to category successfully"})
}

func (h *handler) RemoveProductFromCategory(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	productIdStr := chi.URLParam(r, "productId")
	productId, err := strconv.ParseInt(productIdStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.RemoveProductFromCategory(r.Context(), categoryId, productId); err != nil {
		slog.ErrorContext(r.Context(), "failed to remove product from category", "error", err, "category_id", categoryId, "product_id", productId)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

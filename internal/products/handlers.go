package products

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

func (h *handler) ListProduct(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()

	page, _ := strconv.Atoi(q.Get("page"))
	limit, _ := strconv.Atoi(q.Get("limit"))

	categoryID, _ := strconv.Atoi(q.Get("category_id"))

	params := listProductsParams{
		Page:       page,
		Limit:      limit,
		Sort:       q.Get("sort"),
		Search:     q.Get("search"),
		CategoryID: categoryID,
	}

	if err := params.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	result, err := h.service.ListProductsPaginated(r.Context(), params)
	if err != nil {
		slog.Error("failed to list products", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, result)
}

func (h *handler) FindProductById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.FindProductById(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		slog.Error("failed to find product", "error", err, "id", id)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var tempProduct createProductParams
	if err := json.Read(r, &tempProduct); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := tempProduct.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdProduct, err := h.service.CreateProduct(r.Context(), tempProduct)
	if err != nil {
		slog.Error("failed to create product", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusCreated, createdProduct)
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var tempProduct updateProductParams
	if err := json.Read(r, &tempProduct); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tempProduct.ID = id

	if err := tempProduct.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedProduct, err := h.service.UpdateProduct(r.Context(), tempProduct)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		slog.Error("failed to update product", "error", err, "id", id)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, updatedProduct)
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		slog.Error("failed to delete product", "error", err, "id", id)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

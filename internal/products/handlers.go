package products

import (
	"errors"
	"log"
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
	products, err := h.service.ListProducts(r.Context())
	if err != nil {
		log.Printf("Error listing products: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, products)
}

func (h *handler) FindProductById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	product, err := h.service.FindProductById(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("Error finding product: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, product)
}

func (h *handler) CreateProduct(w http.ResponseWriter, r *http.Request) {
	var tempProduct createProductParams
	if err := json.Read(r, &tempProduct); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdProduct, err := h.service.CreateProduct(r.Context(), tempProduct)
	if err != nil {
		log.Printf("Error creating product: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, createdProduct)
}

func (h *handler) UpdateProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	var tempProduct updateProductParams
	if err := json.Read(r, &tempProduct); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tempProduct.ID = id

	updatedProduct, err := h.service.UpdateProduct(r.Context(), tempProduct)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("Error updating product: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, updatedProduct)
}

func (h *handler) DeleteProduct(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	err = h.service.DeleteProduct(r.Context(), id)
	if err != nil {
		if errors.Is(err, ErrProductNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		log.Printf("Error deleting product: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

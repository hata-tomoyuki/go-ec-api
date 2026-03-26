package products

import (
	"log"
	"net/http"

	"example.com/ecommerce/internal/json"
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
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusOK, products)
}

package categories

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

func (h *handler) CreateCategories(w http.ResponseWriter, r *http.Request) {
	var tempCategory createCategoryParams
	if err := json.Read(r, &tempCategory); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdCategory, err := h.service.CreateCategories(r.Context(), tempCategory.Name, tempCategory.Description)
	if err != nil {
		log.Printf("Error creating category: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, createdCategory)
}

func (h *handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.service.ListCategories(r.Context())
	if err != nil {
		log.Printf("Error listing categories: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, categories)
}

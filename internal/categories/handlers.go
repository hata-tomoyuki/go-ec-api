package categories

import (
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

func (h *handler) FindCategoryById(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing category ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	category, err := h.service.FindCategoryById(r.Context(), id)

	if err != nil {
		log.Printf("Error finding category: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, category)
}

func (h *handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing category ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var tempCategory createCategoryParams
	if err := json.Read(r, &tempCategory); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedCategory, err := h.service.UpdateCategories(r.Context(), id, tempCategory.Name, tempCategory.Description)
	if err != nil {
		log.Printf("Error updating category: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, updatedCategory)
}

func (h *handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.ParseInt(idStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing category ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	if err := h.service.DeleteCategory(r.Context(), id); err != nil {
		log.Printf("Error deleting category: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

func (h *handler) AddProductToCategory(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing category ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	var temp struct {
		ProductID int64 `json:"product_id"`
	}

	if err := json.Read(r, &temp); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := h.service.AddProductToCategory(r.Context(), categoryId, temp.ProductID); err != nil {
		log.Printf("Error adding product to category: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Product added to category successfully"})
}

func (h *handler) RemoveProductFromCategory(w http.ResponseWriter, r *http.Request) {
	categoryIdStr := chi.URLParam(r, "id")
	categoryId, err := strconv.ParseInt(categoryIdStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing category ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid category ID")
		return
	}

	productIdStr := chi.URLParam(r, "productId")
	productId, err := strconv.ParseInt(productIdStr, 10, 64)
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := h.service.RemoveProductFromCategory(r.Context(), categoryId, productId); err != nil {
		log.Printf("Error removing product from category: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}
}

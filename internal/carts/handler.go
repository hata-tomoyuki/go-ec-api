package carts

import (
	"log"
	"net/http"
	"strconv"

	"example.com/ecommerce/internal/json"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/jwtauth/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateCart(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	sub, ok := claims["sub"].(string)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid customer ID in token claims")
		return
	}

	createdCart, err := h.service.CreateCart(r.Context(), userID)
	if err != nil {
		log.Printf("Error creating cart: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, createdCart)
}

func (h *handler) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	var params addItemToCartParams
	if err := json.Read(r, &params); err != nil {
		log.Printf("Error reading request body: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	addedItem, err := h.service.AddItemToCart(r.Context(), params.CartID, params.ProductID, params.Quantity)
	if err != nil {
		log.Printf("Error adding item to cart: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, addedItem)
}

func (h *handler) ShowCartItems(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	sub, ok := claims["sub"].(string)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid customer ID in token claims")
		return
	}

	items, err := h.service.ListCartItemsByUserId(r.Context(), userID)
	if err != nil {
		log.Printf("Error listing cart items: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, items)
}

func (h *handler) UpdateCartItemQuantity(w http.ResponseWriter, r *http.Request) {
	var params addItemToCartParams
	if err := json.Read(r, &params); err != nil {
		log.Printf("Error reading request body: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedItem, err := h.service.UpdateCartItemQuantity(r.Context(), params.ProductID, params.Quantity)
	if err != nil {
		log.Printf("Error updating cart item quantity: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, updatedItem)
}

func (h *handler) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "id")
	productID, err := strconv.ParseInt(productId, 10, 64)
	if err != nil {
		log.Printf("Error parsing product ID: %v", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	removedItem, err := h.service.RemoveItemFromCart(r.Context(), productID)
	if err != nil {
		log.Printf("Error removing item from cart: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, removedItem)
}

func (h *handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	sub, ok := claims["sub"].(string)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid customer ID in token claims")
		return
	}

	if err := h.service.ClearCart(r.Context(), userID); err != nil {
		log.Printf("Error clearing cart: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Cart cleared successfully"})
}

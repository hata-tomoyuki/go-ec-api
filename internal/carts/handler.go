package carts

import (
	"log/slog"
	"net/http"
	"strconv"

	"example.com/ecommerce/internal/auth"
	"example.com/ecommerce/internal/json"
	"github.com/go-chi/chi/v5"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{service: service}
}

func (h *handler) CreateCart(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	createdCart, err := h.service.CreateCart(r.Context(), userID)
	if err != nil {
		slog.Error("failed to create cart", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusCreated, createdCart)
}

func (h *handler) AddItemToCart(w http.ResponseWriter, r *http.Request) {
	var params addItemToCartParams
	if err := json.Read(r, &params); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	addedItem, err := h.service.AddItemToCart(r.Context(), params.CartID, params.ProductID, params.Quantity)
	if err != nil {
		slog.Error("failed to add item to cart", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, addedItem)
}

func (h *handler) ShowCartItems(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	items, err := h.service.ListCartItemsByUserId(r.Context(), userID)
	if err != nil {
		slog.Error("failed to list cart items", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, items)
}

func (h *handler) UpdateCartItemQuantity(w http.ResponseWriter, r *http.Request) {
	var params addItemToCartParams
	if err := json.Read(r, &params); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedItem, err := h.service.UpdateCartItemQuantity(r.Context(), params.ProductID, params.Quantity)
	if err != nil {
		slog.Error("failed to update cart item quantity", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, updatedItem)
}

func (h *handler) RemoveItemFromCart(w http.ResponseWriter, r *http.Request) {
	productId := chi.URLParam(r, "id")
	productID, err := strconv.ParseInt(productId, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	removedItem, err := h.service.RemoveItemFromCart(r.Context(), productID)
	if err != nil {
		slog.Error("failed to remove item from cart", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, removedItem)
}

func (h *handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	if err := h.service.ClearCart(r.Context(), userID); err != nil {
		slog.Error("failed to clear cart", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Cart cleared successfully"})
}

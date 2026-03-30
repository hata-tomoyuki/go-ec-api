package orders

import (
	"errors"
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
	return &handler{
		service: service,
	}
}

func (h *handler) ListOrdersByCustomerID(w http.ResponseWriter, r *http.Request) {
	customerID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	orders, err := h.service.ListOrdersByCustomerID(r.Context(), customerID)
	if err != nil {
		slog.Error("failed to list orders", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, orders)
}

func (h *handler) ListAllOrders(w http.ResponseWriter, r *http.Request) {
	orders, err := h.service.ListAllOrders(r.Context())
	if err != nil {
		slog.Error("failed to list all orders", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, orders)
}

func (h *handler) FindOrderById(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	customerID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	order, err := h.service.FindOrderById(r.Context(), orderID)
	if err != nil {
		if errors.Is(err, ErrOrderNotFound) {
			json.WriteError(w, http.StatusNotFound, err.Error())
			return
		}
		slog.Error("failed to find order", "error", err, "order_id", orderID)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	if order.CustomerID != customerID {
		json.WriteError(w, http.StatusForbidden, ErrOrderForbidden.Error())
		return
	}

	json.Write(w, http.StatusOK, order)
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	customerID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// リクエストボディの customer_id は無視し、JWTから取得したIDで上書きする。
	// これにより、認証済みユーザー以外の名義で注文が作られることを防ぐ。
	tempOrder.CustomerID = customerID

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderEmptyItems), errors.Is(err, ErrOrderItemValidation):
			json.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrorProductNotFound):
			json.WriteError(w, http.StatusNotFound, err.Error())
		default:
			slog.Error("failed to place order", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}

func (h *handler) CancelOrder(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	customerID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	_, err = h.service.CancelOrder(r.Context(), orderID, customerID)
	if err != nil {
		switch {
		case errors.Is(err, ErrOrderNotFound):
			json.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrOrderForbidden):
			json.WriteError(w, http.StatusForbidden, err.Error())
		case errors.Is(err, ErrOrderNotPending):
			json.WriteError(w, http.StatusBadRequest, err.Error())
		default:
			slog.Error("failed to cancel order", "error", err, "order_id", orderID)
			json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) UpdateOrderStatus(w http.ResponseWriter, r *http.Request) {
	orderID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid order ID")
		return
	}

	var req struct {
		Status string `json:"status"`
	}
	if err := json.Read(r, &req); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedOrder, err := h.service.UpdateOrderStatus(r.Context(), orderID, req.Status)
	if err != nil {
		switch {
		case errors.Is(err, ErrInvalidStatus):
			json.WriteError(w, http.StatusBadRequest, err.Error())
		case errors.Is(err, ErrOrderNotFound):
			json.WriteError(w, http.StatusNotFound, err.Error())
		default:
			slog.Error("failed to update order status", "error", err, "order_id", orderID)
			json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	json.Write(w, http.StatusOK, updatedOrder)
}

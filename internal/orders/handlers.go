package orders

import (
	"log"
	"net/http"
	"strconv"

	"example.com/ecommerce/internal/json"
	"github.com/go-chi/jwtauth/v5"
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
	_, claims, _ := jwtauth.FromContext(r.Context())

	sub, ok := claims["sub"].(string)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	customerID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid customer ID in token claims")
		return
	}

	orders, err := h.service.ListOrdersByCustomerID(r.Context(), customerID)
	if err != nil {
		log.Printf("Error listing orders: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, orders)
}

func (h *handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	var tempOrder createOrderParams
	if err := json.Read(r, &tempOrder); err != nil {
		log.Println("Error reading request body:", err)
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	createdOrder, err := h.service.PlaceOrder(r.Context(), tempOrder)
	if err != nil {
		log.Printf("Error placing order: %v", err)

		if err == ErrorProductNotFound {
			http.Error(w, err.Error(), http.StatusNotFound)
			return
		}

		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	json.Write(w, http.StatusCreated, createdOrder)
}

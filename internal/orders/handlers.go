package orders

import (
	"log"
	"net/http"

	"example.com/ecommerce/internal/json"
)

type handler struct {
	service Service
}

func NewHandler(service Service) *handler {
	return &handler{
		service: service,
	}
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

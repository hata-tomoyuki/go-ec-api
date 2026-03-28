package auth

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

func (h *handler) RegisterUser(w http.ResponseWriter, r *http.Request) {
	var tempUser registerParams
	if err := json.Read(r, &tempUser); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdUser, err := h.service.RegisterUser(r.Context(), tempUser)
	if err != nil {
		log.Printf("Error registering user: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, userResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
	})
}

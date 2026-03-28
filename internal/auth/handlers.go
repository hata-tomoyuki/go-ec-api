package auth

import (
	"errors"
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

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var tempUser loginParams
	if err := json.Read(r, &tempUser); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	token, err := h.service.Login(r.Context(), tempUser)
	if err != nil {
		log.Printf("Error logging in user: %v", err)
		if errors.Is(err, ErrInvalidCredentials) {
			json.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"token": token})
}

package address

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
	return &handler{service: service}
}

func (h *handler) FindAddressByUserId(w http.ResponseWriter, r *http.Request) {
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

	address, err := h.service.FindAddressByUserId(r.Context(), userID)
	if err != nil {
		log.Printf("Error finding address: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, address)
}

func (h *handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
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

	var tempAddress createAddressParams
	if err := json.Read(r, &tempAddress); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdAddress, err := h.service.CreateAddress(r.Context(), userID, tempAddress)
	if err != nil {
		log.Printf("Error creating address: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, createdAddress)
}

func (h *handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
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

	var tempAddress createAddressParams
	if err := json.Read(r, &tempAddress); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedAddress, err := h.service.UpdateAddress(r.Context(), userID, tempAddress)
	if err != nil {
		log.Printf("Error updating address: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, updatedAddress)
}

func (h *handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
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

	deletedAddress, err := h.service.DeleteAddress(r.Context(), userID)
	if err != nil {
		log.Printf("Error deleting address: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, deletedAddress)
}

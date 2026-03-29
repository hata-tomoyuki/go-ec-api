package address

import (
	"errors"
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

// getUserID extracts the user ID from JWT claims.
func getUserID(r *http.Request) (int64, error) {
	_, claims, _ := jwtauth.FromContext(r.Context())
	sub, ok := claims["sub"].(string)
	if !ok {
		return 0, errors.New("invalid token claims")
	}
	return strconv.ParseInt(sub, 10, 64)
}

func (h *handler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	addresses, err := h.service.ListAddressesByUserId(r.Context(), userID)
	if err != nil {
		log.Printf("Error listing addresses: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, addresses)
}

func (h *handler) FindAddressById(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	addressID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	addr, err := h.service.FindAddressById(r.Context(), userID, int32(addressID))
	if err != nil {
		switch {
		case errors.Is(err, ErrAddressNotFound):
			json.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrForbidden):
			json.WriteError(w, http.StatusForbidden, err.Error())
		default:
			log.Printf("Error finding address: %v", err)
			json.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	json.Write(w, http.StatusOK, addr)
}

func (h *handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	var params createAddressParams
	if err := json.Read(r, &params); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := params.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdAddress, err := h.service.CreateAddress(r.Context(), userID, params)
	if err != nil {
		log.Printf("Error creating address: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusCreated, createdAddress)
}

func (h *handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	addressID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	var params createAddressParams
	if err := json.Read(r, &params); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := params.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	updatedAddress, err := h.service.UpdateAddress(r.Context(), userID, int32(addressID), params)
	if err != nil {
		switch {
		case errors.Is(err, ErrAddressNotFound):
			json.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrForbidden):
			json.WriteError(w, http.StatusForbidden, err.Error())
		default:
			log.Printf("Error updating address: %v", err)
			json.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	json.Write(w, http.StatusOK, updatedAddress)
}

func (h *handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := getUserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	addressID, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 32)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid address ID")
		return
	}

	err = h.service.DeleteAddress(r.Context(), userID, int32(addressID))
	if err != nil {
		switch {
		case errors.Is(err, ErrAddressNotFound):
			json.WriteError(w, http.StatusNotFound, err.Error())
		case errors.Is(err, ErrForbidden):
			json.WriteError(w, http.StatusForbidden, err.Error())
		default:
			log.Printf("Error deleting address: %v", err)
			json.WriteError(w, http.StatusInternalServerError, err.Error())
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

package address

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
	return &handler{service: service}
}

func (h *handler) ListAddresses(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	addresses, err := h.service.ListAddressesByUserId(r.Context(), userID)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to list addresses", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, addresses)
}

func (h *handler) FindAddressById(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserID(r)
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
			slog.ErrorContext(r.Context(), "failed to find address", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	json.Write(w, http.StatusOK, addr)
}

func (h *handler) CreateAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	var params createAddressParams
	if err := json.Read(r, &params); err != nil {
		slog.ErrorContext(r.Context(), "failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	if err := params.validate(); err != nil {
		json.WriteError(w, http.StatusBadRequest, err.Error())
		return
	}

	createdAddress, err := h.service.CreateAddress(r.Context(), userID, params)
	if err != nil {
		slog.ErrorContext(r.Context(), "failed to create address", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusCreated, createdAddress)
}

func (h *handler) UpdateAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserID(r)
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
		slog.ErrorContext(r.Context(), "failed to read request body", "error", err)
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
			slog.ErrorContext(r.Context(), "failed to update address", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	json.Write(w, http.StatusOK, updatedAddress)
}

func (h *handler) DeleteAddress(w http.ResponseWriter, r *http.Request) {
	userID, err := auth.UserID(r)
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
			slog.ErrorContext(r.Context(), "failed to delete address", "error", err)
			json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

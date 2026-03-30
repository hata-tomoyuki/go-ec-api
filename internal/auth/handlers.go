package auth

import (
	"errors"
	"log/slog"
	"net/http"

	"example.com/ecommerce/internal/json"
	"github.com/jackc/pgx/v5"
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
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	createdUser, err := h.service.RegisterUser(r.Context(), tempUser)
	if err != nil {
		slog.Error("failed to register user", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusCreated, userResponse{
		ID:    createdUser.ID,
		Name:  createdUser.Name,
		Email: createdUser.Email,
		Role:  string(createdUser.Role),
	})
}

func (h *handler) Login(w http.ResponseWriter, r *http.Request) {
	var tempUser loginParams
	if err := json.Read(r, &tempUser); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tokens, err := h.service.Login(r.Context(), tempUser)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			json.WriteError(w, http.StatusUnauthorized, "Invalid email or password")
			return
		}
		slog.Error("failed to login", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, tokens)
}

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	lc, err := GetLogoutClaims(r)
	if err != nil {
		json.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	if err := h.service.Logout(r.Context(), lc.JTI, lc.ExpiredAt, lc.RefreshTokenID); err != nil {
		slog.Error("failed to logout", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (h *handler) Refresh(w http.ResponseWriter, r *http.Request) {
	var body struct {
		RefreshToken string `json:"refresh_token"`
	}
	if err := json.Read(r, &body); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	tokens, err := h.service.Refresh(r.Context(), body.RefreshToken)
	if err != nil {
		if errors.Is(err, ErrInvalidRefreshToken) {
			json.WriteError(w, http.StatusUnauthorized, "Invalid or expired refresh token")
			return
		}
		slog.Error("failed to refresh token", "error", err)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, tokens)
}

func (h *handler) GetMe(w http.ResponseWriter, r *http.Request) {
	userID, err := UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	user, err := h.service.GetProfile(r.Context(), userID)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			json.WriteError(w, http.StatusNotFound, "User not found")
			return
		}
		slog.Error("failed to load profile", "error", err, "user_id", userID)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, userResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
		Role:  string(user.Role),
	})
}

func (h *handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	userID, err := UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	var params updateUserParams
	if err := json.Read(r, &params); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), userID, params)
	if err != nil {
		slog.Error("failed to update user", "error", err, "user_id", userID)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, userResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Role:  string(updatedUser.Role),
	})
}

func (h *handler) UpdatePassword(w http.ResponseWriter, r *http.Request) {
	userID, err := UserID(r)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	var params struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.Read(r, &params); err != nil {
		slog.Error("failed to read request body", "error", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedUser, err := h.service.UpdateUserPassword(r.Context(), userID, params.CurrentPassword, params.NewPassword)
	if err != nil {
		if errors.Is(err, ErrInvalidCredentials) {
			json.WriteError(w, http.StatusUnauthorized, "Current password is incorrect")
			return
		}
		slog.Error("failed to update password", "error", err, "user_id", userID)
		json.WriteError(w, http.StatusInternalServerError, "Internal server error")
		return
	}

	json.Write(w, http.StatusOK, userResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Role:  string(updatedUser.Role),
	})
}

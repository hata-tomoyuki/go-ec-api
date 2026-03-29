package auth

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

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
		Role:  string(createdUser.Role),
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

func (h *handler) Logout(w http.ResponseWriter, r *http.Request) {
	token, claims, err := jwtauth.FromContext(r.Context())
	if err != nil || token == nil {
		json.WriteError(w, http.StatusUnauthorized, "Unauthorized")
		return
	}

	jti, ok := claims["jti"].(string)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	exp, ok := claims["exp"].(float64)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	expiredAt := time.Unix(int64(exp), 0)
	if err := h.service.Logout(r.Context(), jti, expiredAt); err != nil {
		log.Printf("Error logging out user: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, map[string]string{"message": "Logged out successfully"})
}

func (h *handler) GetMe(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	sub, ok := claims["sub"].(string)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	name, _ := claims["name"].(string)
	email, _ := claims["email"].(string)
	role, _ := claims["role"].(string)

	json.Write(w, http.StatusOK, userResponse{
		ID:    userID,
		Name:  name,
		Email: email,
		Role:  role,
	})
}

func (h *handler) UpdateMe(w http.ResponseWriter, r *http.Request) {
	_, claims, _ := jwtauth.FromContext(r.Context())

	sub, ok := claims["sub"].(string)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	var params updateUserParams
	if err := json.Read(r, &params); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedUser, err := h.service.UpdateUser(r.Context(), userID, params)
	if err != nil {
		log.Printf("Error updating user: %v", err)
		json.WriteError(w, http.StatusInternalServerError, err.Error())
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
	_, claims, _ := jwtauth.FromContext(r.Context())

	sub, ok := claims["sub"].(string)
	if !ok {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	userID, err := strconv.ParseInt(sub, 10, 64)
	if err != nil {
		json.WriteError(w, http.StatusBadRequest, "Invalid token claims")
		return
	}

	var params struct {
		CurrentPassword string `json:"current_password"`
		NewPassword     string `json:"new_password"`
	}
	if err := json.Read(r, &params); err != nil {
		log.Println("Error reading request body:", err)
		json.WriteError(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	updatedUser, err := h.service.UpdateUserPassword(r.Context(), userID, params.CurrentPassword, params.NewPassword)
	if err != nil {
		log.Printf("Error updating user password: %v", err)
		if errors.Is(err, ErrInvalidCredentials) {
			json.WriteError(w, http.StatusUnauthorized, "Current password is incorrect")
			return
		}
		json.WriteError(w, http.StatusInternalServerError, err.Error())
		return
	}

	json.Write(w, http.StatusOK, userResponse{
		ID:    updatedUser.ID,
		Name:  updatedUser.Name,
		Email: updatedUser.Email,
		Role:  string(updatedUser.Role),
	})
}

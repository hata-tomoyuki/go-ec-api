package auth

import (
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

func generateJWT(ja *jwtauth.JWTAuth, userID int64, email, role string) (string, error) {
	claims := map[string]interface{}{
		"sub":   fmt.Sprintf("%d", userID),
		"email": email,
		"role":  role,
		"jti":   uuid.New().String(),
	}
	jwtauth.SetExpiryIn(claims, 24*time.Hour)

	_, tokenString, err := ja.Encode(claims)
	return tokenString, err
}

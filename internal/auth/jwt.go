package auth

import (
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

func generateJWT(ja *jwtauth.JWTAuth, userID int64, email string) (string, error) {
	claims := map[string]interface{}{
		"sub":   fmt.Sprintf("%d", userID),
		"email": email,
		"jti":   uuid.New().String(),
	}
	jwtauth.SetExpiryIn(claims, 24*time.Hour)

	_, tokenString, err := ja.Encode(claims)
	return tokenString, err
}

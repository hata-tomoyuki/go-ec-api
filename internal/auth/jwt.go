package auth

import (
	"fmt"
	"time"

	"github.com/go-chi/jwtauth/v5"
	"github.com/google/uuid"
)

const accessTokenTTL = 15 * time.Minute

func generateJWT(ja *jwtauth.JWTAuth, userID int64, name, email, role string, refreshTokenID int64) (string, error) {
	claims := map[string]interface{}{
		"sub":   fmt.Sprintf("%d", userID),
		"name":  name,
		"email": email,
		"role":  role,
		"jti":   uuid.New().String(),
		"rtid":  fmt.Sprintf("%d", refreshTokenID),
	}
	jwtauth.SetExpiryIn(claims, accessTokenTTL)

	_, tokenString, err := ja.Encode(claims)
	return tokenString, err
}

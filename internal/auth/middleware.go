package auth

import (
	"encoding/json"
	"net/http"

	"github.com/go-chi/jwtauth/v5"
)

func JWTAuthenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, _, err := jwtauth.FromContext(r.Context())

		w.Header().Set("Content-Type", "application/json")

		if err != nil || token == nil {
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(map[string]string{"error": "Unauthorized"})
			return
		}

		next.ServeHTTP(w, r)
	})
}

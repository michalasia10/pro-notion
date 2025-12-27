package middleware

import (
	"fmt"
	"net/http"
	"strings"

	"src/internal/config"
	"src/internal/pkg/httpx"

	"github.com/golang-jwt/jwt/v5"
)

// JWTAuthMiddleware validates JWT tokens and sets user context
func JWTAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authHeader := r.Header.Get("Authorization")
		if authHeader == "" {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Missing authorization header"})
			return
		}

		// Extract token from "Bearer <token>" format
		tokenString := strings.TrimPrefix(authHeader, "Bearer ")
		if tokenString == authHeader || tokenString == "" {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid authorization header format"})
			return
		}

		// Parse and validate token
		cfg := config.Get()
		token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
			if token.Method != jwt.SigningMethodHS256 {
				return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
			}
			return []byte(cfg.JWT.Secret), nil
		})

		if err != nil {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token"})
			return
		}

		if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
			if cfg.JWT.Issuer != "" && claims.Issuer != cfg.JWT.Issuer {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token issuer"})
				return
			}
			if cfg.JWT.Audience != "" && !containsAudience(claims.Audience, cfg.JWT.Audience) {
				httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token audience"})
				return
			}
			// Set user ID in context
			ctx := SetUserID(r.Context(), claims.UserID)
			next.ServeHTTP(w, r.WithContext(ctx))
		} else {
			httpx.WriteJSON(w, http.StatusUnauthorized, map[string]string{"error": "Invalid token claims"})
		}
	})
}

func containsAudience(aud jwt.ClaimStrings, want string) bool {
	for _, a := range aud {
		if a == want {
			return true
		}
	}
	return false
}

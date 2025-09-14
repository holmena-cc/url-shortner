package server

import (
	"context"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
)
type contextKey string
const userIDKey = contextKey("userId")

var jwtKey = []byte(os.Getenv("JWT_KEY")) // load key from env

type Claims struct {
	UserID int32 `json:"userId"`
	jwt.RegisteredClaims
}

// Generate JWT token
func GenerateToken(userID int32) (string, error) {
	claims := &Claims{
		UserID: userID,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(jwtKey)
}

// AuthMiddleware checks JWT in HttpOnly cookie
func AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		claims := &Claims{}
		tkn, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})
		if err != nil || !tkn.Valid {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}
		ctx := context.WithValue(r.Context(),userIDKey, claims.UserID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// AuthMiddlewareOptional attaches userID to context if token exists but never forces login
func AuthMiddlewareOptional(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err == nil {
			claims := &Claims{}
			tkn, err := jwt.ParseWithClaims(cookie.Value, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})
			if err == nil && tkn.Valid {
				// Token is valid → attach userID to context
				r = r.WithContext(context.WithValue(r.Context(), userIDKey, claims.UserID))
			}
		}
		// Continue to handler no matter what
		next.ServeHTTP(w, r)
	})
}

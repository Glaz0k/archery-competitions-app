package delivery

import (
	"context"
	"net/http"
	"strings"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key_my_secret_key_my_secret_key") // where to get

type Claims struct {
	UserID int    `json:"user_id"`
	Role   string `json:"role"`
	jwt.StandardClaims
}

func JWTRoleMiddleware(roles string) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			tokenString := r.Header.Get("Authorization")
			if tokenString == "" {
				http.Error(w, "Missing authorization token", http.StatusUnauthorized)
				return
			}

			if strings.ToUpper(tokenString[0:7]) == "BEARER " {
				tokenString = tokenString[7:]
			}

			claims := &Claims{}
			token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
				return jwtKey, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Invalid token", http.StatusUnauthorized)
				return
			}

			if !strings.Contains(roles, claims.Role) {
				http.Error(w, "Invalid role", http.StatusForbidden)
				return
			}

			r = r.WithContext(context.WithValue(r.Context(), "user_id", claims.UserID))
			r = r.WithContext(context.WithValue(r.Context(), "role", claims.Role))
			next.ServeHTTP(w, r)
		})
	}
}

package delivery

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key") // where to get

type Claims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	jwt.StandardClaims
}

// without Bearer
func JWTMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		tokenString := r.Header.Get("Authorization")
		if tokenString == "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Missing authorization token"))
			return
		}

		claims := &Claims{}
		token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
			return jwtKey, nil
		})

		if err != nil || !token.Valid {
			log.Print(err)
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("Invalid token"))
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "claims", claims)
		r = r.WithContext(ctx)

		next.ServeHTTP(w, r)
	})
}

func RoleMiddleware(role string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := r.Context().Value("claims").(*Claims)
			if !ok || claims.Role != role {
				w.WriteHeader(http.StatusForbidden)
				w.Write([]byte(fmt.Sprintf("Access denied. Required role: %s", role)))
				return
			}

			next.ServeHTTP(w, r)
		})
	}
}

// for local use
// eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VybmFtZSI6IkpvbiIsInJvbGUiOiJ1c2VyIiwiZXhwIjoxNzQxODkzMDg3fQ.U8hdL0AouhRC8NgfMowuK8ZaAz_h-F-Fpic3nGb8zrU
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	username := r.URL.Query().Get("username")
	role := r.URL.Query().Get("role")

	if username == "" || (role != "user" && role != "admin") {
		w.WriteHeader(http.StatusBadRequest)
		w.Write([]byte("Invalid username or role"))
		return
	}

	expirationTime := time.Now().Add(15 * time.Minute)
	claims := &Claims{
		Username: username,
		Role:     role,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Error generating token"))
		return
	}

	w.Write([]byte(tokenString))
}

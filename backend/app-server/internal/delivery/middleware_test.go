package delivery_test

import (
	"app-server/internal/delivery"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/dgrijalva/jwt-go"
)

var jwtKey = []byte("my_secret_key") // where to get

func generateTestToken(role string) string {
	claims := &delivery.Claims{
		Role: role,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		panic(err)
	}
	return tokenString
}
func TestJWTRoleMiddleware(t *testing.T) {
	handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("Access granted"))
	})

	tests := []struct {
		name           string
		role           string
		token          string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "No token",
			role:           "admin",
			token:          "",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Missing authorization token",
		},
		{
			name:           "Invalid token",
			role:           "admin",
			token:          "invalid-token",
			expectedStatus: http.StatusUnauthorized,
			expectedBody:   "Invalid token",
		},
		{
			name:           "Wrong role",
			role:           "admin",
			token:          generateTestToken("user"),
			expectedStatus: http.StatusForbidden,
			expectedBody:   "Access denied. Required role: admin",
		},
		{
			name:           "Correct role",
			role:           "admin",
			token:          generateTestToken("admin"),
			expectedStatus: http.StatusOK,
			expectedBody:   "Access granted",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаем запрос
			req := httptest.NewRequest("GET", "/", nil)
			if tt.token != "" {
				req.Header.Set("Authorization", "Bearer "+tt.token)
			}

			rr := httptest.NewRecorder()

			middleware := delivery.JWTRoleMiddleware(tt.role)
			middleware(handler).ServeHTTP(rr, req)

			if status := rr.Code; status != tt.expectedStatus {
				t.Errorf("Ожидался статус %d, получен %d", tt.expectedStatus, status)
			}

			if body := rr.Body.String(); body != tt.expectedBody {
				t.Errorf("Ожидалось тело ответа %q, получено %q", tt.expectedBody, body)
			}
		})
	}
}

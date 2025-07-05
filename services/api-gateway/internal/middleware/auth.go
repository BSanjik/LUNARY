package middleware

import (
	"context"
	"errors"
	"net/http"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

type ctxKey string

const UserContextKey = ctxKey("user")

func JWTMiddleware(secret []byte, exemptPaths []string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Освобождаемые маршруты
			if isExempt(r.URL.Path, exemptPaths) {
				next.ServeHTTP(w, r)
				return
			}

			// Проверка заголовка
			authHeader := r.Header.Get("Authorization")
			if authHeader == "" {
				http.Error(w, "Unauthorized: missing token", http.StatusUnauthorized)
				return
			}

			if !strings.HasPrefix(authHeader, "Bearer ") {
				http.Error(w, "Unauthorized: invalid token format", http.StatusUnauthorized)
				return
			}

			tokenString := strings.TrimPrefix(authHeader, "Bearer ")

			// Парсинг токена
			token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, errors.New("unexpected signing method")
				}
				return secret, nil
			})

			if err != nil || !token.Valid {
				http.Error(w, "Unauthorized: invalid token", http.StatusUnauthorized)
				return
			}

			// Добавляем claims в контекст
			ctx := context.WithValue(r.Context(), UserContextKey, token.Claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Проверка, освобождён ли путь от проверки токена
func isExempt(path string, exemptPaths []string) bool {
	for _, prefix := range exemptPaths {
		if strings.HasPrefix(path, prefix) {
			return true
		}
	}
	return false
}

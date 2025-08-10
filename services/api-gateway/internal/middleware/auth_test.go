package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/stretchr/testify/assert"
)

func TestJWTMiddleware(t *testing.T) {
	secret := []byte("test-secret")
	exemptPaths := []string{"/auth/login", "/auth/registration"}

	// Создаём валидный JWT токен для тестов
	validToken := createTestJWT(t, secret, time.Now().Add(time.Hour))
	expiredToken := createTestJWT(t, secret, time.Now().Add(-time.Hour))

	tests := []struct {
		name           string
		path           string
		authHeader     string
		expectedStatus int
		expectedError  string
	}{
		{
			name:           "Успешная авторизация с валидным токеном",
			path:           "/protected",
			authHeader:     "Bearer " + validToken,
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Освобождённый путь - логин",
			path:           "/auth/login",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Освобождённый путь - регистрация",
			path:           "/auth/registration",
			authHeader:     "",
			expectedStatus: http.StatusOK,
		},
		{
			name:           "Отсутствует токен",
			path:           "/protected",
			authHeader:     "",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Требуется авторизация",
		},
		{
			name:           "Неверный формат токена",
			path:           "/protected",
			authHeader:     "Invalid token-format",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Неверный формат токена авторизации",
		},
		{
			name:           "Истёкший токен",
			path:           "/protected",
			authHeader:     "Bearer " + expiredToken,
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Токен авторизации истёк",
		},
		{
			name:           "Невалидный токен",
			path:           "/protected",
			authHeader:     "Bearer invalid.jwt.token",
			expectedStatus: http.StatusUnauthorized,
			expectedError:  "Недействительный токен авторизации",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём handler для тестирования
			handler := JWTMiddleware(secret, exemptPaths)(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				w.WriteHeader(http.StatusOK)
				w.Write([]byte("success"))
			}))

			// Создаём HTTP запрос
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			if tt.authHeader != "" {
				req.Header.Set("Authorization", tt.authHeader)
			}

			// Создаём ResponseRecorder
			rr := httptest.NewRecorder()

			// Выполняем запрос
			handler.ServeHTTP(rr, req)

			// Проверяем результат
			assert.Equal(t, tt.expectedStatus, rr.Code)

			if tt.expectedError != "" {
				assert.Contains(t, rr.Body.String(), tt.expectedError)
			} else if tt.expectedStatus == http.StatusOK {
				assert.Equal(t, "success", rr.Body.String())
			}
		})
	}
}

func TestIsExempt(t *testing.T) {
	exemptPaths := []string{"/auth", "/health", "/docs"}

	tests := []struct {
		path     string
		expected bool
	}{
		{"/auth/login", true},
		{"/auth/registration", true},
		{"/health", true},
		{"/health/live", true},
		{"/docs/swagger", true},
		{"/protected", false},
		{"/user/profile", false},
		{"/scenario/search", false},
	}

	for _, tt := range tests {
		t.Run(tt.path, func(t *testing.T) {
			result := isExempt(tt.path, exemptPaths)
			assert.Equal(t, tt.expected, result)
		})
	}
}

// Вспомогательная функция для создания тестового JWT токена
func createTestJWT(t *testing.T, secret []byte, exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": "test-user-123",
		"exp": exp.Unix(),
		"iat": time.Now().Unix(),
	})

	tokenString, err := token.SignedString(secret)
	assert.NoError(t, err)
	return tokenString
}

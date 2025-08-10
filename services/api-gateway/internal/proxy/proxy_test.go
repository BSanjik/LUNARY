package proxy

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestProxy_ServeHTTP(t *testing.T) {
	// Создаём тестовый backend сервер
	backendServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("backend response: " + r.URL.Path))
	}))
	defer backendServer.Close()

	// Конфигурация сервисов
	services := map[string]string{
		"auth":     backendServer.URL,
		"scenario": backendServer.URL,
	}

	proxy := NewProxy(services)

	tests := []struct {
		name           string
		path           string
		expectedStatus int
		expectedBody   string
	}{
		{
			name:           "Успешное проксирование auth сервиса",
			path:           "/auth/login",
			expectedStatus: http.StatusOK,
			expectedBody:   "backend response: /login",
		},
		{
			name:           "Успешное проксирование scenario сервиса",
			path:           "/scenario/search",
			expectedStatus: http.StatusOK,
			expectedBody:   "backend response: /search",
		},
		{
			name:           "Пустой путь",
			path:           "/",
			expectedStatus: http.StatusBadRequest,
			expectedBody:   "Некорректный путь запроса",
		},
		{
			name:           "Несуществующий сервис",
			path:           "/unknown/endpoint",
			expectedStatus: http.StatusNotFound,
			expectedBody:   "Сервис не найден",
		},
		{
			name:           "Глубокий путь",
			path:           "/auth/user/profile/settings",
			expectedStatus: http.StatusOK,
			expectedBody:   "backend response: /user/profile/settings",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, tt.path, nil)
			rr := httptest.NewRecorder()

			proxy.ServeHTTP(rr, req)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Contains(t, rr.Body.String(), tt.expectedBody)
		})
	}
}

func TestNewProxy(t *testing.T) {
	services := map[string]string{
		"auth":     "http://auth-service:8081",
		"scenario": "http://scenario-service:8082",
	}

	proxy := NewProxy(services)

	assert.NotNil(t, proxy)
	assert.Equal(t, services, proxy.Services)
}

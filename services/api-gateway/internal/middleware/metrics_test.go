package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestMetricsMiddleware(t *testing.T) {
	tests := []struct {
		name             string
		handlerStatus    int
		handlerResponse  string
		handlerDelay     time.Duration
		expectedStatus   int
		expectedResponse string
	}{
		{
			name:             "Успешный запрос",
			handlerStatus:    http.StatusOK,
			handlerResponse:  "success",
			expectedStatus:   http.StatusOK,
			expectedResponse: "success",
		},
		{
			name:             "Ошибка клиента",
			handlerStatus:    http.StatusBadRequest,
			handlerResponse:  "bad request",
			expectedStatus:   http.StatusBadRequest,
			expectedResponse: "bad request",
		},
		{
			name:             "Ошибка сервера",
			handlerStatus:    http.StatusInternalServerError,
			handlerResponse:  "internal error",
			expectedStatus:   http.StatusInternalServerError,
			expectedResponse: "internal error",
		},
		{
			name:             "Медленный запрос",
			handlerStatus:    http.StatusOK,
			handlerResponse:  "slow response",
			handlerDelay:     2 * time.Second,
			expectedStatus:   http.StatusOK,
			expectedResponse: "slow response",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Создаём тестовый handler
			testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
				if tt.handlerDelay > 0 {
					time.Sleep(tt.handlerDelay)
				}
				w.WriteHeader(tt.handlerStatus)
				w.Write([]byte(tt.handlerResponse))
			})

			// Оборачиваем в metrics middleware
			handler := MetricsMiddleware()(testHandler)

			req := httptest.NewRequest(http.MethodGet, "/test", nil)
			rr := httptest.NewRecorder()

			start := time.Now()
			handler.ServeHTTP(rr, req)
			duration := time.Since(start)

			assert.Equal(t, tt.expectedStatus, rr.Code)
			assert.Equal(t, tt.expectedResponse, rr.Body.String())

			// Проверяем, что медленные запросы действительно медленные
			if tt.handlerDelay > 0 {
				assert.True(t, duration >= tt.handlerDelay)
			}
		})
	}
}

func TestRequestIDMiddleware(t *testing.T) {
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("test"))
	})

	handler := RequestIDMiddleware()(testHandler)

	req := httptest.NewRequest(http.MethodGet, "/test", nil)
	rr := httptest.NewRecorder()

	handler.ServeHTTP(rr, req)

	assert.Equal(t, http.StatusOK, rr.Code)
	assert.Equal(t, "test", rr.Body.String())

	// Проверяем, что добавлен X-Request-ID header
	requestID := rr.Header().Get("X-Request-ID")
	assert.NotEmpty(t, requestID)
	assert.True(t, len(requestID) > 10) // Простая проверка формата
}

func TestGenerateRequestID(t *testing.T) {
	id1 := generateRequestID()
	id2 := generateRequestID()

	assert.NotEmpty(t, id1)
	assert.NotEmpty(t, id2)
	assert.NotEqual(t, id1, id2)  // Разные ID для разных вызовов
	assert.True(t, len(id1) > 10) // Минимальная длина
}

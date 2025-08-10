package middleware

import (
	"net/http"
	"time"

	"github.com/BSanjik/LUNARY/services/common/shared"
)

// responseWriter обёртка для захвата статус кода
type responseWriter struct {
	http.ResponseWriter
	statusCode int
	written    int64
}

func (rw *responseWriter) WriteHeader(code int) {
	rw.statusCode = code
	rw.ResponseWriter.WriteHeader(code)
}

func (rw *responseWriter) Write(b []byte) (int, error) {
	if rw.statusCode == 0 {
		rw.statusCode = 200
	}
	n, err := rw.ResponseWriter.Write(b)
	rw.written += int64(n)
	return n, err
}

// MetricsMiddleware добавляет мониторинг производительности
func MetricsMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			// Обернуть ResponseWriter для захвата метрик
			wrapped := &responseWriter{
				ResponseWriter: w,
				statusCode:     200,
			}

			// Обработка запроса
			next.ServeHTTP(wrapped, r)

			// Вычисление метрик
			duration := time.Since(start)

			// Structured logging для HTTP метрик
			shared.LogHTTPRequest(
				r.Method,
				r.URL.Path,
				r.UserAgent(),
				r.RemoteAddr,
				wrapped.statusCode,
				duration.Milliseconds(),
			)

			// Дополнительные метрики для мониторинга
			shared.LogInfo("HTTP запрос завершён", map[string]interface{}{
				"method":        r.Method,
				"path":          r.URL.Path,
				"status_code":   wrapped.statusCode,
				"duration_ms":   duration.Milliseconds(),
				"bytes_written": wrapped.written,
				"remote_addr":   r.RemoteAddr,
				"user_agent":    r.UserAgent(),
				"type":          "http_metrics",
			})

			// Логирование медленных запросов (> 1 секунды)
			if duration > time.Second {
				shared.LogWarn("Медленный HTTP запрос", map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"duration_ms": duration.Milliseconds(),
					"status_code": wrapped.statusCode,
					"type":        "slow_request",
				})
			}

			// Логирование ошибок 4xx и 5xx
			if wrapped.statusCode >= 400 {
				logLevel := "warn"
				if wrapped.statusCode >= 500 {
					logLevel = "error"
				}

				fields := map[string]interface{}{
					"method":      r.Method,
					"path":        r.URL.Path,
					"status_code": wrapped.statusCode,
					"duration_ms": duration.Milliseconds(),
					"remote_addr": r.RemoteAddr,
					"type":        "http_error",
				}

				if logLevel == "error" {
					shared.LogError(nil, "HTTP ошибка сервера", fields)
				} else {
					shared.LogWarn("HTTP ошибка клиента", fields)
				}
			}
		})
	}
}

// RequestIDMiddleware добавляет уникальный ID к каждому запросу
func RequestIDMiddleware() func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Генерируем простой request ID (можно использовать UUID)
			requestID := generateRequestID()

			// Добавляем в заголовок ответа
			w.Header().Set("X-Request-ID", requestID)

			// Добавляем в контекст для дальнейшего использования
			ctx := shared.WithRequestID(r.Context(), requestID)

			// Логируем с контекстом (request_id автоматически добавится)
			shared.LogInfoWithContext(ctx, "Новый HTTP запрос", map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
			})

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

// Простая генерация request ID (в production лучше использовать UUID)
func generateRequestID() string {
	return time.Now().Format("20060102150405") + "-" +
		string(rune(time.Now().Nanosecond()%26+97)) +
		string(rune(time.Now().Nanosecond()%26+97))
}

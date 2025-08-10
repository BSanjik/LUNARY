package middleware

import (
	"net/http"
	"time"

	"github.com/BSanjik/LUNARY/services/common/shared"
)

// loggingResponseWriter обёртка для логирования ответов
type loggingResponseWriter struct {
	http.ResponseWriter
	statusCode int
	size       int
}

func (lrw *loggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *loggingResponseWriter) Write(b []byte) (int, error) {
	if lrw.statusCode == 0 {
		lrw.statusCode = 200
	}
	size, err := lrw.ResponseWriter.Write(b)
	lrw.size += size
	return size, err
}

// LoggingMiddleware логирует все входящие HTTP запросы
func LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		// Логирование входящего запроса
		shared.LogInfo("Входящий HTTP запрос", map[string]interface{}{
			"method":          r.Method,
			"uri":             r.RequestURI,
			"path":            r.URL.Path,
			"remote_addr":     r.RemoteAddr,
			"user_agent":      r.UserAgent(),
			"content_type":    r.Header.Get("Content-Type"),
			"content_length":  r.ContentLength,
			"host":            r.Host,
			"referer":         r.Referer(),
			"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
			"type":            "request_start",
		})

		// Обернуть ResponseWriter для логирования ответа
		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		// Обработка запроса
		next.ServeHTTP(lrw, r)

		// Вычисление времени обработки
		duration := time.Since(start)

		// Определение уровня логирования по статус коду
		logFields := map[string]interface{}{
			"method":        r.Method,
			"uri":           r.RequestURI,
			"path":          r.URL.Path,
			"status_code":   lrw.statusCode,
			"duration_ms":   duration.Milliseconds(),
			"response_size": lrw.size,
			"remote_addr":   r.RemoteAddr,
			"user_agent":    r.UserAgent(),
			"type":          "request_completed",
		}

		// Логирование завершения запроса с соответствующим уровнем
		switch {
		case lrw.statusCode >= 500:
			shared.LogError(nil, "HTTP запрос завершён с ошибкой сервера", logFields)
		case lrw.statusCode >= 400:
			shared.LogWarn("HTTP запрос завершён с ошибкой клиента", logFields)
		case lrw.statusCode >= 300:
			shared.LogInfo("HTTP запрос завершён с перенаправлением", logFields)
		default:
			shared.LogInfo("HTTP запрос успешно завершён", logFields)
		}

		// Дополнительное логирование для аномально медленных запросов
		if duration > 5*time.Second {
			shared.LogWarn("Очень медленный HTTP запрос", map[string]interface{}{
				"method":      r.Method,
				"path":        r.URL.Path,
				"duration_ms": duration.Milliseconds(),
				"threshold":   "5000ms",
				"type":        "performance_alert",
			})
		}

		// Логирование подозрительной активности
		if r.Header.Get("User-Agent") == "" {
			shared.LogWarn("Запрос без User-Agent", map[string]interface{}{
				"path":        r.URL.Path,
				"remote_addr": r.RemoteAddr,
				"type":        "security_alert",
			})
		}

		// Логирование больших запросов/ответов
		if r.ContentLength > 10*1024*1024 { // > 10MB
			shared.LogWarn("Большой входящий запрос", map[string]interface{}{
				"path":           r.URL.Path,
				"content_length": r.ContentLength,
				"type":           "large_request",
			})
		}

		if lrw.size > 10*1024*1024 { // > 10MB
			shared.LogWarn("Большой исходящий ответ", map[string]interface{}{
				"path":          r.URL.Path,
				"response_size": lrw.size,
				"type":          "large_response",
			})
		}
	})
}

// AccessLogMiddleware специализированное логирование доступа (в стиле Apache/Nginx)
func AccessLogMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		lrw := &loggingResponseWriter{
			ResponseWriter: w,
			statusCode:     200,
		}

		next.ServeHTTP(lrw, r)

		// Логирование в формате access log
		shared.LogInfo("ACCESS_LOG", map[string]interface{}{
			"remote_addr":   r.RemoteAddr,
			"method":        r.Method,
			"uri":           r.RequestURI,
			"protocol":      r.Proto,
			"status_code":   lrw.statusCode,
			"response_size": lrw.size,
			"referer":       r.Referer(),
			"user_agent":    r.UserAgent(),
			"duration_ms":   time.Since(start).Milliseconds(),
			"timestamp":     start.Format("02/Jan/2006:15:04:05 -0700"),
			"type":          "access_log",
		})
	})
}

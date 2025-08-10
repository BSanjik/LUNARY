package middleware

import (
	"net/http"
	"time"

	"golang.org/x/time/rate"

	"github.com/BSanjik/LUNARY/services/common/shared"
)

func RateLimitMiddleware(limiter *rate.Limiter) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			start := time.Now()

			if !limiter.Allow() {
				// Логирование превышения лимита
				shared.LogWarn("Превышен лимит запросов", map[string]interface{}{
					"remote_addr":     r.RemoteAddr,
					"user_agent":      r.UserAgent(),
					"path":            r.URL.Path,
					"method":          r.Method,
					"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
					"type":            "rate_limit_exceeded",
					"timestamp":       start.Format(time.RFC3339),
				})

				// Возвращаем structured error вместо http.Error
				shared.WriteErrorResponse(w, shared.NewRateLimitError("Слишком много запросов. Попробуйте позже"))
				return
			}

			// Логирование разрешённого запроса (только для debug)
			shared.LogDebug("Запрос разрешён rate limiter", map[string]interface{}{
				"remote_addr": r.RemoteAddr,
				"path":        r.URL.Path,
				"method":      r.Method,
			})

			next.ServeHTTP(w, r)
		})
	}
}

// AdvancedRateLimitMiddleware с более сложной логикой (опционально)
func AdvancedRateLimitMiddleware(limiter *rate.Limiter, burst int) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// Получаем количество доступных токенов
			if !limiter.Allow() {
				// Дополнительная информация для мониторинга
				shared.LogWarn("Rate limit: запрос отклонён", map[string]interface{}{
					"remote_addr":     r.RemoteAddr,
					"path":            r.URL.Path,
					"method":          r.Method,
					"user_agent":      r.UserAgent(),
					"referer":         r.Referer(),
					"x_forwarded_for": r.Header.Get("X-Forwarded-For"),
					"burst_limit":     burst,
					"type":            "rate_limit_violation",
				})

				// Добавляем headers для клиента
				w.Header().Set("X-RateLimit-Limit", "5")
				w.Header().Set("X-RateLimit-Remaining", "0")
				w.Header().Set("Retry-After", "60")

				shared.WriteErrorResponse(w, shared.NewRateLimitError("Превышен лимит запросов. Повторите попытку через 60 секунд"))
				return
			}

			// Добавляем информационные headers
			w.Header().Set("X-RateLimit-Limit", "5")
			w.Header().Set("X-RateLimit-Remaining", "4") // Можно рассчитать точнее

			next.ServeHTTP(w, r)
		})
	}
}

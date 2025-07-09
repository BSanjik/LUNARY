// Определение маршрутов и проксирование запросов
package internal

import (
	"net/http"
	"time"

	"github.com/BSanjik/LUNARY/services/api-gateway/internal/config"
	"github.com/BSanjik/LUNARY/services/api-gateway/internal/middleware"
	"github.com/BSanjik/LUNARY/services/api-gateway/internal/proxy"

	"golang.org/x/time/rate"
)

func NewServer(cfg *config.Config) *http.Server {
	limiter := rate.NewLimiter(5, 10)

	mux := http.NewServeMux()

	// Маршруты, которые не требуют JWT
	exempt := []string{"/auth/registration", "/auth/login"}

	// Обёртываем proxy
	proxyHandler := proxy.NewProxy(cfg.Services)

	// Оборачиваем middleware
	handler := middleware.LoggingMiddleware(
		middleware.RateLimitMiddleware(limiter)(
			middleware.JWTMiddleware([]byte(cfg.JWTSecret), exempt)(
				proxyHandler,
			),
		),
	)

	// Регистрируем "/" после middleware
	mux.Handle("/", handler)

	return &http.Server{
		Addr:         cfg.ListnerAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

}

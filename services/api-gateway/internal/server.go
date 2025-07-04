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

	proxyHandler := proxy.NewProxy(cfg.Services)

	mux := http.NewServeMux()

	handler := middleware.LoggingMiddleware(middleware.RateLimitMiddleware(limiter)(middleware.JWTMiddleware([]byte(cfg.JWTSecret), []string{"/auth"})(proxyHandler)))

	mux.Handle("/", handler)

	return &http.Server{
		Addr:         cfg.ListnerAddr,
		Handler:      mux,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

}

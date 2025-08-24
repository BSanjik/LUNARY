// Точка входа user сервиса
package main

import (
	"context"
	"net/http"
	"os"
	"time"
	"user-service/internal/handler"
	"user-service/internal/service"
	"user-service/internal/storage"
	"user-service/internal/token"

	"github.com/BSanjik/LUNARY/services/common/shared"
	"github.com/gorilla/mux"
)

func loadConfig() *storage.Config {
	return &storage.Config{
		Port:           getEnv("PORT", "8083"),
		DBHost:         getEnv("DB_HOST", "localhost"),
		DBPort:         getEnv("DB_PORT", "5432"),
		DBUser:         getEnv("DB_USER", "postgres"),
		DBPassword:     getEnv("DB_PASSWORD", "postgres"),
		DBName:         getEnv("DB_NAME", "lunary_db"),
		DBSSLMode:      getEnv("DB_SSL_MODE", "disable"),
		JWTSecret:      getEnv("JWT_SECRET", "your-secret-key"),
		AuthServiceURL: getEnv("AUTH_SERVICE_URL", "http://localhost:8081"),
	}
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	// Загружаем конфигурацию
	cfg := loadConfig()

	// Инициализируем структурированное логирование
	shared.InitLogger("user-service", "info")

	shared.LogInfo("Starting user service", map[string]interface{}{
		"port":    cfg.Port,
		"version": "1.0.0",
	})

	// Подключаемся к базе данных
	store, err := storage.NewStorage(cfg)
	if err != nil {
		shared.LogError(err, "Failed to connect to database")
		return
	}
	defer store.DB.Close()

	// Инициализируем JWT
	token.Init(cfg.JWTSecret)

	// Создаем сервисы
	userService := service.NewUserService(store)
	h := &handler.Handler{Service: userService}

	// Создаем health manager
	healthManager := shared.NewHealthManager("user-service", "1.0.0")

	// Добавляем проверку БД
	dbChecker := shared.NewDatabaseHealthChecker(store.DB)
	healthManager.AddChecker(dbChecker)

	// Создаем роутер
	r := mux.NewRouter()

	// API endpoints
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/profile", h.AuthMiddleware(h.CreateProfile)).Methods("POST")
	apiRouter.HandleFunc("/profile", h.AuthMiddleware(h.GetProfile)).Methods("GET")
	apiRouter.HandleFunc("/profile", h.AuthMiddleware(h.UpdateProfile)).Methods("PUT")
	apiRouter.HandleFunc("/profile", h.AuthMiddleware(h.DeleteProfile)).Methods("DELETE")

	// Health check endpoints
	r.HandleFunc("/health", healthManager.HealthHandler).Methods("GET")
	r.HandleFunc("/health/live", healthManager.LivenessHandler).Methods("GET")
	r.HandleFunc("/health/ready", healthManager.ReadinessHandler).Methods("GET")

	// CORS middleware
	r.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
			w.Header().Set("Access-Control-Allow-Headers", "Origin, Content-Type, Accept, Authorization")

			if r.Method == "OPTIONS" {
				w.WriteHeader(http.StatusNoContent)
				return
			}

			next.ServeHTTP(w, r)
		})
	})

	// Создаем HTTP сервер
	server := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	// Создаем graceful shutdown менеджер
	graceful := shared.NewGracefulShutdown(server, 30*time.Second)

	// Добавляем функции очистки
	graceful.AddCleanup(shared.DatabaseCleanup(store.DB, "user-database"))

	// Добавляем очистку токенов
	graceful.AddCleanup(shared.ContextCleanup("token-cleanup", func(ctx context.Context) error {
		shared.LogInfo("Cleaning up token cache")
		return nil
	}))

	// Ожидаем завершения активных соединений
	graceful.AddCleanup(shared.WaitForActiveConnections(server, 10*time.Second))

	shared.LogInfo("User service configured successfully", map[string]interface{}{
		"port":             cfg.Port,
		"api_endpoints":    []string{"/api/v1/profile"},
		"health_endpoints": []string{"/health", "/health/live", "/health/ready"},
		"shutdown_timeout": "30s",
		"read_timeout":     "15s",
		"write_timeout":    "15s",
		"idle_timeout":     "60s",
	})

	// Запускаем сервер с graceful shutdown
	if err := graceful.Start(); err != nil {
		shared.LogError(err, "Server startup failed")
	}
}

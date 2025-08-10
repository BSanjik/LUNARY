// Точка входа auth сервиса
// @title LUNARY Auth Service API
// @version 1.0
// @description API сервиса авторизации LUNARY
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"context"
	"net/http"
	"time"

	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"auth-service/internal/token"

	"github.com/BSanjik/LUNARY/services/common/shared"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"

	_ "auth-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.LoadConfig()

	// Инициализируем структурированное логирование
	shared.InitLogger("auth-service", "info")

	shared.LogInfo("Starting auth service", map[string]interface{}{
		"port":    cfg.Port,
		"version": "1.0.0",
	})

	store, err := storage.NewStorage(cfg)
	if err != nil {
		shared.LogError(err, "Failed to connect to database")
		return
	}
	defer store.DB.Close()

	token.Init(cfg.JWTSecret)

	// Создаем validator
	validate := validator.New()

	svc := &service.AuthService{
		Storage:   store,
		Validator: validate,
	}
	h := &handler.Handler{Service: svc}

	// Создаем health manager
	healthManager := shared.NewHealthManager("auth-service", "1.0.0")

	// Добавляем проверку БД
	dbChecker := shared.NewDatabaseHealthChecker(store.DB)
	healthManager.AddChecker(dbChecker)

	r := mux.NewRouter()

	// API endpoints
	apiRouter := r.PathPrefix("/api/v1").Subrouter()
	apiRouter.HandleFunc("/registration", h.Registration).Methods("POST")
	apiRouter.HandleFunc("/login", h.Login).Methods("POST")

	// Health check endpoints
	r.HandleFunc("/health", healthManager.HealthHandler).Methods("GET")
	r.HandleFunc("/health/live", healthManager.LivenessHandler).Methods("GET")
	r.HandleFunc("/health/ready", healthManager.ReadinessHandler).Methods("GET")

	// Swagger
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

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
	graceful.AddCleanup(shared.DatabaseCleanup(store.DB, "auth-database"))

	// Добавляем очистку токенов (если есть кэш)
	graceful.AddCleanup(shared.ContextCleanup("token-cleanup", func(ctx context.Context) error {
		shared.LogInfo("Cleaning up token cache")
		// Здесь можно очистить кэш токенов, если он есть
		return nil
	}))

	// Ожидаем завершения активных соединений
	graceful.AddCleanup(shared.WaitForActiveConnections(server, 10*time.Second))

	shared.LogInfo("Auth service configured successfully", map[string]interface{}{
		"port":             cfg.Port,
		"api_endpoints":    []string{"/api/v1/registration", "/api/v1/login"},
		"health_endpoints": []string{"/health", "/health/live", "/health/ready"},
		"swagger":          "/swagger/",
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

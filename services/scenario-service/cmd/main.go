// Точка входа Scenario сервиса
package main

import (
	"context"
	"database/sql"
	"net/http"
	"scenario-service/internal/config"
	"scenario-service/internal/handler"
	"scenario-service/internal/service"
	"scenario-service/internal/storage"
	"time"

	"github.com/BSanjik/LUNARY/services/common/shared"
	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

func main() {
	// Инициализация логирования
	shared.InitLogger("scenario-service", "info")

	cfg := config.LoadConfig()

	// Подключение к БД
	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		shared.LogError(err, "Не удалось подключиться к базе данных")
		return
	}
	defer db.Close()

	// Проверка подключения к БД
	if err := db.Ping(); err != nil {
		shared.LogError(err, "Не удалось проверить подключение к базе данных")
		return
	}

	// Инициализация валидатора
	validator := validator.New()

	// Инициализация слоёв
	store := &storage.Storage{DB: db}
	svc := &service.ScenarioService{
		Storage:   store,
		Validator: validator,
	}
	h := &handler.Handler{Service: svc}

	// Настройка маршрутов
	r := mux.NewRouter()

	// Бизнес эндпоинты
	r.HandleFunc("/search", h.GetScenario).Methods("POST")
	r.HandleFunc("/feedback", h.SubmitFeedback).Methods("POST")

	// Настройка health checks
	healthManager := shared.NewHealthManager("scenario-service", "1.0.0")
	dbHealthChecker := shared.NewDatabaseHealthChecker(db)
	healthManager.AddChecker(dbHealthChecker)

	// Health check эндпоинты
	r.HandleFunc("/health", healthManager.HealthHandler).Methods("GET")
	r.HandleFunc("/health/live", healthManager.LivenessHandler).Methods("GET")
	r.HandleFunc("/health/ready", healthManager.ReadinessHandler).Methods("GET")

	// Настройка HTTP сервера
	server := &http.Server{
		Addr:    "0.0.0.0:" + cfg.Port,
		Handler: r,
	}

	// Настройка graceful shutdown
	graceful := shared.NewGracefulShutdown(server, 30*time.Second)
	graceful.AddCleanup(func(ctx context.Context) error {
		return db.Close()
	})

	shared.LogInfo("Scenario-service запущен на порту " + cfg.Port)

	// Запуск сервера с graceful shutdown
	graceful.Start()
}

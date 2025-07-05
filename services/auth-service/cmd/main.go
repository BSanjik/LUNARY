// Точка входа auth сервиса
// @title LUNARY Auth Service API
// @version 1.0
// @description API сервиса авторизации LUNARY
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/BSanjik/LUNARY/services/auth-service/internal/config"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/db"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/handler"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/logger"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/token"

	"github.com/gorilla/mux"

	_ "github.com/BSanjik/LUNARY/services/auth-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	// 1. Конфиг и логгер
	cfg := config.LoadConfig()
	logger.InitLogger()
	logger.Log.Info("Logger initialized")

	// 2. Инициализация БД
	conn, err := db.InitDB(cfg)
	if err != nil {
		logger.Log.Fatalw("DB connection failed", "error", err)
	}
	defer conn.Close()

	// 3. Инициализация JWT
	tokenService := token.New(cfg.JWTSecret)

	// 4. Auth handler
	authHandler := handler.NewAuthHandler(conn, tokenService)

	// 5. Router
	r := mux.NewRouter()
	api := r.PathPrefix("/api/v1").Subrouter()
	api.HandleFunc("/register", authHandler.Register).Methods("POST")
	api.HandleFunc("/login", authHandler.Login).Methods("POST")
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	// 6. HTTP Server с graceful shutdown
	srv := &http.Server{
		Addr:         ":" + cfg.Port,
		Handler:      r,
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  30 * time.Second,
	}

	// 7. Завершение по Ctrl+C
	go func() {
		logger.Log.Infow("Starting server", "port", cfg.Port)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Log.Fatalw("Listen failed", "error", err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	logger.Log.Infow("Shutting down server...")
	if err := srv.Close(); err != nil {
		logger.Log.Errorw("Server shutdown error", "error", err)
	}
}

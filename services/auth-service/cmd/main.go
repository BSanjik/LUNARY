// Точка входа auth сервиса
// @title LUNARY Auth Service API
// @version 1.0
// @description API сервиса авторизации LUNARY
// @host localhost:8080
// @BasePath /api/v1
package main

import (
	"log"
	"net/http"

	"auth-service/internal/config"
	"auth-service/internal/handler"
	"auth-service/internal/service"
	"auth-service/internal/storage"
	"auth-service/internal/token"

	"github.com/gorilla/mux"

	_ "auth-service/docs"

	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.LoadConfig()
	store, err := storage.NewStorage(cfg)
	if err != nil {
		log.Fatal("DB error: ", err)
	}
	defer store.DB.Close()

	token.Init(cfg.JWTSecret)

	svc := &service.AuthService{Storage: store}
	h := &handler.Handler{Service: svc}

	r := mux.NewRouter()
	r.HandleFunc("/registration", h.Registration).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Auth service running on port", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}

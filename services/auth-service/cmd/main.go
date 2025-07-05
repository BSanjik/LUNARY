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

	"github.com/BSanjik/LUNARY/services/auth-service/internal/config"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/db"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/handler"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/token"

	"github.com/gorilla/mux"

	_ "github.com/BSanjik/LUNARY/services/auth-service/docs"
	httpSwagger "github.com/swaggo/http-swagger"
)

func main() {
	cfg := config.LoadConfig()
	conn, err := db.InitDB(cfg)
	if err != nil {
		log.Fatal("DB error: ", err)
	}
	defer conn.Close()

	token.Init(cfg.JWTSecret)

	h := &handler.AuthHandler{DB: conn}

	r := mux.NewRouter()
	r.HandleFunc("/register", h.Register).Methods("POST")
	r.HandleFunc("/login", h.Login).Methods("POST")

	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler)

	log.Println("Auth service running on port", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}

// Точка входа auth сервиса
package main

import (
	"log"
	"net/http"

	"github.com/BSanjik/LUNARY/services/auth-service/internal/config"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/db"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/handler"
	"github.com/BSanjik/LUNARY/services/auth-service/internal/token"
	"github.com/gorilla/mux"
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

	log.Println("Auth service running on port", cfg.Port)
	http.ListenAndServe(":"+cfg.Port, r)
}

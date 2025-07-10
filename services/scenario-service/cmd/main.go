// Точка входа Scenario сервиса
package main

import (
	"database/sql"
	"log"
	"net/http"
	"scenario-service/internal/config"
	"scenario-service/internal/handler"
	"scenario-service/internal/service"
	"scenario-service/internal/storage"

	"github.com/gorilla/mux"
)

func main() {
	cfg := config.LoadConfig()

	db, err := sql.Open("postgres", cfg.DBURL)
	if err != nil {
		log.Fatalf("Не удалось подключиться к БД: %v", err)
	}
	defer db.Close()

	store := &storage.Storage{DB: db}
	svc := &service.ScenarioService{Storage: store}
	h := &handler.Handler{Service: svc}

	r := mux.NewRouter()
	r.HandleFunc("/scenario", h.GetScenario).Methods("POST")

	log.Printf("Scenario-service started on port %s", cfg.Port)
	log.Fatal(http.ListenAndServe(":"+cfg.Port, r))
}

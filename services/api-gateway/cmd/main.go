// Точка входа, запускает сервер API Gateway
package main

import (
	"log"

	"api-gateway/internal"
	"api-gateway/internal/config"
)

func main() {
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	srv := internal.NewServer(cfg)
	log.Printf("API Gateway started at %s\n", cfg.ListnerAddr)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}

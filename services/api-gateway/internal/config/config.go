// Конфигурация сервиса
package config

import (
	"errors"
	"log"
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	JWTSecret   string
	ListnerAddr string
	Services    map[string]string
}

func LoadConfig() (*Config, error) {
	if err := godotenv.Load(); err != nil {
		log.Println(".env not found, continue loading from environment")
	}

	jwtSecret := os.Getenv("JWT_SECRET")
	if jwtSecret == "" {
		return nil, errors.New("JWT_SECRET env var required")
	}

	listenAddr := os.Getenv("LISTEN_ADDR")
	if listenAddr == " " {
		listenAddr = ":8080"
	}

	services := map[string]string{
		"auth":     getEnvOrDefault("SERVICE_AUTH", "http://localhost:8081"),
		"scenario": getEnvOrDefault("SERVICE_AI", "http://localhost:8082"),
	}
	return &Config{
		JWTSecret:   jwtSecret,
		ListnerAddr: listenAddr,
		Services:    services,
	}, nil
}

func getEnvOrDefault(key, def string) string {
	if val := strings.TrimSpace(os.Getenv(key)); val != "" {
		return val
	}
	return def
}

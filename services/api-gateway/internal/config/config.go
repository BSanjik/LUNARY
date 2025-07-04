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
		"auth":  getEnvOrDefault("SERVICE_AUTH", "http://localhost:8001"),
		"user":  getEnvOrDefault("SERVICE_USER", "http://localhost:8002"),
		"ad":    getEnvOrDefault("SERVICE_AD", "http://localhost:8003"),
		"ai":    getEnvOrDefault("SERVICE_AI", "http://localhost:8004"),
		"media": getEnvOrDefault("SERVICE_MEDIA", "http://localhost:8005"),
		"notif": getEnvOrDefault("SERVICE_NOTIF", "http://localhost:8006"),
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

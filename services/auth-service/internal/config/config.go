// Конфигурация
package config

import (
	"log"
	"os"
)

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
}

func LoadConfig() *Config {
	// Пытаемся загрузить .env, но не обязательно
	// Используем только переменные окружения из docker-compose.yml
	log.Println("📋 Загружаем конфигурацию из переменных окружения")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Println("[Config] PORT по умолчанию 8080")
	}

	db := os.Getenv("DB_URL")
	if db == "" {
		log.Fatal("[Config] Не задана переменная DB_URL")
	}

	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		log.Fatal("[Config] Не задана переменная JWT_SECRET")
	}

	return &Config{
		Port:      port,
		DBUrl:     db,
		JWTSecret: secret,
	}
}

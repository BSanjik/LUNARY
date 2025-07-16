// Конфигурация
package config

import (
	"log"
	"os"
)

type Config struct {
	Port  string
	DBURL string
}

func LoadConfig() *Config {
	// Используем только переменные окружения из docker-compose.yml
	log.Println("📋 Загружаем конфигурацию из переменных окружения")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8082"
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("❌ DB_URL is not set in environment")
	}

	return &Config{
		Port:  port,
		DBURL: dbURL,
	}
}

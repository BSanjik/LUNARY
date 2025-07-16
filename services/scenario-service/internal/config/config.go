// Конфигурация
package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBURL string
}

func LoadConfig() *Config {
	err := godotenv.Load()
	if err != nil {
		log.Println("⚠️  .env файл не найден, используем переменные окружения")
	}
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

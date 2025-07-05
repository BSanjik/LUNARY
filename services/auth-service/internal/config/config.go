// Конфигурация
package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port      string
	DBUrl     string
	JWTSecret string
}

func LoadConfig() *Config {
	godotenv.Load(".env")
	return &Config{
		Port:      os.Getenv("PORT"),
		DBUrl:     os.Getenv("DB_URL"),
		JWTSecret: os.Getenv("JWT_SECRET"),
	}
}

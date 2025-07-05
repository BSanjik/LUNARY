// Доступ к БД (PostgreSQL)
package db

import (
	"database/sql"
	"fmt"

	"github.com/BSanjik/LUNARY/services/auth-service/internal/config"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	connsStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.DBHost, cfg.DBPort, cfg.DBUser, cfg.DBPassword, cfg.DBName)
	return sql.Open("postgres", connsStr)
}

// Доступ к БД (PostgreSQL)
package db

import (
	"database/sql"
	"log"
	"time"

	"github.com/BSanjik/LUNARY/services/auth-service/internal/config"

	_ "github.com/lib/pq"
)

func InitDB(cfg *config.Config) (*sql.DB, error) {
	log.Println("[DB] Connecting to:", cfg.DBUrl) //после проверки удалить

	db, err := sql.Open("postgres", cfg.DBUrl)
	if err != nil {
		log.Println("[DB] Failed to open DB:", err)
		return nil, err
	}

	//connection check
	for i := 1; i <= 10; i++ {
		err := db.Ping()
		if err == nil {
			log.Println("[DB] Successfully connected")
			return db, nil
		}
		log.Printf("[DB] Ping attempt %d failed: %v", i, err)
		time.Sleep(2 * time.Second)
	}

	log.Println("[DB] Successfully connected")
	return db, nil
}

// Доступ к БД (PostgreSQL)
package storage

import (
	"database/sql"
	"log"
	"time"

	"auth-service/internal/config"

	_ "github.com/lib/pq"
)

type Storage struct {
	DB *sql.DB
}

func NewStorage(cfg *config.Config) (*Storage, error) {

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
			return &Storage{DB: db}, nil
		}
		log.Printf("[DB] Ping attempt %d failed: %v", i, err)
		time.Sleep(2 * time.Second)
	}

	return &Storage{DB: db}, nil
}

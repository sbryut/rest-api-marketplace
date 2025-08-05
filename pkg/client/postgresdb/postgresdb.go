package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"rest-api-marketplace/internal/config"

	_ "github.com/lib/pq"
)

func NewClient(ctx context.Context, cfg config.PostgresConfig) (*sql.DB, error) {

	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		cfg.Host,
		cfg.Port,
		cfg.Username,
		cfg.Password,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	log.Println("Successfully connected to the database")
	return db, nil
}

func CloseDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("error closing database connection: %v", err)
	}
}

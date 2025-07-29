package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log"

	"rest-api-marketplace/internal/config"
)

func NewClient(ctx context.Context, cfg config.PostgresConfig) (*sql.DB, error) {

	connStr := fmt.Sprintf("user=%s port=%s password=%s dbname=%s sslmode=disable",
		cfg.Username,
		cfg.Port,
		cfg.Password,
		cfg.DBName,
	)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %v", err)
	}

	if err := db.PingContext(ctx); err != nil {
		return nil, fmt.Errorf("failed to ping database: %v", err)
	}

	return db, nil
}

func CloseDatabase(db *sql.DB) {
	if err := db.Close(); err != nil {
		log.Printf("error closing database connection: %v", err)
	}
}

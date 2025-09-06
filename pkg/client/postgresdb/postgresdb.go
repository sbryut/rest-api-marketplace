// Package postgres is a package for postgresDB, the service data source,
package postgres

import (
	"context"
	"database/sql"
	"fmt"
	"log/slog"

	"rest-api-marketplace/internal/config"
)

// NewClient initializes and configures a connection to PostgresDB
func NewClient(ctx context.Context, cfg config.PostgresConfig, log *slog.Logger) (*sql.DB, error) {
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
		if cerr := db.Close(); cerr != nil {
			log.Error("failed to close database after ping error", slog.Any("error", cerr))
		}
		return nil, fmt.Errorf("failed to ping database: %w", err)
	}

	log.Info("Successfully connected to the database")
	return db, nil
}

// CloseDatabase closes the connection to the PostgresDB
func CloseDatabase(db *sql.DB, log *slog.Logger) {
	if err := db.Close(); err != nil {
		log.Error("failed to close database connection", slog.Any("error", err))
	}
}

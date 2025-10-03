package main

import (
	"database/sql"
	"fmt"
	"log/slog"
	"os"

	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/thenopholo/shorty_url/internal/server"
	"github.com/thenopholo/shorty_url/internal/store"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute the code", "error", err)
		return
	}
	slog.Info("All systems offline")
}

func run() error {
	// Get database configuration from environment variables
	dbHost := getEnv("DB_HOST", "localhost")
	dbPort := getEnv("DB_PORT", "5432")
	dbUser := getEnv("DB_USER", "admin")
	dbPassword := getEnv("DB_PASSWORD", "password")
	dbName := getEnv("DB_NAME", "shorty_url")
	port := getEnv("PORT", "8080")

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable",
		dbUser, dbPassword, dbHost, dbPort, dbName)

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		return err
	}
	defer db.Close()

	if err := db.Ping(); err != nil {
		slog.Error("failed to connect to database", "error", err)
		return err
	}
	slog.Info("Successfully connected to the database")

	// Create table if not exists
	if err := createTable(db); err != nil {
		return err
	}

	st := store.NewStore(db)

	srv := server.NewServer(server.Config{
		Port:   port,
		Logger: slog.Default(),
		Store:  st,
	})
	srv.SetupRoutes()
	return srv.Start()
}

func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id SERIAL PRIMARY KEY,
		code VARCHAR(10) UNIQUE NOT NULL,
		original_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	)`

	_, err := db.Exec(query)
	if err != nil {
		return fmt.Errorf("failed to create table: %w", err)
	}

	slog.Info("Database table ready")
	return nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

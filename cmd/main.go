package main

import (
	"database/sql"
	"log/slog"

	"github.com/thenopholo/shorty_url/internal/server"
)

func main() {
	if err := run(); err != nil {
		slog.Error("failed to execute the code", "error", err)
		return
	}
	slog.Info("All systems offline")
}

func run() error {
	connStr := "postgres://admin:password@localhost:5432/shorty_url?sslmode=disable"
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

	srv := server.NewServer(server.Config{
		Port: ":9786",
		Logger: slog.Default(),
	})
	srv.SetupRoutes()
	return srv.Start()
}

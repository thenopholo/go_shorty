package store

import (
	"database/sql"
	"errors"
	"fmt"
)

type Store struct {
	db *sql.DB
}

func NewStrore(db *sql.DB) *Store {
	return &Store{db}
}

func (s *Store) SaveURL(code, originalURL string) error {
	query := `INSERT INTO urls (code, original_url) VALUES($1, $2)`
	_, err := s.db.Exec(query, code, originalURL)
	if err != nil {
		return fmt.Errorf("failed to save url: %w", err)
	}

	return nil
}

func (s *Store) GetURL(code string) (string, error) {
	var originalURL string
	query := `SELECT original_url WHERE urls code = $1`
	err := s.db.QueryRow(query, code).Scan(&originalURL)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return "", fmt.Errorf("url not found")
		}
		return "", fmt.Errorf("failed to get url: %w", err)
	}
	return originalURL, nil
}

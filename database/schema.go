package database

import (
	"database/sql"
	"fmt"
)

// CreateTables crée toutes les tables nécessaires si elles n'existent pas.
func CreateTables(db *sql.DB) error {
	queries := []string{
		// Table des utilisateurs
		`CREATE TABLE IF NOT EXISTS users (
			id       INTEGER PRIMARY KEY AUTOINCREMENT,
			username TEXT    NOT NULL UNIQUE,
			email    TEXT    NOT NULL UNIQUE,
			password TEXT    NOT NULL
		)`,

		// Table des sessions
		`CREATE TABLE IF NOT EXISTS sessions (
			id         TEXT    PRIMARY KEY,
			user_id    INTEGER NOT NULL,
			expires_at DATETIME NOT NULL,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// Table des posts
		`CREATE TABLE IF NOT EXISTS posts (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			user_id    INTEGER NOT NULL,
			title      TEXT    NOT NULL,
			content    TEXT    NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("erreur création table : %w", err)
		}
	}

	return nil
}

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
			category   TEXT    NOT NULL DEFAULT 'general',
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,

		// Table des commentaires
		`CREATE TABLE IF NOT EXISTS comments (
			id         INTEGER PRIMARY KEY AUTOINCREMENT,
			post_id    INTEGER NOT NULL,
			user_id    INTEGER NOT NULL,
			content    TEXT    NOT NULL,
			created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
			FOREIGN KEY(post_id) REFERENCES posts(id) ON DELETE CASCADE,
			FOREIGN KEY(user_id) REFERENCES users(id) ON DELETE CASCADE
		)`,
	}

	for _, query := range queries {
		if _, err := db.Exec(query); err != nil {
			return fmt.Errorf("erreur création table : %w", err)
		}
	}

	if err := ensureColumnExists(db, "posts", "category", "TEXT NOT NULL DEFAULT 'general'"); err != nil {
		return err
	}

	return nil
}

func ensureColumnExists(db *sql.DB, table, column, columnDefinition string) error {
	rows, err := db.Query(fmt.Sprintf("PRAGMA table_info(%s)", table))
	if err != nil {
		return fmt.Errorf("impossible de lire les colonnes de %s : %w", table, err)
	}
	defer rows.Close()

	for rows.Next() {
		var cid int
		var name, ctype string
		var notnull int
		var dfltValue sql.NullString
		var pk int
		if err := rows.Scan(&cid, &name, &ctype, &notnull, &dfltValue, &pk); err != nil {
			return fmt.Errorf("impossible de lire la définition des colonnes de %s : %w", table, err)
		}
		if name == column {
			return nil
		}
	}

	if _, err := db.Exec(fmt.Sprintf("ALTER TABLE %s ADD COLUMN %s %s", table, column, columnDefinition)); err != nil {
		return fmt.Errorf("impossible d'ajouter la colonne %s à %s : %w", column, table, err)
	}

	return nil
}

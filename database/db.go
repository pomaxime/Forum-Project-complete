package database

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
)

// Init ouvre la connexion à la base de données SQLite.
func Init() (*sql.DB, error) {
	db, err := sql.Open("sqlite3", "./database.db")
	if err != nil {
		return nil, fmt.Errorf("impossible d'ouvrir la base de données : %w", err)
	}

	// Vérification que la connexion fonctionne vraiment
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("impossible de contacter la base de données : %w", err)
	}

	// Activation des foreign keys (désactivé par défaut dans SQLite)
	if _, err := db.Exec("PRAGMA foreign_keys = ON"); err != nil {
		return nil, fmt.Errorf("impossible d'activer les foreign keys : %w", err)
	}

	return db, nil
}

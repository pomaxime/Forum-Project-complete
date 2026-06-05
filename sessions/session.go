package sessions

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/google/uuid"
)

const DefaultTTL = 24 * time.Hour

type Store struct {
	db  *sql.DB
	ttl time.Duration
}

func NewStore(db *sql.DB, ttl time.Duration) *Store {
	if ttl == 0 {
		ttl = DefaultTTL
	}
	return &Store{db: db, ttl: ttl}
}

func (s *Store) Create(userID int) (string, time.Time, error) {
	sessionID := uuid.New().String()
	expiresAt := time.Now().Add(s.ttl)

	_, err := s.db.Exec(
		"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
		sessionID, userID, expiresAt,
	)
	if err != nil {
		return "", time.Time{}, fmt.Errorf("sessions.Create : %w", err)
	}
	return sessionID, expiresAt, nil
}

func (s *Store) GetUserID(sessionID string) (int, error) {
	var userID int
	var expiresAt time.Time
	err := s.db.QueryRow(
		"SELECT user_id, expires_at FROM sessions WHERE id = ?", sessionID,
	).Scan(&userID, &expiresAt)
	if err != nil {
		return 0, err
	}
	if time.Now().After(expiresAt) {
		// Nettoyage de la session expirée
		s.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
		return 0, sql.ErrNoRows
	}
	return userID, nil
}

func (s *Store) Delete(sessionID string) error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE id = ?", sessionID)
	if err != nil {
		return fmt.Errorf("sessions.Delete : %w", err)
	}
	return nil
}

func (s *Store) DeleteExpired() error {
	_, err := s.db.Exec("DELETE FROM sessions WHERE expires_at < ?", time.Now())
	if err != nil {
		return fmt.Errorf("sessions.DeleteExpired : %w", err)
	}
	return nil
}

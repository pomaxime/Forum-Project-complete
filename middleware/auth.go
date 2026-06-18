package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"time"
)

type contextKey string

const userIDKey contextKey = "userID"

func Auth(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			loginURL := "/login?next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}

		var userID int
		var expiresAt time.Time
		err = db.QueryRow(
			"SELECT user_id, expires_at FROM sessions WHERE id = ?",
			cookie.Value,
		).Scan(&userID, &expiresAt)

		if err == sql.ErrNoRows {
			clearSessionCookie(w)
			loginURL := "/login?next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		if time.Now().After(expiresAt) {
			db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
			clearSessionCookie(w)
			loginURL := "/login?expired=1&next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserID(r *http.Request) int {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0
	}
	return userID
}

func GetUserIDFromCookie(w http.ResponseWriter, db *sql.DB, r *http.Request) int {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return 0
	}

	var userID int
	var expiresAt time.Time
	err = db.QueryRow(
		"SELECT user_id, expires_at FROM sessions WHERE id = ?",
		cookie.Value,
	).Scan(&userID, &expiresAt)
	if err != nil {
		clearSessionCookie(w)
		return 0
	}

	if time.Now().After(expiresAt) {
		db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
		clearSessionCookie(w)
		return 0
	}

	return userID
}

func clearSessionCookie(w http.ResponseWriter) {
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
}

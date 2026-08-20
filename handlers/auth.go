package handlers

import (
	"database/sql"
	"net/http"
	"time"

	"forum/middleware"
)

func deleteSession(db *sql.DB, w http.ResponseWriter, r *http.Request) {
	cookie, err := r.Cookie("session_id")
	if err != nil {
		return
	}
	db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Path:     "/",
		Expires:  time.Unix(0, 0),
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})
}

func isLoggedIn(db *sql.DB, w http.ResponseWriter, r *http.Request) bool {
	return middleware.GetUserIDFromCookie(w, db, r) != 0
}

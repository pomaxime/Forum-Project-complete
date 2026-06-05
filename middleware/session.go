package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"time"
)

func InjectUser(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("session_id")
		if err != nil {
			next.ServeHTTP(w, r)
			return
		}

		var userID int
		var expiresAt time.Time
		err = db.QueryRow(
			"SELECT user_id, expires_at FROM sessions WHERE id = ?",
			cookie.Value,
		).Scan(&userID, &expiresAt)

		if err != nil || time.Now().After(expiresAt) {
			if err == nil {
				db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
			}
			clearSessionCookie(w)
			next.ServeHTTP(w, r)
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

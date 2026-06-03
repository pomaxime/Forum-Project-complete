package middleware

import (
	"context"
	"database/sql"
	"net/http"
	"net/url"
	"time"
)

// contextKey est un type privé pour éviter les collisions dans le contexte HTTP.
type contextKey string

const userIDKey contextKey = "userID"

// Auth est un middleware qui protège les routes nécessitant une authentification.
func Auth(db *sql.DB, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// 1. Lecture du cookie
		cookie, err := r.Cookie("session_id")
		if err != nil {
			loginURL := "/login?next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}

		// 2. Vérification de la session en base de données
		var userID int
		var expiresAt time.Time
		err = db.QueryRow(
			"SELECT user_id, expires_at FROM sessions WHERE id = ?",
			cookie.Value,
		).Scan(&userID, &expiresAt)

		if err == sql.ErrNoRows {
			// Session inconnue
			clearSessionCookie(w)
			loginURL := "/login?next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		// 3. Vérification de l'expiration
		if time.Now().After(expiresAt) {
			// Session expirée : on la supprime de la BDD
			db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
			clearSessionCookie(w)
			loginURL := "/login?expired=1&next=" + url.QueryEscape(r.URL.RequestURI())
			http.Redirect(w, r, loginURL, http.StatusSeeOther)
			return
		}

		// 4. Session valide : injection du user_id dans le contexte de la requête
		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

// GetUserID extrait le user_id du contexte de la requête. Retourne 0 si non trouvé.
func GetUserID(r *http.Request) int {
	userID, ok := r.Context().Value(userIDKey).(int)
	if !ok {
		return 0
	}
	return userID
}

// GetUserIDFromCookie vérifie le cookie de session et renvoie le user_id si la session est valide.
// Si la session est invalide ou expirée, le cookie est effacé.
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
		// Session expirée : suppression côté serveur et client
		db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)
		clearSessionCookie(w)
		return 0
	}

	return userID
}

// clearSessionCookie expire immédiatement le cookie côté navigateur.
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

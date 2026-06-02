package handlers

import (
	"net/http"
	"time"
)

// Logout déconnecte l'utilisateur en supprimant sa session et son cookie.
func (h *AuthHandler) Logout(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	// Lecture du cookie de session
	cookie, err := r.Cookie("session_id")
	if err != nil {
		// Pas de cookie → déjà déconnecté
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	// Suppression de la session en base de données
	h.db.Exec("DELETE FROM sessions WHERE id = ?", cookie.Value)

	// Expiration immédiate du cookie côté navigateur
	http.SetCookie(w, &http.Cookie{
		Name:     "session_id",
		Value:    "",
		HttpOnly: true,
		Secure:   r.TLS != nil,
		Path:     "/",
		Expires:  time.Unix(0, 0), // Date dans le passé → suppression immédiate
		MaxAge:   -1,
		SameSite: http.SameSiteLaxMode,
	})

	http.Redirect(w, r, "/login", http.StatusSeeOther)
}

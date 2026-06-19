package handlers

import (
	"bytes"
	"database/sql"
	"html/template"
	"net/http"
	"strings"
	"time"

	"forum/middleware"
	"forum/models"
	"forum/utils"

	"github.com/google/uuid"
)

func (h *AuthHandler) Login(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/login.html",
	))

	switch r.Method {
	case http.MethodGet:
		next := r.URL.Query().Get("next")
		data := models.TemplateData{
			IsLoggedIn: middleware.GetUserIDFromCookie(w, h.db, r) != 0,
			Next:       next,
		}

		if r.URL.Query().Get("registered") == "1" {
			data.Success = "Inscription réussie ! Connectez-vous."
		}

		var buf bytes.Buffer
		if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
			http.Error(w, "Erreur rendu template", http.StatusInternalServerError)
			return
		}
		buf.WriteTo(w)

	case http.MethodPost:
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		if err := utils.ValidateLogin(email, password); err != nil {
			data := models.TemplateData{Error: err.Error()}
			var buf bytes.Buffer
			if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
				http.Error(w, "Erreur rendu template", http.StatusInternalServerError)
				return
			}
			buf.WriteTo(w)
			return
		}

		id, _, hashedPassword, err := getUserByEmail(h.db, email)
		if err == sql.ErrNoRows {
			data := models.TemplateData{Error: "Email ou mot de passe incorrect"}
			var buf bytes.Buffer
			if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
				http.Error(w, "Erreur rendu template", http.StatusInternalServerError)
				return
			}
			buf.WriteTo(w)
			return
		}
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		if err := utils.CheckPassword(hashedPassword, password); err != nil {
			data := models.TemplateData{Error: "Email ou mot de passe incorrect"}
			var buf bytes.Buffer
			if err := tmpl.ExecuteTemplate(&buf, "base", data); err != nil {
				http.Error(w, "Erreur rendu template", http.StatusInternalServerError)
				return
			}
			buf.WriteTo(w)
			return
		}

		next := strings.TrimSpace(r.FormValue("next"))
		if next == "" || !strings.HasPrefix(next, "/") {
			next = "/"
		}

		// Création d'un UUID pour la session
		sessionID := uuid.New().String()
		expiresAt := time.Now().Add(24 * time.Hour)

		_, err = h.db.Exec(
			"INSERT INTO sessions (id, user_id, expires_at) VALUES (?, ?, ?)",
			sessionID, id, expiresAt,
		)
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		http.SetCookie(w, &http.Cookie{
			Name:     "session_id",
			Value:    sessionID,
			HttpOnly: true,
			Secure:   r.TLS != nil,
			Path:     "/",
			Expires:  expiresAt,
			MaxAge:   24 * 60 * 60,
			SameSite: http.SameSiteLaxMode,
		})

		http.Redirect(w, r, next, http.StatusSeeOther)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

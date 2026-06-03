package handlers

import (
	"database/sql"
	"html/template"
	"net/http"

	"forum/middleware"
	"forum/models"
)

// AuthHandler regroupe tous les handlers d'authentification.
type AuthHandler struct {
	db *sql.DB
}

// NewAuthHandler crée un AuthHandler avec la connexion DB fournie.
func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

// Profile affiche la page de profil de l'utilisateur connecté.
func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/profile.html",
	))

	switch r.Method {
	case http.MethodGet:
		userID := middleware.GetUserID(r)
		if userID == 0 {
			http.Redirect(w, r, "/login", http.StatusSeeOther)
			return
		}

		user, err := getUserByID(h.db, userID)
		if err == sql.ErrNoRows {
			http.Error(w, "Utilisateur introuvable", http.StatusNotFound)
			return
		}
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "base", map[string]any{
			"User":       user,
			"IsLoggedIn": true,
		})

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func getUserByID(db *sql.DB, id int) (models.User, error) {
	var user models.User
	err := db.QueryRow(
		"SELECT username, email FROM users WHERE id = ?",
		id,
	).Scan(&user.Username, &user.Email)
	return user, err
}

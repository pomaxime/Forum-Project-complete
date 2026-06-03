package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"

	"forum/utils"
)

// RegisterHandler gère l'inscription des nouveaux utilisateurs.
func (h *AuthHandler) Register(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/register.html",
	))

	switch r.Method {
	case http.MethodGet:
		tmpl.ExecuteTemplate(w, "base", nil)

	case http.MethodPost:
		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(r.FormValue("email"))
		password := r.FormValue("password")

		// Validation des champs
		if err := utils.ValidateRegister(username, email, password); err != nil {
			tmpl.ExecuteTemplate(w, "base", map[string]string{
				"Error":    err.Error(),
				"Username": username,
				"Email":    email,
				"Password": password,
			})
			return
		}

		// Hash du mot de passe
		hashedPassword, err := utils.HashPassword(password)
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		// Insertion en base de données (requête préparée → protection SQL Injection)
		_, err = h.db.Exec(
			"INSERT INTO users (username, email, password) VALUES (?, ?, ?)",
			username, email, hashedPassword,
		)
		if err != nil {
			// Vérifie si l'email ou le username est déjà utilisé
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.email") {
				tmpl.ExecuteTemplate(w, "base", map[string]string{
					"Error":      "Cet e-mail est déjà utilisé",
					"Username":   username,
					"Email":      email,
					"Password":   password,
					"EmailClass": "input-error",
				})
				return
			}
			if strings.Contains(err.Error(), "UNIQUE constraint failed: users.username") {
				tmpl.ExecuteTemplate(w, "base", map[string]string{
					"Error":         "Ce nom d'utilisateur est déjà pris",
					"Username":      username,
					"Email":         email,
					"Password":      password,
					"UsernameClass": "input-error",
				})
				return
			}
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		// Inscription réussie → redirection vers la page de login
		http.Redirect(w, r, "/login?registered=1", http.StatusSeeOther)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// getUserByEmail récupère un utilisateur par son email.
// Retourne sql.ErrNoRows si aucun utilisateur trouvé.
func getUserByEmail(db *sql.DB, email string) (int, string, string, error) {
	var id int
	var username, hashedPassword string
	err := db.QueryRow(
		"SELECT id, username, password FROM users WHERE email = ?", email,
	).Scan(&id, &username, &hashedPassword)
	return id, username, hashedPassword, err
}

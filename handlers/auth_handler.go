package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"

	"forum/middleware"
	"forum/models"
	"forum/repository"
	"forum/utils"
)

type AuthHandler struct {
	db *sql.DB
}

func NewAuthHandler(db *sql.DB) *AuthHandler {
	return &AuthHandler{db: db}
}

func (h *AuthHandler) Profile(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/profile.html",
	))

	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Redirect(w, r, "/login", http.StatusSeeOther)
		return
	}

	switch r.Method {
	case http.MethodGet:
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

	case http.MethodPost:
		if err := r.ParseForm(); err != nil {
			http.Error(w, "Données du formulaire invalides", http.StatusBadRequest)
			return
		}

		username := strings.TrimSpace(r.FormValue("username"))
		email := strings.TrimSpace(r.FormValue("email"))
		avatarURL := strings.TrimSpace(r.FormValue("avatar_url"))
		gender := strings.TrimSpace(r.FormValue("gender"))
		newPassword := r.FormValue("new_password")
		confirmPassword := r.FormValue("confirm_password")

		if username == "" || email == "" {
			tmpl.ExecuteTemplate(w, "base", map[string]any{
				"User": models.User{
					ID:        userID,
					Username:  username,
					Email:     email,
					AvatarURL: avatarURL,
					Gender:    gender,
				},
				"IsLoggedIn": true,
				"Error":      "Le nom d'utilisateur et l'email sont obligatoires",
			})
			return
		}

		if !utils.ValidateEmail(email) {
			tmpl.ExecuteTemplate(w, "base", map[string]any{
				"User": models.User{
					ID:        userID,
					Username:  username,
					Email:     email,
					AvatarURL: avatarURL,
					Gender:    gender,
				},
				"IsLoggedIn": true,
				"Error":      "Adresse e-mail invalide",
			})
			return
		}

		currentUser, err := getUserByID(h.db, userID)
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		if email != currentUser.Email {
			if exists, err := repository.NewUserRepo(h.db).EmailExists(email); err != nil {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			} else if exists {
				tmpl.ExecuteTemplate(w, "base", map[string]any{
					"User": models.User{
						ID:        userID,
						Username:  username,
						Email:     email,
						AvatarURL: avatarURL,
						Gender:    gender,
					},
					"IsLoggedIn": true,
					"Error":      "Cet e-mail est déjà utilisé",
				})
				return
			}
		}

		if username != currentUser.Username {
			if exists, err := repository.NewUserRepo(h.db).UsernameExists(username); err != nil {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			} else if exists {
				tmpl.ExecuteTemplate(w, "base", map[string]any{
					"User": models.User{
						ID:        userID,
						Username:  username,
						Email:     email,
						AvatarURL: avatarURL,
						Gender:    gender,
					},
					"IsLoggedIn": true,
					"Error":      "Ce nom d'utilisateur est déjà pris",
				})
				return
			}
		}

		if newPassword != "" {
			if newPassword != confirmPassword {
				tmpl.ExecuteTemplate(w, "base", map[string]any{
					"User": models.User{
						ID:        userID,
						Username:  username,
						Email:     email,
						AvatarURL: avatarURL,
						Gender:    gender,
					},
					"IsLoggedIn": true,
					"Error":      "Les mots de passe ne correspondent pas",
				})
				return
			}
			if len(newPassword) < 8 {
				tmpl.ExecuteTemplate(w, "base", map[string]any{
					"User": models.User{
						ID:        userID,
						Username:  username,
						Email:     email,
						AvatarURL: avatarURL,
						Gender:    gender,
					},
					"IsLoggedIn": true,
					"Error":      "Le mot de passe doit contenir au moins 8 caractères",
				})
				return
			}
		}

		r := repository.NewUserRepo(h.db)
		if err := r.UpdateProfile(userID, username, email, avatarURL, gender); err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		if newPassword != "" {
			hashedPassword, err := utils.HashPassword(newPassword)
			if err != nil {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
			if err := r.UpdatePassword(userID, hashedPassword); err != nil {
				http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
				return
			}
		}

		user, err := getUserByID(h.db, userID)
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		tmpl.ExecuteTemplate(w, "base", map[string]any{
			"User":       user,
			"IsLoggedIn": true,
			"Success":    "Profil mis à jour avec succès",
		})

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func getUserByID(db *sql.DB, id int) (models.User, error) {
	var user models.User
	err := db.QueryRow(
		"SELECT id, username, email, avatar_url, gender FROM users WHERE id = ?",
		id,
	).Scan(&user.ID, &user.Username, &user.Email, &user.AvatarURL, &user.Gender)
	return user, err
}

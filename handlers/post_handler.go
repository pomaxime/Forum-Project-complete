package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"
	"time"

	"forum/middleware"
	"forum/utils"
)

// PostHandler gère les opérations sur les posts.
type PostHandler struct {
	db *sql.DB
}

// NewPostHandler crée un PostHandler avec la connexion DB fournie.
func NewPostHandler(db *sql.DB) *PostHandler {
	return &PostHandler{db: db}
}

// Post représente un post du forum.
type Post struct {
	ID        int
	Username  string
	Title     string
	Content   string
	CreatedAt time.Time
}

// Index affiche la liste de tous les posts (page d'accueil).
func (h *PostHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles("templates/index.html"))

	rows, err := h.db.Query(`
		SELECT p.id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		ORDER BY p.created_at DESC
	`)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Username, &p.Title, &p.Content, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	// Vérifie si l'utilisateur est connecté (pour afficher les boutons Create/Logout)
	userID := middleware.GetUserIDFromCookie(w, h.db, r)
	tmpl.Execute(w, map[string]interface{}{
		"Posts":    posts,
		"LoggedIn": userID != 0,
	})
}

// Create gère la création d'un nouveau post (réservé aux utilisateurs connectés).
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/create_post.html"))

	switch r.Method {
	case http.MethodGet:
		tmpl.Execute(w, nil)

	case http.MethodPost:
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))

		if err := utils.ValidatePost(title, content); err != nil {
			tmpl.Execute(w, map[string]string{"Error": err.Error()})
			return
		}

		// Récupération de l'utilisateur connecté via le contexte
		userID := middleware.GetUserID(r)

		_, err := h.db.Exec(
			"INSERT INTO posts (user_id, title, content) VALUES (?, ?, ?)",
			userID, title, content,
		)
		if err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

// Show affiche un post individuel.
func (h *PostHandler) Show(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/post.html"))

	// Extraction de l'ID depuis l'URL : /post/42
	id := strings.TrimPrefix(r.URL.Path, "/post/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	var p Post
	err := h.db.QueryRow(`
		SELECT p.id, u.username, p.title, p.content, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(&p.ID, &p.Username, &p.Title, &p.Content, &p.CreatedAt)

	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	data := struct {
		Post       Post
		IsLoggedIn bool
	}{
		Post:       p,
		IsLoggedIn: middleware.GetUserIDFromCookie(w, h.db, r) != 0,
	}

	tmpl.Execute(w, data)
}

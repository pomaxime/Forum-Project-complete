package handlers

import (
	"database/sql"
	"html/template"
	"net/http"
	"strings"
	"time"

	"forum/middleware"
	"forum/models"
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
	Category  string
	CreatedAt time.Time
}

type CategoryOption struct {
	Value string
	Label string
}

type CreatePostData struct {
	models.TemplateData
	Title         string
	Category      string
	Content       string
	TitleClass    string
	CategoryClass string
	ContentClass  string
	Categories    []CategoryOption
}

type IndexData struct {
	models.TemplateData
	Posts            []Post
	SelectedCategory string
}

type PostPageData struct {
	models.TemplateData
	Post     Post
	Comments []Comment
}

var postCategoryOptions = []CategoryOption{
	{Value: "general", Label: "Général"},
	{Value: "tech", Label: "Tech"},
	{Value: "jeux", Label: "Jeux"},
	{Value: "business", Label: "Business"},
	{Value: "écologie", Label: "Écologie"},
	{Value: "santé", Label: "Santé"},
	{Value: "sport", Label: "Sport"},
	{Value: "culture", Label: "Culture"},
	{Value: "éducation", Label: "Éducation"},
}

var allowedPostCategories = make(map[string]bool)

func init() {
	for _, option := range postCategoryOptions {
		allowedPostCategories[option.Value] = true
	}
}

// Index affiche la liste de tous les posts (page d'accueil).
func (h *PostHandler) Index(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/index.html",
	))

	category := strings.TrimSpace(r.URL.Query().Get("category"))

	query := `
		SELECT p.id, u.username, p.title, p.content, p.category, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
	`
	var rows *sql.Rows
	var err error
	if category != "" {
		if !allowedPostCategories[category] {
			http.Error(w, "Catégorie invalide", http.StatusBadRequest)
			return
		}
		query += " WHERE p.category = ?"
		query += " ORDER BY p.created_at DESC"
		rows, err = h.db.Query(query, category)
	} else {
		query += " ORDER BY p.created_at DESC"
		rows, err = h.db.Query(query)
	}
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var posts []Post
	for rows.Next() {
		var p Post
		if err := rows.Scan(&p.ID, &p.Username, &p.Title, &p.Content, &p.Category, &p.CreatedAt); err != nil {
			continue
		}
		posts = append(posts, p)
	}

	userID := middleware.GetUserIDFromCookie(w, h.db, r)
	data := IndexData{
		TemplateData:     models.TemplateData{IsLoggedIn: userID != 0},
		Posts:            posts,
		SelectedCategory: category,
	}
	tmpl.ExecuteTemplate(w, "base", data)
}

// Create gère la création d'un nouveau post (réservé aux utilisateurs connectés).
func (h *PostHandler) Create(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/create_post.html",
	))

	switch r.Method {
	case http.MethodGet:
		data := CreatePostData{
			TemplateData: models.TemplateData{IsLoggedIn: true},
			Category:     "general",
			Categories:   postCategoryOptions,
		}
		tmpl.ExecuteTemplate(w, "base", data)

	case http.MethodPost:
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		category := strings.TrimSpace(r.FormValue("category"))
		if category == "" {
			category = "general"
		}

		data := CreatePostData{
			TemplateData: models.TemplateData{IsLoggedIn: true},
			Title:        title,
			Content:      content,
			Category:     category,
			Categories:   postCategoryOptions,
		}

		if !allowedPostCategories[category] {
			data.Error = "Catégorie invalide"
			data.CategoryClass = "input-error"
			tmpl.ExecuteTemplate(w, "base", data)
			return
		}

		if err := utils.ValidatePost(title, content); err != nil {
			data.Error = err.Error()
			if strings.Contains(err.Error(), "titre") {
				data.TitleClass = "input-error"
			}
			if strings.Contains(err.Error(), "contenu") {
				data.ContentClass = "input-error"
			}
			if strings.Contains(err.Error(), "obligatoires") {
				data.TitleClass = "input-error"
				data.ContentClass = "input-error"
			}
			tmpl.ExecuteTemplate(w, "base", data)
			return
		}

		// Récupération de l'utilisateur connecté via le contexte
		userID := middleware.GetUserID(r)

		if err := h.createPost(userID, title, content, category); err != nil {
			http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
			return
		}

		http.Redirect(w, r, "/", http.StatusSeeOther)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func (h *PostHandler) createPost(userID int, title, content, category string) error {
	_, err := h.db.Exec(
		"INSERT INTO posts (user_id, title, content, category) VALUES (?, ?, ?, ?)",
		userID, title, content, category,
	)
	return err
}

// Show affiche un post individuel.
func (h *PostHandler) Show(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/post.html",
	))

	// Extraction de l'ID depuis l'URL : /post/42
	id := strings.TrimPrefix(r.URL.Path, "/post/")
	if id == "" {
		http.NotFound(w, r)
		return
	}

	var p Post
	err := h.db.QueryRow(`
		SELECT p.id, u.username, p.title, p.content, p.category, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(&p.ID, &p.Username, &p.Title, &p.Content, &p.Category, &p.CreatedAt)

	if err == sql.ErrNoRows {
		http.NotFound(w, r)
		return
	}
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	commentsRows, err := h.db.Query(`
		SELECT c.id, u.username, c.content, c.created_at
		FROM comments c
		JOIN users u ON c.user_id = u.id
		WHERE c.post_id = ?
		ORDER BY c.created_at ASC
	`, id)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}
	defer commentsRows.Close()

	var comments []Comment
	for commentsRows.Next() {
		var comment Comment
		if err := commentsRows.Scan(&comment.ID, &comment.Username, &comment.Content, &comment.CreatedAt); err != nil {
			continue
		}
		comments = append(comments, comment)
	}

	data := PostPageData{
		TemplateData: models.TemplateData{IsLoggedIn: middleware.GetUserIDFromCookie(w, h.db, r) != 0},
		Post:         p,
		Comments:     comments,
	}

	tmpl.ExecuteTemplate(w, "base", data)
}

func (h *PostHandler) Categories(w http.ResponseWriter, r *http.Request) {
	tmpl := template.Must(template.ParseFiles("templates/categories.html"))
	tmpl.Execute(w, nil)
}

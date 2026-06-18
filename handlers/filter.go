package handlers

import (
	"html/template"
	"net/http"
	"strings"

	"forum/middleware"
	"forum/models"
)

type FilterHandler struct {
	*PostHandler
}

func NewFilterHandler(ph *PostHandler) *FilterHandler {
	return &FilterHandler{PostHandler: ph}
}

func (h *FilterHandler) ByUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	username := strings.TrimSpace(r.URL.Query().Get("user"))
	if username == "" {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rows, err := h.db.Query(`
		SELECT p.id, u.username, p.title, p.content, p.category, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE u.username = ?
		ORDER BY p.created_at DESC
	`, username)
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

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/index.html",
	))

	userID := middleware.GetUserIDFromCookie(w, h.db, r)
	tmpl.ExecuteTemplate(w, "base", IndexData{
		TemplateData:     models.TemplateData{IsLoggedIn: userID != 0},
		Posts:            posts,
		SelectedCategory: "",
	})
}

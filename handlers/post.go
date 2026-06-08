package handlers

import (
	"html/template"
	"net/http"
	"strconv"
	"strings"

	"forum/middleware"
	"forum/utils"
)

func (h *PostHandler) Edit(w http.ResponseWriter, r *http.Request) {
	userID := middleware.GetUserID(r)

	idStr := strings.TrimSpace(r.URL.Query().Get("id"))
	if r.Method == http.MethodPost {
		idStr = strings.TrimSpace(r.FormValue("id"))
	}
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	var ownerID int
	err = h.db.QueryRow("SELECT user_id FROM posts WHERE id = ?", id).Scan(&ownerID)
	if err != nil {
		http.NotFound(w, r)
		return
	}
	if ownerID != userID {
		renderError(w, http.StatusForbidden, "Vous n'êtes pas autorisé à modifier ce post")
		return
	}

	switch r.Method {
	case http.MethodGet:
		post, err := h.getPostByID(id)
		if err != nil {
			http.NotFound(w, r)
			return
		}

		tmpl, err := parseTemplates("templates/edit_post.html")
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Erreur interne du serveur")
			return
		}
		tmpl.ExecuteTemplate(w, "base", map[string]any{
			"Post":       post,
			"IsLoggedIn": true,
		})

	case http.MethodPost:
		title := strings.TrimSpace(r.FormValue("title"))
		content := strings.TrimSpace(r.FormValue("content"))
		category := strings.TrimSpace(r.FormValue("category"))
		if category == "" {
			category = "general"
		}

		if err := utils.ValidatePost(title, content); err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}
		if !allowedPostCategories[category] {
			http.Error(w, "Catégorie invalide", http.StatusBadRequest)
			return
		}

		_, err = h.db.Exec(
			"UPDATE posts SET title = ?, content = ?, category = ? WHERE id = ? AND user_id = ?",
			title, content, category, id, userID,
		)
		if err != nil {
			renderError(w, http.StatusInternalServerError, "Erreur interne du serveur")
			return
		}
		http.Redirect(w, r, "/post/"+strconv.Itoa(id), http.StatusSeeOther)

	default:
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
	}
}

func (h *PostHandler) Delete(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	userID := middleware.GetUserID(r)
	id, err := strconv.Atoi(strings.TrimSpace(r.FormValue("id")))
	if err != nil || id <= 0 {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	res, err := h.db.Exec(
		"DELETE FROM posts WHERE id = ? AND user_id = ?", id, userID,
	)
	if err != nil {
		renderError(w, http.StatusInternalServerError, "Erreur interne du serveur")
		return
	}
	n, _ := res.RowsAffected()
	if n == 0 {
		renderError(w, http.StatusForbidden, "Post introuvable ou accès refusé")
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}

func (h *PostHandler) getPostByID(id int) (*Post, error) {
	var p Post
	err := h.db.QueryRow(`
		SELECT p.id, u.username, p.title, p.content, p.category, p.created_at
		FROM posts p
		JOIN users u ON p.user_id = u.id
		WHERE p.id = ?
	`, id).Scan(&p.ID, &p.Username, &p.Title, &p.Content, &p.Category, &p.CreatedAt)
	if err != nil {
		return nil, err
	}
	return &p, nil
}

func parseTemplates(contentTemplate string) (*template.Template, error) {
	return template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		contentTemplate,
	)
}

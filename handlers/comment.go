package handlers

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"forum/middleware"
)

type Comment struct {
	ID        int
	Username  string
	Content   string
	CreatedAt time.Time
}

func (h *PostHandler) Comment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	postID, err := strconv.Atoi(r.FormValue("post_id"))
	if err != nil || postID <= 0 {
		http.Error(w, "ID de post invalide", http.StatusBadRequest)
		return
	}

	content := strings.TrimSpace(r.FormValue("content"))
	if len(content) < 3 {
		http.Error(w, "Le commentaire doit contenir au moins 3 caractères", http.StatusBadRequest)
		return
	}

	userID := middleware.GetUserID(r)
	if userID == 0 {
		http.Error(w, "Vous devez être connecté pour commenter", http.StatusUnauthorized)
		return
	}

	_, err = h.db.Exec(
		"INSERT INTO comments (post_id, user_id, content) VALUES (?, ?, ?)",
		postID, userID, content,
	)
	if err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/post/"+strconv.Itoa(postID), http.StatusSeeOther)
}
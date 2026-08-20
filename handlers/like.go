package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
	"strconv"

	"forum/middleware"
	"forum/repository"
)

type LikeHandler struct {
	repo *repository.LikeRepository
}

func NewLikeHandler(db *sql.DB) *LikeHandler {
	return &LikeHandler{repo: repository.NewLikeRepository(db)}
}

func (h *LikeHandler) Like(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Données du formulaire invalides", http.StatusBadRequest)
		return
	}

	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	userID := middleware.GetUserID(r)
	if postID == 0 || userID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if authorID, err := h.repo.GetPostAuthorID(postID); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	} else if authorID == userID {
		http.Error(w, "Vous ne pouvez pas réagir à votre propre post", http.StatusForbidden)
		return
	}

	if err := h.repo.Like(postID, userID); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	redirectURL := r.Referer()
	if redirectURL == "" {
		redirectURL = fmt.Sprintf("/post/%d", postID)
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *LikeHandler) Dislike(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Données du formulaire invalides", http.StatusBadRequest)
		return
	}

	postID, _ := strconv.Atoi(r.FormValue("post_id"))
	userID := middleware.GetUserID(r)
	if postID == 0 || userID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if authorID, err := h.repo.GetPostAuthorID(postID); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	} else if authorID == userID {
		http.Error(w, "Vous ne pouvez pas réagir à votre propre post", http.StatusForbidden)
		return
	}

	if err := h.repo.Dislike(postID, userID); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	redirectURL := r.Referer()
	if redirectURL == "" {
		redirectURL = fmt.Sprintf("/post/%d", postID)
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *LikeHandler) LikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Données du formulaire invalides", http.StatusBadRequest)
		return
	}

	commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
	userID := middleware.GetUserID(r)
	if commentID == 0 || userID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if authorID, err := h.repo.GetCommentAuthorID(commentID); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	} else if authorID == userID {
		http.Error(w, "Vous ne pouvez pas réagir à votre propre commentaire", http.StatusForbidden)
		return
	}

	if err := h.repo.LikeComment(commentID, userID); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	redirectURL := r.Referer()
	if redirectURL == "" {
		if postID, err := h.repo.GetCommentPostID(commentID); err == nil {
			redirectURL = fmt.Sprintf("/post/%d", postID)
		} else {
			redirectURL = "/"
		}
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

func (h *LikeHandler) DislikeComment(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	if err := r.ParseForm(); err != nil {
		http.Error(w, "Données du formulaire invalides", http.StatusBadRequest)
		return
	}

	commentID, _ := strconv.Atoi(r.FormValue("comment_id"))
	userID := middleware.GetUserID(r)
	if commentID == 0 || userID == 0 {
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if authorID, err := h.repo.GetCommentAuthorID(commentID); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	} else if authorID == userID {
		http.Error(w, "Vous ne pouvez pas réagir à votre propre commentaire", http.StatusForbidden)
		return
	}

	if err := h.repo.DislikeComment(commentID, userID); err != nil {
		http.Error(w, "Erreur interne du serveur", http.StatusInternalServerError)
		return
	}

	redirectURL := r.Referer()
	if redirectURL == "" {
		if postID, err := h.repo.GetCommentPostID(commentID); err == nil {
			redirectURL = fmt.Sprintf("/post/%d", postID)
		} else {
			redirectURL = "/"
		}
	}
	http.Redirect(w, r, redirectURL, http.StatusSeeOther)
}

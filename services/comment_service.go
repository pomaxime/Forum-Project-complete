package services

import (
	"database/sql"
	"errors"
	"strings"

	"forum/models"
	"forum/repository"
)

type CommentService struct {
	comments *repository.CommentRepo
}

func NewCommentService(db *sql.DB) *CommentService {
	return &CommentService{comments: repository.NewCommentRepo(db)}
}

var ErrCommentTooShort = errors.New("le commentaire doit contenir au moins 3 caractères")

func (s *CommentService) GetByPostID(postID int) ([]models.Comment, error) {
	return s.comments.GetByPostID(postID)
}

func (s *CommentService) Create(postID, userID int, content string) error {
	content = strings.TrimSpace(content)
	if len(content) < 3 {
		return ErrCommentTooShort
	}
	_, err := s.comments.Create(postID, userID, content)
	return err
}

package services

import (
	"database/sql"
	"errors"

	"forum/models"
	"forum/repository"
	"forum/utils"
)

type PostService struct {
	posts      *repository.PostRepo
	categories *repository.CategoryRepo
}

func NewPostService(db *sql.DB) *PostService {
	return &PostService{
		posts:      repository.NewPostRepo(db),
		categories: repository.NewCategoryRepo(),
	}
}

var ErrInvalidCategory = errors.New("catégorie invalide")

func (s *PostService) GetAll(category string) ([]models.Post, error) {
	if category != "" && !s.categories.IsValid(category) {
		return nil, ErrInvalidCategory
	}
	return s.posts.GetAll(category)
}

func (s *PostService) GetByID(id int) (*models.Post, error) {
	return s.posts.GetByID(id)
}

func (s *PostService) Create(userID int, title, content, category string) error {
	if err := utils.ValidatePost(title, content); err != nil {
		return err
	}
	if !s.categories.IsValid(category) {
		return ErrInvalidCategory
	}
	_, err := s.posts.Create(userID, title, content, category)
	return err
}

func (s *PostService) GetCategories() []models.Category {
	return s.categories.GetAll()
}

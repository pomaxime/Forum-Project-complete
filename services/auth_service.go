package services

import (
	"database/sql"
	"errors"
	"strings"

	"forum/models"
	"forum/repository"
	"forum/utils"
)

type AuthService struct {
	users *repository.UserRepo
}

func NewAuthService(db *sql.DB) *AuthService {
	return &AuthService{users: repository.NewUserRepo(db)}
}

var ErrEmailTaken = errors.New("cet e-mail est déjà utilisé")

var ErrUsernameTaken = errors.New("ce nom d'utilisateur est déjà pris")

var ErrInvalidCredentials = errors.New("email ou mot de passe incorrect")

func (s *AuthService) Register(username, email, password string) error {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if err := utils.ValidateRegister(username, email, password); err != nil {
		return err
	}

	emailExists, err := s.users.EmailExists(email)
	if err != nil {
		return err
	}
	if emailExists {
		return ErrEmailTaken
	}

	usernameExists, err := s.users.UsernameExists(username)
	if err != nil {
		return err
	}
	if usernameExists {
		return ErrUsernameTaken
	}

	hashedPassword, err := utils.HashPassword(password)
	if err != nil {
		return err
	}

	_, err = s.users.Create(username, email, hashedPassword)
	return err
}

func (s *AuthService) Authenticate(email, password string) (*models.User, error) {
	user, err := s.users.GetByEmail(strings.TrimSpace(email))
	if err == sql.ErrNoRows {
		return nil, ErrInvalidCredentials
	}
	if err != nil {
		return nil, err
	}

	if err := utils.CheckPassword(user.Password, password); err != nil {
		return nil, ErrInvalidCredentials
	}

	return user, nil
}

func (s *AuthService) GetUserByID(id int) (*models.User, error) {
	return s.users.GetByID(id)
}

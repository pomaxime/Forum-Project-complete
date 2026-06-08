package utils

import (
	"errors"
	"regexp"
	"strings"
)

var emailRegex = regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)

func ValidateRegister(username, email, password string) error {
	username = strings.TrimSpace(username)
	email = strings.TrimSpace(email)

	if username == "" || email == "" || password == "" {
		return errors.New("tous les champs sont obligatoires")
	}
	if len(username) < 3 || len(username) > 30 {
		return errors.New("le nom d'utilisateur doit contenir entre 3 et 30 caractères")
	}
	if !emailRegex.MatchString(email) {
		return errors.New("adresse e-mail invalide")
	}
	if len(password) < 8 {
		return errors.New("le mot de passe doit contenir au moins 8 caractères")
	}
	return nil
}

func ValidateLogin(email, password string) error {
	if strings.TrimSpace(email) == "" || password == "" {
		return errors.New("tous les champs sont obligatoires")
	}
	return nil
}

func ValidateEmail(email string) bool {
	email = strings.TrimSpace(email)
	return emailRegex.MatchString(email)
}

func ValidatePost(title, content string) error {
	title = strings.TrimSpace(title)
	content = strings.TrimSpace(content)

	if title == "" || content == "" {
		return errors.New("le titre et le contenu sont obligatoires")
	}
	if len(title) < 3 || len(title) > 200 {
		return errors.New("le titre doit contenir entre 3 et 200 caractères")
	}
	if len(content) < 10 {
		return errors.New("le contenu doit contenir au moins 10 caractères")
	}
	return nil
}

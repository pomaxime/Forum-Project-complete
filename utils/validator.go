package utils

import (
	"errors"
	"strings"
)

func ValidateComment(content string) error {
	content = strings.TrimSpace(content)
	if content == "" {
		return errors.New("le commentaire ne peut pas être vide")
	}
	if len(content) < 3 {
		return errors.New("le commentaire doit contenir au moins 3 caractères")
	}
	if len(content) > 2000 {
		return errors.New("le commentaire ne peut pas dépasser 2000 caractères")
	}
	return nil
}

func ValidateID(id int) error {
	if id <= 0 {
		return errors.New("identifiant invalide")
	}
	return nil
}
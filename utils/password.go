package utils

import (
	"fmt"

	"golang.org/x/crypto/bcrypt"
)

// Génère un hash bcrypt à partir d'un mdp.
func HashPassword(password string) (string, error) {
	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", fmt.Errorf("erreur hachage mot de passe : %w", err)
	}
	return string(hashed), nil
}

// Compare mdp avec un hash bcrypt.
func CheckPassword(hash, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
}

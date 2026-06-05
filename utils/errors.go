package utils

import "net/http"

type AppError struct {
	Code    int
	Message string
}

func (e *AppError) Error() string {
	return e.Message
}

func NewError(code int, message string) *AppError {
	return &AppError{Code: code, Message: message}
}

var (
	ErrNotFound   = NewError(http.StatusNotFound, "Page introuvable")
	ErrForbidden  = NewError(http.StatusForbidden, "Accès refusé")
	ErrInternal   = NewError(http.StatusInternalServerError, "Erreur interne du serveur")
	ErrBadRequest = NewError(http.StatusBadRequest, "Requête invalide")
)
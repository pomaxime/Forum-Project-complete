package handlers

import (
	"fmt"
	"net/http"
)

func NotFound(w http.ResponseWriter, r *http.Request) {
	renderError(w, http.StatusNotFound, "Page introuvable")
}

func MethodNotAllowed(w http.ResponseWriter, r *http.Request) {
	renderError(w, http.StatusMethodNotAllowed, "Méthode non autorisée")
}

func renderError(w http.ResponseWriter, status int, message string) {
	w.Header().Set("Content-Type", "text/plain; charset=utf-8")
	w.WriteHeader(status)
	fmt.Fprintln(w, message)
}

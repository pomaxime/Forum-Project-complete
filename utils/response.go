package utils

import (
	"html/template"
	"net/http"
)

func RenderError(w http.ResponseWriter, code int, message string) {
	tmpl, err := template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/error.html",
	)
	if err != nil {
		http.Error(w, message, code)
		return
	}
	w.WriteHeader(code)
	tmpl.ExecuteTemplate(w, "base", map[string]any{
		"Code":    code,
		"Message": message,
	})
}

func RenderNotFound(w http.ResponseWriter) {
	RenderError(w, http.StatusNotFound, "Page introuvable")
}

func RenderInternalError(w http.ResponseWriter) {
	RenderError(w, http.StatusInternalServerError, "Erreur interne du serveur")
}
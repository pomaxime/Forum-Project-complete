package handlers

import (
	"html/template"
	"net/http"

	"forum/models"
)

type HomeHandler struct{}

func NewHomeHandler() *HomeHandler {
	return &HomeHandler{}
}

type homeData struct {
	models.TemplateData
	AppName string
}

func (h *HomeHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	if r.Method != http.MethodGet {
		http.Error(w, "Méthode non autorisée", http.StatusMethodNotAllowed)
		return
	}

	tmpl := template.Must(template.ParseFiles(
		"templates/base.html",
		"templates/partials/navbar.html",
		"templates/partials/alerts.html",
		"templates/partials/footer.html",
		"templates/home.html",
	))

	tmpl.ExecuteTemplate(w, "base", homeData{
		TemplateData: models.TemplateData{IsLoggedIn: false},
		AppName:      "Reddick",
	})
}

package handlers

import (
	"log"
	"net/http"

	"github.com/jordanhw34/ambershouse/pkg/config"
	"github.com/jordanhw34/ambershouse/pkg/models"
	"github.com/jordanhw34/ambershouse/pkg/render"
)

// Repo => variable representing the Repository type
var Repo *Repository

// Repository => the Repository type
type Repository struct {
	App *config.AppConfig
}

// NewRepo => Creates a new Repository
func NewRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
	}
}

// NewHandlers => sets the repository for the handlers
func NewHandlers(r *Repository) {
	Repo = r
}

// Home is the home page handler
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// store remote IP address
	remoteIP := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.RenderTemplate(w, "home.page.html", &models.TemplateData{})
}

// About is the about page handler
func (repo *Repository) About(w http.ResponseWriter, r *http.Request) {
	// perform some busines logic
	stringMap := make(map[string]string)
	stringMap["test"] = "Data coming in via a handler bitch"
	stringMap["pageHeader"] = "This is the mother fucking About Page Mother Fuckers"

	// Pull the session info that was captured on the home page
	remoteIP := repo.App.Session.GetString(r.Context(), "remote_ip")
	log.Println("remote ip =", remoteIP)
	stringMap["remote_ip"] = remoteIP
	// Handling this logic in the template itself
	// if len(remoteIP) > 0 {
	// 	stringMap["remote_ip"] = remoteIP
	// } else {
	// 	stringMap["remote_ip"] = "we do not have the remote ip address captured in a cookie/session"
	// }

	render.RenderTemplate(w, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) Other(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, "other.page.html", &models.TemplateData{})
}

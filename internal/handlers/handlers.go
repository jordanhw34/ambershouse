package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/models"
	"github.com/jordanhw34/ambershouse/internal/render"
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

	render.RenderTemplate(w, r, "home.page.html", &models.TemplateData{})
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

	render.RenderTemplate(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) Bilbo(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "bilbo.page.html", &models.TemplateData{})
}

// Renders the Frodo room page
func (repo *Repository) Frodo(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "frodo.page.html", &models.TemplateData{})
}

// Renders the Reservations page
func (repo *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "reservations.page.html", &models.TemplateData{})
}

// Renders the PostReservations Handler
func (repo *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	// when pulling data out of forms it is always a string, often will need to cast to something else
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	log.Printf(" > Start Date: %s  - - -  End Date: %s", start, end)
	w.Write([]byte(fmt.Sprintf("Posted Form Data: start = %s -- end = %s", start, end)))
}

// jsonResponse struct => intentionally lower case since only used here
type jsonResponse struct {
	OK      bool   `json:"ok"`
	Message string `json:"message"`
}

// Renders the PostReservationsRoom to handle post request from room pages
func (repo *Repository) PostReservationsRoom(w http.ResponseWriter, r *http.Request) {
	resp := jsonResponse{
		OK:      true,
		Message: "Available!",
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		log.Println("Error Message:", err)
	}

	// Create a header so the browser knows what type of data it is receiving
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Renders the Reservations page
func (repo *Repository) Confirm(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "confirm.page.html", &models.TemplateData{})
}

// Renders the Contact page
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.RenderTemplate(w, r, "contact.page.html", &models.TemplateData{})
}

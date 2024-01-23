package handlers

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/drivers"
	"github.com/jordanhw34/ambershouse/internal/forms"
	"github.com/jordanhw34/ambershouse/internal/helpers"
	"github.com/jordanhw34/ambershouse/internal/models"
	"github.com/jordanhw34/ambershouse/internal/render"
	"github.com/jordanhw34/ambershouse/internal/repository"
	"github.com/jordanhw34/ambershouse/internal/repository/dbrepo"
)

// Repository => the Repository type
type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

// Repo => variable representing the Repository type
var Repo *Repository

// NewRepo => Creates a new Repository
func NewRepo(app *config.AppConfig, db *drivers.DB) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewPostgresRepo(db.SQL, app),
	}
}

// NewHandlers => sets the repository for the handlers
func NewHandlers(repo *Repository) {
	Repo = repo
}

// Home is the home page handler
func (repo *Repository) Home(w http.ResponseWriter, r *http.Request) {
	// This is just demonstrating how to access the methods defined in our DB repo
	isAllUsers := repo.DB.AllUsers()
	log.Println("is AllUsers()?", isAllUsers)

	// store remote IP address
	remoteIP := r.RemoteAddr
	repo.App.Session.Put(r.Context(), "remote_ip", remoteIP)

	render.Template(w, r, "home.page.html", &models.TemplateData{})
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

	render.Template(w, r, "about.page.html", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (repo *Repository) Bilbo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "bilbo.page.html", &models.TemplateData{})
}

// Renders the Frodo room page
func (repo *Repository) Frodo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "frodo.page.html", &models.TemplateData{})
}

// Renders the Reservations page
func (repo *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "reservations.page.html", &models.TemplateData{})
}

// Renders the PostReservations Handler
func (repo *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	// when pulling data out of forms it is always a string, often will need to cast to something else
	start := r.Form.Get("start")
	end := r.Form.Get("end")
	log.Println("Form Data - Start Date:", start)
	log.Println("Form Data - End Date:", end)

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, start)
	if err != nil {
		log.Println(" > Error in Parsing startDate")
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, end)
	if err != nil {
		log.Println(" > Error in Parsing endDate")
		helpers.ServerError(w, err)
		return
	}

	log.Println("Parsed - Start Date:", startDate)
	log.Println("Parsed - End Date:", endDate)

	rooms, err := repo.DB.AvailableByDates(startDate, endDate)
	if err != nil {
		log.Println(" > Error in AvailableByDates() Function")
		helpers.ServerError(w, err)
		return
	}

	if len(rooms) == 0 {
		repo.App.Session.Put(r.Context(), "error", "No availability in those dates")
		http.Redirect(w, r, "/reservations", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: startDate,
		EndDate:   endDate,
	}

	// Putting the Reservation prototype we created above into the Session so can use it elsewhere, don't have to pass it
	repo.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.html", &models.TemplateData{
		Data: data,
	})
}

func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the Res out of the Session that we populated in the ReservationsPost() Handler
	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from session in ChooseRoom() Handler Function"))
		return
	}

	// Update the Reservation with the RoomID that was selected and pulled out of the URL Params
	res.RoomID = roomID

	// Put the Res back into the Session for the Confirm handler
	repo.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/confirm", http.StatusSeeOther)
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
		//log.Println("Error Message:", err)
		helpers.ServerError(w, err)
	}

	// Create a header so the browser knows what type of data it is receiving
	log.Println(string(out))
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)
}

// Confirm renders the Reservations page
func (repo *Repository) Confirm(w http.ResponseWriter, r *http.Request) {
	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		helpers.ServerError(w, errors.New("cannot get reservation from sesssion in Confirm() Handler Function"))
		return
	}

	layout := "2006-01-02"

	startDate := res.StartDate.Format(layout)
	endDate := res.EndDate.Format(layout)

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "confirm.page.html", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// Handles the Form Post to confirm reservation
func (repo *Repository) PostConfirm(w http.ResponseWriter, r *http.Request) {
	log.Println(" > PostConfirm Handler")
	err := r.ParseForm()
	if err != nil {
		//log.Println("PostConfirm Error Message:", err.Error())
		helpers.ServerError(w, err)
		return
	}

	sd := r.Form.Get("start_date")
	ed := r.Form.Get("end_date")
	room_id := r.Form.Get("room_id")

	log.Printf(" > start_date = %s", sd)
	log.Printf(" > end_date = %s", ed)
	log.Printf(" > room_id = %s", room_id)

	layout := "2006-01-02"

	startDate, err := time.Parse(layout, sd)
	if err != nil {
		log.Println(" > Error in Parsing startDate")
		helpers.ServerError(w, err)
		return
	}

	endDate, err := time.Parse(layout, ed)
	if err != nil {
		log.Println(" > Error in Parsing endDate")
		helpers.ServerError(w, err)
		return
	}

	roomID, err := strconv.Atoi(room_id)
	if err != nil {
		log.Println(" > Error in Converting String to Int for Room ID")
		helpers.ServerError(w, err)
		return
	}

	reservation := models.Reservation{
		FirstName: r.Form.Get("first_name"),
		LastName:  r.Form.Get("last_name"),
		Email:     r.Form.Get("email"),
		Phone:     r.Form.Get("phone"),
		StartDate: startDate,
		EndDate:   endDate,
		RoomID:    roomID,
	}

	form := forms.New(r.PostForm)

	log.Printf("FirstName = %s \n", reservation.FirstName)
	log.Printf("LastName = %s \n", reservation.LastName)
	log.Printf("Phone = %s \n", reservation.Phone)
	log.Printf("Email = %s \n", reservation.Email)

	// If this field is blank it will add an error message that it cannot be blank
	// Instead of calling the Has() method for each field, we created a new method
	//form.Has("first_name", r)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.IsValid() {
		formData := make(map[string]interface{})
		formData["reservation"] = reservation

		render.Template(w, r, "confirm.page.html", &models.TemplateData{
			Form: form,
			Data: formData,
		})

		return
	}

	// WWrite Reservation to DB
	newResID, err := repo.DB.InsertReservation(reservation)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Before we redirect we need to populate a room_restriction record because these dates are no longer available
	restriction := models.RoomRestriction{
		StartDate:     startDate,
		EndDate:       endDate,
		RoomID:        roomID,
		ReservationID: newResID,
		RestrictionID: 1,
	}

	err = repo.DB.InsertRoomRestriction(restriction)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	repo.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/summary", http.StatusSeeOther)
}

// Renders the Contact page
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

func (repo *Repository) Summary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.ErrorLog.Println("cannot get error from session")
		log.Println("Cannot get item from session")
		repo.App.Session.Put(r.Context(), "error", "cannot get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Remove data from Session after we use it, I noticed after leaving the summary page and going back my name was still there
	repo.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	render.Template(w, r, "summary.page.html", &models.TemplateData{
		Data: data,
	})
}

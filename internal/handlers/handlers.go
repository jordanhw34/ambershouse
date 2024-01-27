package handlers

import (
	"encoding/json"
	"fmt"
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

// NewTestingRepo creates a repository for unit testing
func NewTestingRepo(app *config.AppConfig) *Repository {
	return &Repository{
		App: app,
		DB:  dbrepo.NewDBTestRepo(app),
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

// Bilbo renders the room page for Bilbo's Room
func (repo *Repository) Bilbo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "bilbo.page.html", &models.TemplateData{})
}

// Frodo renders the room page for Frodo's Room
func (repo *Repository) Frodo(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "frodo.page.html", &models.TemplateData{})
}

// Renders the Reservations page
func (repo *Repository) Reservations(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "reservations.page.html", &models.TemplateData{})
}

// Renders the PostReservations Handler
func (repo *Repository) PostReservations(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not ParseForm() in PostReservations() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	startDate, err := time.Parse(repo.App.DateFormat, start)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not parse startDate in PostReservations() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(repo.App.DateFormat, end)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not parse endDate in PostReservations() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

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

// ChooseRoom is used to choose a room after using the BookNow link
func (repo *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	roomID, err := strconv.Atoi(chi.URLParam(r, "id"))
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	// Get the Res out of the Session that we populated in the ReservationsPost() Handler
	res, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "Could not get reservation from Session in ChooseRoom() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
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
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// Renders the PostReservationsRoom to handle post request from room pages
func (repo *Repository) PostReservationsRoom(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not ParseForm() in PostReservationsRoom() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")
	room_id := r.Form.Get("room_id")
	roomID, err := strconv.Atoi(room_id)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not get RoomID in PostReservationsRoom() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	startDate, err := time.Parse(repo.App.DateFormat, start)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not parse startDate in PostReservationsRoom() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	endDate, err := time.Parse(repo.App.DateFormat, end)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not parse endDate in PostReservationsRoom() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	available, err := repo.DB.AvailableByRoomIDAndDates(startDate, endDate, roomID)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not get availability bool in PostReservationsRoom() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	resp := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: start,
		EndDate:   end,
		RoomID:    room_id,
	}

	out, err := json.MarshalIndent(resp, "", "     ")
	if err != nil {
		//log.Println("Error Message:", err)
		helpers.ServerError(w, err)
		return // added this myself but I think it's right
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
		repo.App.Session.Put(r.Context(), "error", "Could not get reservation from Session in Confirm() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := repo.DB.GetRoomByID(res.RoomID)
	if err != nil {
		repo.App.Session.Put(r.Context(), "error", "Could not find Room in database in Confirm() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	res.Room.RoomName = room.RoomName

	repo.App.Session.Put(r.Context(), "reservation", res)

	startDate := res.StartDate.Format(repo.App.DateFormat)
	endDate := res.EndDate.Format(repo.App.DateFormat)

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
	err := r.ParseForm()
	if err != nil {
		log.Println("PostConfirm() Function - ParseForm()")
		repo.App.Session.Put(r.Context(), "error", "Could not ParseForm() in PostConfirm() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		log.Println("PostConfirm() Function - Get reservation from context")
		repo.App.Session.Put(r.Context(), "error", "Could not get reservation from Session in PostConfirm() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	form.Required("first_name", "last_name", "email")
	form.MinLength("first_name", 3)
	form.IsEmail("email")

	if !form.IsValid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		render.Template(w, r, "confirm.page.html", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	// WWrite Reservation to DB
	newResID, err := repo.DB.InsertReservation(reservation)
	if err != nil {
		log.Println("PostConfirm() Function - Could not insert reservation")
		repo.App.Session.Put(r.Context(), "error", "Could not insert new reservation into DB in PostConfirm() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Before we redirect we need to populate a room_restriction record because these dates are no longer available
	restriction := models.RoomRestriction{
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		RoomID:        reservation.RoomID,
		ReservationID: newResID,
		RestrictionID: 1,
	}

	err = repo.DB.InsertRoomRestriction(restriction)
	if err != nil {
		log.Println("PostConfirm() Function - Could not insert room_restriction")
		repo.App.Session.Put(r.Context(), "error", "Could not insert Room Restriction in PostConfirm() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	guestBody := fmt.Sprintf(`
		<strong>Reservation Confirmation</strong><br /><br />
		Dear %s %s,<br /><br />
		This message is confirming your reservation from %s to %s.<br /><br />
		Thank you,<br /><br />
		Amber's House Bed & Breakfast

	`, reservation.FirstName, reservation.LastName, reservation.StartDate.Format(repo.App.DateFormat), reservation.EndDate.Format(repo.App.DateFormat))

	// Sendnotification to guest
	guestMsg := models.MailData{
		To:       reservation.Email,
		From:     "notify@ambershouse.com",
		Subject:  "Reservation Confirmation",
		Body:     guestBody,
		Template: "basic.html",
	}
	repo.App.MailChan <- guestMsg

	// Sendnotification to Property Owners
	ownerMsg := models.MailData{
		To:       "owners@ambershouse.com",
		From:     "notify@ambershouse.com",
		Subject:  "Reservation Confirmation",
		Body:     guestBody,
		Template: "basic.html",
	}
	repo.App.MailChan <- ownerMsg

	repo.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/summary", http.StatusSeeOther)
}

// Renders the Contact page
func (repo *Repository) Contact(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "contact.page.html", &models.TemplateData{})
}

// Summary display the Reservation summary page once booked
func (repo *Repository) Summary(w http.ResponseWriter, r *http.Request) {
	reservation, ok := repo.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {
		repo.App.Session.Put(r.Context(), "error", "Could not get reservation from Session in Summary() Function")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	// Remove data from Session after we use it, I noticed after leaving the summary page and going back my name was still there
	repo.App.Session.Remove(r.Context(), "reservation")

	data := make(map[string]interface{})
	data["reservation"] = reservation

	startDate := reservation.StartDate.Format(repo.App.DateFormat)
	endDate := reservation.EndDate.Format(repo.App.DateFormat)

	stringMap := make(map[string]string)
	stringMap["start_date"] = startDate
	stringMap["end_date"] = endDate

	render.Template(w, r, "summary.page.html", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

// ConfirmRoom handles the GET request from a room specific page to confirm a reservation
func (repo *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {
	roomID, _ := strconv.Atoi(r.URL.Query().Get("id"))
	// if err != nil {
	// 	helpers.ServerError(w, err)
	// 	return
	// }

	start := r.URL.Query().Get("s")
	log.Println("From Form, start = ", start)
	end := r.URL.Query().Get("e")
	log.Println("From Form, end = ", end)
	startDate, err := time.Parse(repo.App.DateFormat, start)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}
	endDate, err := time.Parse(repo.App.DateFormat, end)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	var res models.Reservation

	room, err := repo.DB.GetRoomByID(roomID)
	if err != nil {
		helpers.ServerError(w, err)
		return
	}

	res.RoomID = roomID
	res.StartDate = startDate
	res.EndDate = endDate
	res.Room.RoomName = room.RoomName

	repo.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/confirm", http.StatusSeeOther)
}

// LoginGet displays the login page for back-end users
func (repo *Repository) LoginGet(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "login.page.html", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// LoginPost handles logging the user in
func (repo *Repository) LoginPost(w http.ResponseWriter, r *http.Request) {
	_ = repo.App.Session.RenewToken(r.Context()) // good practice

	err := r.ParseForm()
	if err != nil {
		log.Println()
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)
	form.Required("email", "password")
	form.IsEmail("email")
	if !form.IsValid() {
		// TODO: take user back to login page and show a message that credentials not correct
		log.Println(" > Form is not valid")
		render.Template(w, r, "login.page.html", &models.TemplateData{
			Form: form,
		})
		return
	}

	id, _, err := repo.DB.Authenticate(email, password)
	if err != nil {
		log.Println("Error in LoginPost() Handler => Authenticate Function")
		repo.App.Session.Put(r.Context(), "error", "Invalid Login Credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}

	// Login the user
	repo.App.Session.Put(r.Context(), "user_id", id)
	repo.App.Session.Put(r.Context(), "flash", "Logged In Successfull")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

// Logout logs a user out by removing their user_id from the Session
func (repo *Repository) Logout(w http.ResponseWriter, r *http.Request) {
	//repo.App.Session.Remove(r.Context(), "user_id")
	err := repo.App.Session.Destroy(r.Context())
	if err != nil {
		log.Println("error destroying the Session in Logout() Handler Func", err)
	}

	err = repo.App.Session.RenewToken(r.Context())
	if err != nil {
		log.Println("error renewing Session token in Logout() Handler Func", err)
	}

	http.Redirect(w, r, "/user/login", http.StatusSeeOther)
}

// AdminDashboard renders the Admin Dashboard if a user is authenticated
func (repo *Repository) AdminDashboard(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-dashboard.page.html", &models.TemplateData{})
}

// AdminReservationsNew renders a page displaying new Reservations
func (repo *Repository) AdminReservationsNew(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-reservations-new.page.html", &models.TemplateData{})
}

// AdminReservationsAll renders a page displaying all Reservations
func (repo *Repository) AdminReservationsAll(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-reservations-all.page.html", &models.TemplateData{})
}

// AdminReservationsCalendar renders a page displaying the reservation calendar
func (repo *Repository) AdminReservationsCalendar(w http.ResponseWriter, r *http.Request) {
	render.Template(w, r, "admin-reservations-calendar.page.html", &models.TemplateData{})
}

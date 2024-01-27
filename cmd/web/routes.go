package main

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/handlers"
)

func routes(app *config.AppConfig) http.Handler {

	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/bilbo", handlers.Repo.Bilbo)
	mux.Get("/frodo", handlers.Repo.Frodo)

	mux.Get("/reservations", handlers.Repo.Reservations)      // /search-availability => handlers.Repo.Availability
	mux.Post("/reservations", handlers.Repo.PostReservations) // /search-availability => handlers.Repo.PostAvailability
	mux.Get("/choose-room/{id}", handlers.Repo.ChooseRoom)

	mux.Post("/reservations-room", handlers.Repo.PostReservationsRoom) // /search-availability-json => handlers.Repo.AvailabilityJSON

	// Confirm Reservation with form
	mux.Get("/confirm", handlers.Repo.Confirm)      // /make-reservation => handlers.Repo.Reservation
	mux.Post("/confirm", handlers.Repo.PostConfirm) // /make-reservation => handlers.Repo.PostReservation
	mux.Get("/confirm-room", handlers.Repo.BookRoom)

	// Reservation Summary
	mux.Get("/summary", handlers.Repo.Summary) // /reservation-summary => handlers.Repo.ReservationSummary

	mux.Get("/contact", handlers.Repo.Contact)

	mux.Get("/user/login", handlers.Repo.LoginGet)
	mux.Post("/user/login", handlers.Repo.LoginPost)
	mux.Get("/user/logout", handlers.Repo.Logout)

	// Serve up static content from the ./static/* directory
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	mux.Route("/admin", func(mux chi.Router) {
		mux.Use(Auth)
		mux.Get("/dashboard", handlers.Repo.AdminDashboard)
		mux.Get("/reservations/new", handlers.Repo.AdminReservationsNew)
		mux.Get("/reservations/all", handlers.Repo.AdminReservationsAll)
		mux.Get("/reservations/calendar", handlers.Repo.AdminReservationsCalendar)
	})

	return mux
}

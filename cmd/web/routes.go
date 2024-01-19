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
	mux.Use(WriteToConsole)
	mux.Use(NoSurf)
	mux.Use(SessionLoad)

	mux.Get("/", handlers.Repo.Home)
	mux.Get("/about", handlers.Repo.About)
	mux.Get("/bilbo", handlers.Repo.Bilbo)
	mux.Get("/frodo", handlers.Repo.Frodo)
	mux.Get("/reservations", handlers.Repo.Reservations)
	mux.Post("/reservations", handlers.Repo.PostReservations)
	mux.Post("/reservations-room", handlers.Repo.PostReservationsRoom) // change back to POST later

	mux.Get("/confirm", handlers.Repo.Confirm)
	mux.Get("/contact", handlers.Repo.Contact)

	// Serve up static content from the ./static/* directory
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux
}

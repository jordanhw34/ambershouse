package main

import (
	"log"
	"net/http"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/handlers"
	"github.com/jordanhw34/ambershouse/internal/render"
)

const host = "localhost:"
const port = "3000"

// Create Application Config
var app config.AppConfig

// Global Session Manager
var session *scs.SessionManager

// main is the entry point (main method)
func main() {

	// Change this to true when in product
	app.InProduction = false

	// Session Management
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // when they close the window or browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // needs to be true in PROD

	// Setting the session details in App Config so is available everyhere I might need it
	app.Session = session

	// Get Site TemplateCache from render package
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err.Error())
	}

	// If found, set template cache in Application Config
	app.TemplateCache = templateCache
	app.UseCache = true

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	// Now we need to give render access to the Application Config Variable [app] - passing in reference to app config variable
	render.NewTemplates(&app)

	if app.InProduction {
		log.Println(" > The app is in Production Mode")
	} else {
		log.Println(" > The app is in Development Mode")
	}
	log.Println(" > Attemping to start server on port", port)

	server := &http.Server{
		Addr:    host + port,
		Handler: routes(&app),
	}

	err = server.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}

}

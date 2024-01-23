package handlers

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"text/template"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/drivers"
	"github.com/jordanhw34/ambershouse/internal/models"
	"github.com/jordanhw34/ambershouse/internal/render"
	"github.com/justinas/nosurf"
)

var app config.AppConfig
var session *scs.SessionManager
var pathToTemplates = "./../../templates"

func getRoutes() http.Handler {

	gob.Register(models.Reservation{})

	app.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// error.log will eventually write to a file
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // when they close the window or browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // needs to be true in PROD

	app.Session = session

	// Connect to DB
	log.Println("Attemping to connect to DB...")
	db, err := drivers.ConnectSQL("host=localhost port=5432 dbname=connect_db user=postgres password=password")
	if err != nil {
		log.Fatal("cannot connect to DB", err)
		return nil
	}

	templateCache, err := CreateTestTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err.Error())
	}

	// If found, set template cache in Application Config
	app.TemplateCache = templateCache
	app.UseCache = true

	repo := NewRepo(&app, db)
	NewHandlers(repo)
	render.NewRenderer(&app)

	// NOW routes()
	mux := chi.NewRouter()

	mux.Use(middleware.Recoverer)
	//mux.Use(NoSurf)	// posts requests were failing with bad request, we could pass in the token in our tests but we are not testing NoSurf
	mux.Use(SessionLoad)

	mux.Get("/", Repo.Home)
	mux.Get("/about", Repo.About)
	mux.Get("/bilbo", Repo.Bilbo)
	mux.Get("/frodo", Repo.Frodo)

	mux.Get("/reservations", Repo.Reservations)
	mux.Post("/reservations", Repo.PostReservations)

	mux.Post("/reservations-room", Repo.PostReservationsRoom)

	// Confirm Reservation with form
	mux.Get("/confirm", Repo.Confirm)
	mux.Post("/confirm", Repo.PostConfirm)

	// Reservation Summary
	mux.Get("/summary", Repo.Summary)

	mux.Get("/contact", Repo.Contact)

	// Serve up static content from the ./static/* directory
	fileServer := http.FileServer(http.Dir("./static/"))
	mux.Handle("/static/*", http.StripPrefix("/static", fileServer))

	return mux

}

// NoSurf => adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {
	csrfHandler := nosurf.New(next)

	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",              // entire site
		Secure:   app.InProduction, // not running https change to true later in PROD
		SameSite: http.SameSiteLaxMode,
	})

	return csrfHandler
}

// SessionLoad => loads and saves the session on every request
func SessionLoad(next http.Handler) http.Handler {
	return session.LoadAndSave(next)
}

func CreateTestTemplateCache() (map[string]*template.Template, error) {
	//templateCache := make(map[string]*template.Template) => This syntax is the same as using the "make" keyword, creates and initializes an empty slice
	templateCache := map[string]*template.Template{}

	log.Println(" > Building the template cache")
	// look at templates folder and get all files named *.page.html
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		log.Println("no error getting files from the ./templates folder matching *.page.html :-) ")
		return templateCache, err
	}

	// Range through files that were found with filepath.Glob
	for _, page := range pages {
		//log.Println("ranging template pages, current page name =", page)		// e.x. templates\about.page.html
		// page is the full path to the template but we only want the base name
		name := filepath.Base(page)
		log.Println(" > Building Template Cache for Page =", name)

		// now we need to parse the file
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		// Now we need to look for layouts used in these templates
		layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return templateCache, err
		}

		// If at least 1 layout is found
		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = templateSet
	}

	return templateCache, nil
}

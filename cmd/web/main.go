package main

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/drivers"
	"github.com/jordanhw34/ambershouse/internal/handlers"
	"github.com/jordanhw34/ambershouse/internal/helpers"
	"github.com/jordanhw34/ambershouse/internal/models"
	"github.com/jordanhw34/ambershouse/internal/render"
)

const host = "localhost:"
const port = "3000"

// Create Application Config
var app config.AppConfig

// Global Session Manager
var session *scs.SessionManager

// Loggers
var infoLog *log.Logger
var errorLog *log.Logger

// main is the entry point (main method)
func main() {

	db, err := run()

	if err != nil {
		log.Fatal(err)
	}

	defer db.SQL.Close()

	defer close(app.MailChan)

	log.Println(" > Starting Email Listener Channel...")
	listenForMail()

	// Send Mail Capability
	// from := "me@here.com"
	// auth := smtp.PlainAuth("", from, "", "localhost")
	// err = smtp.SendMail("localhost:1025", auth, from, []string{"you@there.com"}, []byte("Hello this is the body"))
	// if err != nil {
	// 	log.Println("Error sending email", err)
	// }

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

func run() (*drivers.DB, error) {

	// What will we put into the session? we can store primitives but we need to tell it about structs we've created
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})

	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	// Set globally the format string to convert time.Time to YYYY-MM-DD
	app.DateFormat = "2006-01-02"

	// Change this to true when in product
	app.InProduction = false

	if app.InProduction {
		log.Println(" > The app is in Production Mode")
	} else {
		log.Println(" > The app is in Development Mode")
	}

	infoLog = log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = infoLog

	// error.log will eventually write to a file
	errorLog = log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	// Session Management
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // when they close the window or browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction // needs to be true in PROD

	// Setting the session details in App Config so is available everyhere I might need it
	app.Session = session

	// Connect to DB
	log.Println(" > Attemping to connect to DB...")
	db, err := drivers.ConnectSQL("host=localhost port=5432 dbname=ambershouse user=postgres password=password")
	if err != nil {
		log.Fatal("cannot connect to DB", err)
		return nil, err
	}
	log.Println(" > Connected to DB! :-) ")

	// Get Site TemplateCache from render package
	templateCache, err := render.CreateTemplateCache()
	if err != nil {
		log.Fatal("cannot create template cache:", err)
		return nil, err
	}

	// If found, set template cache in Application Config
	app.TemplateCache = templateCache
	app.UseCache = false

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	// Now we need to give render access to the Application Config Variable [app] - passing in reference to app config variable
	render.NewRenderer(&app)
	helpers.NewHelpers(&app)

	return db, nil
}

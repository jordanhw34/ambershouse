package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/alexedwards/scs/v2"
	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/models"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// What will we put into the session? we can store primitives but we need to tell it about structs we've created
	gob.Register(models.Reservation{})

	// Change this to true when in product
	testApp.InProduction = false

	infoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = infoLog

	// error.log will eventually write to a file
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	// Session Management
	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true // when they close the window or browser
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false

	// Setting the session details in App Config so is available everyhere I might need it
	testApp.Session = session

	// This makes sure the app setup in main.go is set to the testApp we just created
	app = &testApp

	os.Exit(m.Run())
}

// We are using this dummy type in place of a http.ResponseWriter because the test package does not have a way to create one
type responseWriter struct{}

func (rw *responseWriter) Header() http.Header {
	var header http.Header
	return header
}

func (rw *responseWriter) WriteHeader(i int) {

}

func (rw *responseWriter) Write(b []byte) (int, error) {
	length := len(b)
	return length, nil
}

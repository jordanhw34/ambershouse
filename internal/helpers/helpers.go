package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/jordanhw34/ambershouse/internal/config"
)

var app *config.AppConfig

// NewHelpers => sets up appConfig for Helpers
func NewHelpers(appConfig *config.AppConfig) {
	app = appConfig
}

func ClientError(w http.ResponseWriter, status int) {
	app.InfoLog.Println("client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	app.ErrorLog.Println(trace)
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

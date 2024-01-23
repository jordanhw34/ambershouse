package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// Early example of simple middleware that doesn't do anything useful
// func WriteToConsole(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		log.Println(" > Hit the page")
// 		next.ServeHTTP(w, r)
// 	})
// }

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

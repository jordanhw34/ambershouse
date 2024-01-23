package main

import (
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/jordanhw34/ambershouse/internal/config"
)

func TestRoutes(t *testing.T) {
	var app config.AppConfig

	mux := routes(&app)

	switch v := mux.(type) {
	case *chi.Mux:
		// do nothing, test passed
	default:
		t.Errorf("returned type from routes() is not *chi.Mux , but is type %T", v)
	}
}
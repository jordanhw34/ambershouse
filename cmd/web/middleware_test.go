package main

import (
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {
	var myH myHandler

	handler := NoSurf(&myH)

	// should return a handler
	switch v := handler.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("returned type from NoSurf() is not an http.Handler, but is type %T", v)
	}

}

func TestSessionLoad(t *testing.T) {
	var myH myHandler

	handler := SessionLoad(&myH)

	// should return a handler
	switch v := handler.(type) {
	case http.Handler:
		// do nothing
	default:
		t.Errorf("returned type from SessionLoad() is not an http.Handler, but is type %T", v)
	}

}

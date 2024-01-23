package main

import (
	"net/http"
	"os"
	"testing"
)

// The setup_test.go file is a special name that will execute before any other test files
// TestMain() is a special function name
func TestMain(m *testing.M) {

	os.Exit(m.Run())
}

type myHandler struct{}

func (mh *myHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {

}

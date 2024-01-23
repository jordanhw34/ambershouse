package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"
	"testing"
)

type postData struct {
	key   string
	value string
}

var testCases = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{
	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"bilbo", "/bilbo", "GET", []postData{}, http.StatusOK},
	{"frodo", "/frodo", "GET", []postData{}, http.StatusOK},
	{"reservations get", "/reservations", "GET", []postData{}, http.StatusOK},
	{"confirm", "/confirm", "GET", []postData{}, http.StatusOK},
	{"summary", "/summary", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"reservations post", "/reservations", "POST", []postData{
		{key: "start", value: "2024-06-01"},
		{key: "end", value: "2024-06-05"},
	}, http.StatusOK},
	{"reservations room post", "/reservations-room", "POST", []postData{
		{key: "start", value: "2024-06-01"},
		{key: "end", value: "2024-06-05"},
	}, http.StatusOK},
	{"confirm res post", "/confirm", "POST", []postData{
		{key: "first_name", value: "jordan"},
		{key: "last_name", value: "wilcox"},
		{key: "email", value: "jordan@me.com"},
		{key: "phone", value: "555-555-5555"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range testCases {
		if test.method == "GET" {
			response, err := testServer.Client().Get(testServer.URL + test.url)
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("for the %s test, expected %d status code but got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		} else {
			values := url.Values{}
			for _, x := range test.params {
				values.Add(x.key, x.value)
			}

			response, err := testServer.Client().PostForm(testServer.URL+test.url, values)
			if err != nil {
				t.Fatal(err)
			}
			if response.StatusCode != test.expectedStatusCode {
				t.Errorf("for the %s test, expected %d status code but got %d", test.name, test.expectedStatusCode, response.StatusCode)
			}
		}
	}
}

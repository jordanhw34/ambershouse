package handlers

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"

	"github.com/jordanhw34/ambershouse/internal/models"
)

// Not using anymore apparently after we last updated the tests
// type postData struct {
// 	key   string
// 	value string
// }

var testCases = []struct {
	name               string
	url                string
	method             string
	expectedStatusCode int
}{
	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"bilbo", "/bilbo", "GET", http.StatusOK},
	{"frodo", "/frodo", "GET", http.StatusOK},
	{"reservations get", "/reservations", "GET", http.StatusOK},
	// {"confirm", "/confirm", "GET", []postData{}, http.StatusOK},
	{"summary", "/summary", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	// {"reservations post", "/reservations", "POST", []postData{
	// 	{key: "start", value: "2024-06-01"},
	// 	{key: "end", value: "2024-06-05"},
	// }, http.StatusOK},
	// {"reservations room post", "/reservations-room", "POST", []postData{
	// 	{key: "start", value: "2024-06-01"},
	// 	{key: "end", value: "2024-06-05"},
	// }, http.StatusOK},
	// {"confirm res post", "/confirm", "POST", []postData{
	// 	{key: "first_name", value: "jordan"},
	// 	{key: "last_name", value: "wilcox"},
	// 	{key: "email", value: "jordan@me.com"},
	// 	{key: "phone", value: "555-555-5555"},
	// }, http.StatusOK},
}

func TestHandlers(t *testing.T) {
	routes := getRoutes()
	testServer := httptest.NewTLSServer(routes)
	defer testServer.Close()

	for _, test := range testCases {
		response, err := testServer.Client().Get(testServer.URL + test.url)
		if err != nil {
			t.Fatal(err)
		}
		if response.StatusCode != test.expectedStatusCode {
			t.Errorf("for the %s test, expected %d status code but got %d", test.name, test.expectedStatusCode, response.StatusCode)
		}
	}
}

func TestRepo_Confirm(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Bilbo's Suite",
		},
	}
	req, _ := http.NewRequest("GET", "/confirm", nil)
	ctx := getContext(req)
	req = req.WithContext(ctx) // need that "X-Session" header

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)

	handler := http.HandlerFunc(Repo.Confirm)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Confirm handler returned wrong code: got %d wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session
	req, _ = http.NewRequest("GET", "/confirm", nil)
	ctx = getContext(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Confirm handler returned wrong code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// test case with non-existent room
	req, _ = http.NewRequest("GET", "/confirm", nil)
	ctx = getContext(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Confirm handler returned wrong code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

// TestRepo_PostConfirm tests the PostConfirm Handler Function
func TestRepo_PostConfirm(t *testing.T) {
	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "Bilbo's Suite",
		},
	}

	// Manual way of building the request body
	reqBody := "start_date=2050-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2050-01-05")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=John")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=Doe")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=john@doe.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=555-555-5555")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")
	log.Println(reqBody)

	// This is a cleaner way to create a Request Body
	postedData := url.Values{}
	postedData.Add("start_date", "2050-01-01")
	postedData.Add("end_date", "2050-01-05")
	postedData.Add("first_name", "john")
	postedData.Add("last_name", "doe")
	postedData.Add("email", "john@doe.com")
	postedData.Add("phone", "555-555-5555")
	postedData.Add("room_id", "1")

	//req, _ := http.NewRequest("POST", "/confirm", strings.NewReader(reqBody))
	req, _ := http.NewRequest("POST", "/confirm", strings.NewReader(postedData.Encode()))
	ctx := getContext(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.PostConfirm)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostConfirm  handler returned wrong code: got %d wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Other test cases...
	// Missing request body
	req, _ = http.NewRequest("POST", "/confirm", nil)
	ctx = getContext(req)
	req = req.WithContext(ctx)
	session.Put(ctx, "reservation", reservation)

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.PostConfirm)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("PostConfirm  handler returned wrong code: got %d wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	// bad start date

}

func getContext(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println("handlers_test.go - getContext() Function", err.Error())
	}
	return ctx
}

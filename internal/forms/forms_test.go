package forms

import (
	"net/http"
	"net/url"
	"testing"
)

// When the function you're testing has a receiver, that has to be baked into the name
// For example, the IsValid() function takes a receiver of Form
func TestForm_IsValid(t *testing.T) {
	request, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	isValid := form.IsValid()
	if !isValid {
		t.Error("the form is invalid but should not have been")
	}
}

func TestForm_Required(t *testing.T) {
	request, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)
	form.Required("a", "b", "c")

	if form.IsValid() {
		t.Error("form shows valid when required fields are missing")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "b")
	postedData.Add("c", "c")

	request, _ = http.NewRequest("POST", "/whatever", nil)

	request.PostForm = postedData
	form = New(request.PostForm)
	form.Required("a", "b", "c")

	if !form.IsValid() {
		t.Error("says does not have all required fields when it in fact does :-( ")
	}
}

func TestForm_Has(t *testing.T) {
	request, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	has := form.Has("whatever")
	if has {
		t.Error("form shows it has a field when it does not")
	}

	postedData := url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)

	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	request, _ := http.NewRequest("POST", "/whatever", nil)
	form := New(request.PostForm)

	form.MinLength("x", 10)
	if form.IsValid() {
		t.Error("form shows MinLength is okay for non-existent field")
	}

	getData := form.Errors.Get("x")
	if getData == "" {
		t.Error("the Get() function should have returned an error but did not")
	}

	postedData := url.Values{}
	postedData.Add("some_field", "some_value")
	form = New(postedData)

	form.MinLength("some_field", 100)
	if form.IsValid() {
		t.Error("shows MinLength() of 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("some_field", "some_other_value")
	form = New(postedData)

	form.MinLength("some_field", 5)
	if !form.IsValid() {
		t.Error("shows MinLength() not met when it was")
	}

	getData = form.Errors.Get("some_field")
	if getData != "" {
		t.Error("the Get() function should not have an error but does")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("DummyField")

	if form.IsValid() {
		t.Error("form shows valid email for non-existent field")
	}

	// GOOD TEST
	postedData = url.Values{}
	postedData.Add("email", "john@doe.com")
	form = New(postedData)

	form.IsEmail("email")

	if !form.IsValid() {
		t.Error("got an invalid email when we should not have")
	}

	postedData = url.Values{}
	postedData.Add("email", "john")
	form = New(postedData)

	form.IsEmail("email")

	if form.IsValid() {
		t.Error("got valid for invalid email address")
	}
}

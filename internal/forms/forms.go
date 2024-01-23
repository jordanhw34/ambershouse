package forms

import (
	"fmt"
	"net/url"
	"strings"

	"github.com/asaskevich/govalidator"
)

// Form struct => creates a custom Form struct and embeds a url.Values object
type Form struct {
	url.Values
	Errors errors
}

// New => initializes a Form struct
func New(data url.Values) *Form {
	return &Form{
		data,
		errors(map[string][]string{}),
	}
}

// IsValid => returns true if the length of the errors slice is equal to 0
func (form *Form) IsValid() bool {
	return len(form.Errors) == 0
}

func (form *Form) Required(fields ...string) {
	for _, field := range fields {
		value := form.Get(field)
		if strings.TrimSpace(value) == "" {
			form.Errors.Add(field, "This field cannot be blank")
		}
	}
}

// Has => checks if Form field is in post and not empty
// We don't actually need the Request in this method since we can pull from the Form receiver
// func (form *Form) Has(field string, r *http.Request) bool {
func (form *Form) Has(field string) bool {
	// Bug here => we should not be checking the Request but rather should be checking the Form in our receiver
	//data := r.Form.Get(field)
	data := form.Get(field)
	return data != ""
}

// MinLength => Can make fields require a minimum length
// We don't actually need the Request in this method since we can pull from the Form receiver
// func (form *Form) MinLength(field string, length int, r *http.Request) bool {
func (form *Form) MinLength(field string, length int) bool {
	// Bug here => we should not be checking the Request but rather should be checking the Form in our receiver
	//data := r.Form.Get(field)
	data := form.Get(field)
	if len(data) < length {
		form.Errors.Add(field, fmt.Sprintf("This field must be at least %d characters long", length))
		return false
	}

	return true
}

// Checks for valid email address
func (form *Form) IsEmail(field string) {
	if !govalidator.IsEmail(form.Get(field)) {
		form.Errors.Add(field, "Invalid email address")
	}
}

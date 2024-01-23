package models

import "github.com/jordanhw34/ambershouse/internal/forms"

// This only exists to be imported by other packages

// TemplateData => Holds data sent from handlers to templates
type TemplateData struct {
	StringMap map[string]string
	IntMap    map[string]int
	FloatMap  map[string]float32
	Data      map[string]interface{} // can be anything
	CSRFToken string
	Flash     string
	Warning   string
	Error     string
	Form      *forms.Form
}

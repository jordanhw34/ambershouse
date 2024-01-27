package config

import (
	"log"
	"text/template"

	"github.com/alexedwards/scs/v2"
	"github.com/jordanhw34/ambershouse/internal/models"
)

// This package should be imported where it needs to be
// But it doesn't import anything from the application itself ... avoid import cycle error

// Holds the application config settings
type AppConfig struct {
	UseCache      bool
	TemplateCache map[string]*template.Template
	InfoLog       *log.Logger
	ErrorLog      *log.Logger
	InProduction  bool
	Session       *scs.SessionManager
	DateFormat    string
	MailChan      chan models.MailData
}

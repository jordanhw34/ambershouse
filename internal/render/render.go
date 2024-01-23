package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jordanhw34/ambershouse/internal/config"
	"github.com/jordanhw34/ambershouse/internal/models"
	"github.com/justinas/nosurf"
)

// Application Config
var app *config.AppConfig
var pathToTemplates = "./templates"

// NewTemplates =>
func NewRenderer(appConfig *config.AppConfig) {
	app = appConfig
}

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {
	// stringMap := make(map[string]string)
	// stringMap["siteTitle"] = "webapp1 go app from udemy"
	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)
	return td
}

// RenderTemplate => renders an html template
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var templateCache map[string]*template.Template

	if app.UseCache {
		// Get Template Cache from Application Config
		templateCache = app.TemplateCache
	} else {
		templateCache, _ = CreateTemplateCache()
	}

	// Get requested template from cache
	template, ok := templateCache[tmpl]
	if !ok {
		return errors.New("cannot get template from cache")
	}

	// Doing this for finer grain error checking => maybe something got parsed but it won't execute
	buf := new(bytes.Buffer)

	// Get Default Data
	td = AddDefaultData(td, r)

	// Execute the template
	_ = template.Execute(buf, td)

	// Render the template => different error message in next lecture?
	bytesWritten, err := buf.WriteTo(w)
	log.Println(" > Number of bytes written:", bytesWritten)
	//log.Println(w)
	if err != nil {
		log.Println("error writing template to browser", err.Error())
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//templateCache := make(map[string]*template.Template) => This syntax is the same as using the "make" keyword, creates and initializes an empty slice
	templateCache := map[string]*template.Template{}

	log.Println(" > Building the template cache")
	// look at templates folder and get all files named *.page.html
	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.html", pathToTemplates))
	if err != nil {
		log.Println("no error getting files from the ./templates folder matching *.page.html :-) ")
		return templateCache, err
	}

	// Range through files that were found with filepath.Glob
	for _, page := range pages {
		//log.Println("ranging template pages, current page name =", page)		// e.x. templates\about.page.html
		// page is the full path to the template but we only want the base name
		name := filepath.Base(page)
		//log.Println(" > Building Template Cache for Page =", name)

		// now we need to parse the file
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		// Now we need to look for layouts used in these templates
		layouts, err := filepath.Glob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
		if err != nil {
			return templateCache, err
		}

		// If at least 1 layout is found
		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob(fmt.Sprintf("%s/*.layout.html", pathToTemplates))
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = templateSet
	}

	return templateCache, nil
}

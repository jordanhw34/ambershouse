package render

import (
	"bytes"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/jordanhw34/ambershouse/pkg/config"
	"github.com/jordanhw34/ambershouse/pkg/models"
)

// Application Config
var app *config.AppConfig

// NewTemplates =>
func NewTemplates(a *config.AppConfig) {
	app = a
}

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	// stringMap := make(map[string]string)
	// stringMap["siteTitle"] = "webapp1 go app from udemy"
	return td
}

// RenderTemplate => renders an html template
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

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
		log.Fatal("could not find this template in the cache")
	}

	// Doing this for finer grain error checking => maybe something got parsed but it won't execute
	buf := new(bytes.Buffer)

	// Get Default Data
	td = AddDefaultData(td)

	// Execute the template
	_ = template.Execute(buf, td)

	// Render the template => different error message in next lecture?
	bytesWritten, err := buf.WriteTo(w)
	log.Println(" > Number of bytes written:", bytesWritten)
	//log.Println(w)
	if err != nil {
		log.Println("error writing template to browser", err.Error())
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {
	//templateCache := make(map[string]*template.Template) => This syntax is the same as using the "make" keyword, creates and initializes an empty slice
	templateCache := map[string]*template.Template{}

	log.Println(" > Building the template cache")
	// look at templates folder and get all files named *.page.html
	pages, err := filepath.Glob("./templates/*.page.html")
	if err != nil {
		log.Println("no error getting files from the ./templates folder matching *.page.html :-) ")
		return templateCache, err
	}

	// Range through files that were found with filepath.Glob
	for _, page := range pages {
		//log.Println("ranging template pages, current page name =", page)		// e.x. templates\about.page.html
		// page is the full path to the template but we only want the base name
		name := filepath.Base(page)
		log.Println(" > ranging template pages, current template name =", name)

		// now we need to parse the file
		templateSet, err := template.New(name).ParseFiles(page)
		if err != nil {
			return templateCache, err
		}

		// Now we need to look for layouts used in these templates
		layouts, err := filepath.Glob("./templates/*.layout.html")
		if err != nil {
			return templateCache, err
		}

		// If at least 1 layout is found
		if len(layouts) > 0 {
			templateSet, err = templateSet.ParseGlob("./templates/*.layout.html")
			if err != nil {
				return templateCache, err
			}
		}

		templateCache[name] = templateSet
	}

	return templateCache, nil
}

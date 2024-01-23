package render

import (
	"net/http"
	"testing"

	"github.com/jordanhw34/ambershouse/internal/models"
)

func TestAddDefaultData(t *testing.T) {
	var td models.TemplateData
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	session.Put(r.Context(), "flash", "flash message")

	result := AddDefaultData(&td, r)

	if result.Flash != "flash message" {
		t.Error("flash message was not found in session")
	}
}

func TestRenderTemplate(t *testing.T) {
	pathToTemplates = "./../../templates"
	templateCache, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = templateCache

	request, err := getSession()
	if err != nil {
		t.Error(err)
	}

	var resWriter responseWriter

	err = Template(&resWriter, request, "home.page.html", &models.TemplateData{})
	if err != nil {
		t.Error("error writing template to test client")
	}

	err = Template(&resWriter, request, "non-existent.page.html", &models.TemplateData{})
	if err == nil {
		t.Error("rendered template that does not exist")
	}
}

func getSession() (*http.Request, error) {
	r, err := http.NewRequest("GET", "/some-url", nil)
	if err != nil {
		return nil, err
	}

	ctx := r.Context()
	ctx, _ = session.Load(ctx, r.Header.Get("X-Session"))
	r = r.WithContext(ctx)

	return r, nil
}

func TestNewTemplates(t *testing.T) {
	NewRenderer(app)
}

func TestCreateTemplateCache(t *testing.T) {
	pathToTemplates = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

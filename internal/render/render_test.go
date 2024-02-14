package render

import (
	"net/http"
	"testing"

	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
)

func TestAddDefaultData(t *testing.T) {

	var td models.TemplateData

	r, err := getSession()
	if err != nil {
		t.Error(err)
	}
	session.Put(r.Context(), "flash", "123")

	result := AddDefaultData(&td, r)
	if result.Flash != "123" {
		t.Error("flash value 123 not found Failed")
	}
}

func TestRenderTemplate(t *testing.T) {

	pathToTemplate = "./../../templates"

	tc, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}

	app.TemplateCache = tc

	// for building request
	r, err := getSession()
	if err != nil {
		t.Error(err)
	}

	// Response writer created (check ssetup_tes for mywriter details)
	var ww myWriter

	err = Template(&ww, r, "home.page.tmpl", &models.TemplateData{})
	if err != nil {
		t.Error("error writing to template")
	}
	// err = RenderTemplate(&ww, r, "no.page.tmpl", &models.TemplateData{})
	// if err != nil {
	// 	t.Error("rendered temp that not exist")
	// }
}

func getSession() (*http.Request, error) {
	// http.NewRequest() trying to GET request at url
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

	pathToTemplate = "./../../templates"
	_, err := CreateTemplateCache()
	if err != nil {
		t.Error(err)
	}
}

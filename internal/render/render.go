package render

import (
	"bytes"
	"errors"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"
	"time"

	"github.com/M-Abdullah-Nazeer/bookings/internal/config"
	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
	"github.com/justinas/nosurf"
)

// holds funcitions that we make available to golang templates
var functions = template.FuncMap{

	"humanDate":  HumanDate,
	"formatDate": FormatDate,
	"iterate":    Iterate,
	"add":        Add,
}

func HumanDate(t time.Time) string {
	return t.Format("2006-01-02")
}

func Add(a, b int) int {
	return a + b
}
func Iterate(count int) []int {
	var i int
	var items []int

	for i = 0; i < count; i++ {
		items = append(items, i)
	}
	return items
}

func FormatDate(t time.Time, format string) string {
	return t.Format(format)
}

var pathToTemplate = "./templates"
var app *config.AppConfig
var err error

func AddDefaultData(td *models.TemplateData, r *http.Request) *models.TemplateData {

	td.Flash = app.Session.PopString(r.Context(), "flash")
	td.Error = app.Session.PopString(r.Context(), "error")
	td.Warning = app.Session.PopString(r.Context(), "warning")
	td.CSRFToken = nosurf.Token(r)

	if app.Session.Exists(r.Context(), "user_id") {
		// 1 means login 0 means logout, 0 is default value of int
		td.IsAuthenticated = 1
	}

	return td
}

func NewRenderer(a *config.AppConfig) {
	app = a
}
func Template(w http.ResponseWriter, r *http.Request, tmpl string, td *models.TemplateData) error {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
	}
	if err != nil {
		log.Fatal("errr in CreateTemplateCache()")
		return err
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("Did not get temp frm temp cache")
		return errors.New("can't get template from cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td, r)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Err writing temp to browser", err)
		return err
	}

	return nil
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob(fmt.Sprintf("%s/*.page.tmpl", pathToTemplate))
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.New(name).Funcs(functions).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob(fmt.Sprintf("%s/*.layout.tmpl", pathToTemplate))
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts

	}
	return myCache, nil
}

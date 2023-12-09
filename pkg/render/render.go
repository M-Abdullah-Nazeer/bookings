package render

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"path/filepath"
	"text/template"

	"github.com/M-Abdullah-Nazeer/LearningGO/pkg/config"
	"github.com/M-Abdullah-Nazeer/LearningGO/pkg/models"
)

var app *config.AppConfig
var err error

func AddDefaultData(td *models.TemplateData) *models.TemplateData {
	return td
}

func NewTemplates(a *config.AppConfig) {
	app = a
}
func RenderTemplate(w http.ResponseWriter, tmpl string, td *models.TemplateData) {

	var tc map[string]*template.Template

	if app.UseCache {
		tc = app.TemplateCache
	} else {
		tc, err = CreateTemplateCache()
	}
	if err != nil {
		log.Fatal("errr in CreateTemplateCache()")
	}

	t, ok := tc[tmpl]

	if !ok {
		log.Fatal("not get temp frm temp cache")
	}

	buf := new(bytes.Buffer)
	td = AddDefaultData(td)

	_ = t.Execute(buf, td)

	_, err := buf.WriteTo(w)
	if err != nil {
		fmt.Println("Err writing temp to browser", err)
	}
}

func CreateTemplateCache() (map[string]*template.Template, error) {

	myCache := map[string]*template.Template{}

	pages, err := filepath.Glob("./templates/*.page.tmpl")
	if err != nil {
		return myCache, err
	}
	for _, page := range pages {

		name := filepath.Base(page)
		ts, err := template.New(name).ParseFiles(page)
		if err != nil {
			return myCache, err
		}
		matches, err := filepath.Glob("./templates/*.layout.tmpl")
		if err != nil {
			return myCache, err
		}
		if len(matches) > 0 {
			ts, err = ts.ParseGlob("./templates/*.layout.tmpl")
			if err != nil {
				return myCache, err
			}
		}
		myCache[name] = ts

	}
	return myCache, nil
}

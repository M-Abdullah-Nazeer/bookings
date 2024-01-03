package main

import (
	"encoding/gob"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/M-Abdullah-Nazeer/bookings/internal/config"
	"github.com/M-Abdullah-Nazeer/bookings/internal/handlers"
	"github.com/M-Abdullah-Nazeer/bookings/internal/helpers"
	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
	"github.com/M-Abdullah-Nazeer/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":8080"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	err := run()
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Starting Application on port", portNumber)

	srv := &http.Server{
		Addr:    portNumber,
		Handler: routes(&app),
	}
	err = srv.ListenAndServe()
	if err != nil {
		log.Fatal(err)
	}
}

func run() error {

	// what am i going to store in the session
	gob.Register(models.Reservation{})

	app.InProduction = false
	// os.Stdout is terminal window, \t is tab space
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	app.InfoLog = InfoLog

	// log.Lshortfile will give info about error
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	app.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = app.InProduction

	app.Session = session

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Fatal("can't create temp cache")
		return err
	}

	app.TemplateCache = tc

	app.UseCache = false

	repo := handlers.NewRepo(&app)
	handlers.NewHandlers(repo)

	render.NewTemplates(&app)

	helpers.NewHelpers(&app)

	return nil
}

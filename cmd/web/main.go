package main

import (
	"encoding/gob"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/M-Abdullah-Nazeer/bookings/internal/config"
	"github.com/M-Abdullah-Nazeer/bookings/internal/driver"
	"github.com/M-Abdullah-Nazeer/bookings/internal/handlers"
	"github.com/M-Abdullah-Nazeer/bookings/internal/helpers"

	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
	"github.com/M-Abdullah-Nazeer/bookings/internal/render"
	"github.com/alexedwards/scs/v2"
)

const portNumber = ":80"

var app config.AppConfig
var session *scs.SessionManager

func main() {
	db, err := run()
	if err != nil {
		log.Fatal(err)
	}
	// close db
	defer db.SQL.Close()

	// close mail channel
	defer close(app.MailChan)

	listenForMail()

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

func run() (*driver.DB, error) {

	// what am i going to store in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	// read flags
	inProduction := flag.Bool("production", true, "Application is in Production")
	useCache := flag.Bool("cache", true, "Use template cache")
	dbHost := flag.String("dbhost", "localhost", "Database host")
	dbName := flag.String("dbname", "", "Database name")
	dbUser := flag.String("dbuser", "", "Database user")
	dbPass := flag.String("dbpass", "", "Database password")
	dbPort := flag.String("dbport", "5432", "Database port")
	dbSSL := flag.String("dbssl", "disable", "Database ssl settnig (disable, prefer, require)")

	flag.Parse()

	if *dbName == "" || *dbUser == "" {
		fmt.Println("Missng required flags")
		os.Exit(1)
	}

	// making mail channel
	mailChan := make(chan models.MailData)
	app.MailChan = mailChan

	app.InProduction = *inProduction
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

	// connect to db

	log.Println("Connecting to db....")

	connectionString := fmt.Sprintf("host=%s port=%s dbname=%s user=%s password=%s sslmode=%s", *dbHost, *dbPort, *dbName, *dbUser, *dbPass, *dbSSL)

	db, error := driver.ConnectSQL(connectionString)

	// db, error := driver.ConnectSQL("host=localhost port=5432 dbname=bookings user=postgres password=admin")
	if error != nil {

		log.Fatal("can't connect to db, dying.....")
		return nil, error
	}

	log.Println("Connected to db")

	tc, err := render.CreateTemplateCache()

	if err != nil {
		log.Println(err)
		log.Fatal("can't create temp cache")
		return nil, err
	}

	app.TemplateCache = tc

	app.UseCache = *useCache

	repo := handlers.NewRepo(&app, db)
	handlers.NewHandlers(repo)

	render.NewRenderer(&app)

	helpers.NewHelpers(&app)

	return db, nil
}

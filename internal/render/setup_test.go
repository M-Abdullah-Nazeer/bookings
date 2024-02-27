package render

import (
	"encoding/gob"
	"log"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/M-Abdullah-Nazeer/bookings/internal/config"
	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
	"github.com/alexedwards/scs/v2"
)

var session *scs.SessionManager
var testApp config.AppConfig

func TestMain(m *testing.M) {

	// what am i going to store in the session
	gob.Register(models.Reservation{})
	gob.Register(models.User{})
	gob.Register(models.Room{})
	gob.Register(models.Restriction{})
	gob.Register(map[string]int{})

	testApp.InProduction = false
	// os.Stdout is terminal window, \t is tab space
	InfoLog := log.New(os.Stdout, "INFO\t", log.Ldate|log.Ltime)
	testApp.InfoLog = InfoLog

	// log.Lshortfile will give info about error
	errorLog := log.New(os.Stdout, "ERROR\t", log.Ldate|log.Ltime|log.Lshortfile)
	testApp.ErrorLog = errorLog

	session = scs.New()
	session.Lifetime = 24 * time.Hour
	session.Cookie.Persist = true
	session.Cookie.SameSite = http.SameSiteLaxMode
	session.Cookie.Secure = false //not https thats why89\]

	testApp.Session = session

	app = &testApp

	os.Exit(m.Run())
}

type myWriter struct{}

func (mw *myWriter) Header() http.Header {

	var h http.Header
	return h
}

func (mw *myWriter) WriteHeader(i int) {}

func (mw *myWriter) Write(b []byte) (int, error) {

	len := len(b)

	return len, nil
}

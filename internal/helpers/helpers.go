package helpers

import (
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/M-Abdullah-Nazeer/bookings/internal/config"
)

var app *config.AppConfig

func NewHelpers(a *config.AppConfig) {

	app = a
}

func ClientError(w http.ResponseWriter, status int) {

	app.InfoLog.Println("Client error with status of", status)
	http.Error(w, http.StatusText(status), status)
}

func ServerError(w http.ResponseWriter, err error) {

	// consist of error message and stack trace associated with it
	trace := fmt.Sprintf("%s\n%s", err.Error(), debug.Stack())
	// writes to the terminal
	app.ErrorLog.Println(trace)

	// gets feedback from user
	http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
}

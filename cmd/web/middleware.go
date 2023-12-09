package main

import (
	"net/http"

	"github.com/justinas/nosurf"
)

// adds CSRF protection to all POST requests
func NoSurf(next http.Handler) http.Handler {

	csrfHandler := nosurf.New(next)
	csrfHandler.SetBaseCookie(http.Cookie{
		HttpOnly: true,
		Path:     "/",
		Secure:   app.InProduction,
		SameSite: http.SameSiteLaxMode,
	})
	return csrfHandler
}

//SessionLoad loads saves session fr every request

func SessionLoad(next http.Handler) http.Handler {

	return session.LoadAndSave(next)
}

package main

import (
	"net/http"

	"github.com/M-Abdullah-Nazeer/bookings/internal/helpers"
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

func Auth(next http.Handler) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		// if user_id is not in session
		if !helpers.IsAuthenticated(r) {
			session.Put(r.Context(), "error", "Login first!")
			http.Redirect(w, r, "/user/login", http.StatusSeeOther)
			return
		}
		// if authenticated
		next.ServeHTTP(w, r)
	})
}

package main

import (
	"fmt"
	"testing"

	"github.com/M-Abdullah-Nazeer/bookings/internal/config"
	"github.com/go-chi/chi/v5"
)

func TestRoutes(t *testing.T) {

	var myapp config.AppConfig

	h := routes(&myapp)

	switch dt := h.(type) {
	case *chi.Mux:
	// do nothing
	default:
		t.Error(fmt.Sprintf("Handler is not of type chi.Mux, it is of type %T", dt))

	}

}

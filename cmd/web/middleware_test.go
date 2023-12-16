package main

import (
	"fmt"
	"net/http"
	"testing"
)

func TestNoSurf(t *testing.T) {

	var myH myHandler

	h := NoSurf(&myH)

	switch dt := h.(type) {

	case http.Handler:
	// do nothing test passed

	default:
		// %T is placeholder for type
		t.Error(fmt.Sprintf("The type returned by NoSurf is not handler instead it is %T", dt))
	}
}
func TestSessionLoad(t *testing.T) {

	var myH myHandler

	h := SessionLoad(&myH)

	switch dt := h.(type) {

	case http.Handler:
	// do nothing test passed

	default:
		// %T is placeholder for type
		t.Error(fmt.Sprintf("The type returned by NoSurf is not handler instead it is %T", dt))
	}
}

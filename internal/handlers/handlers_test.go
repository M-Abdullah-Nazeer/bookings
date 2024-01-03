package handlers

import (
	"net/http"
	"net/http/httptest"
	"net/url"

	"testing"
)

// performing table test

type postData struct {
	key   string
	value string
}

var theTest = []struct {
	name               string
	url                string
	method             string
	params             []postData
	expectedStatusCode int
}{

	{"home", "/", "GET", []postData{}, http.StatusOK},
	{"about", "/about", "GET", []postData{}, http.StatusOK},
	{"gs", "/generals-quarters", "GET", []postData{}, http.StatusOK},
	{"ms", "/majors-suite", "GET", []postData{}, http.StatusOK},
	{"contact", "/contact", "GET", []postData{}, http.StatusOK},
	{"sa", "/search-availability", "GET", []postData{}, http.StatusOK},
	{"rs", "/reservation-summary", "GET", []postData{}, http.StatusOK},
	{"mr", "/make-reservation", "GET", []postData{}, http.StatusOK},
	{"post search-availability", "/search-availability", "POST", []postData{
		{key: "start", value: "2023-03-01"},
		{key: "start", value: "2023-03-02"},
	}, http.StatusOK},
	{"post search-availability-json", "/search-availability-json", "POST", []postData{
		{key: "start", value: "2023-03-01"},
		{key: "start", value: "2023-03-02"},
	}, http.StatusOK},
	{"make reservation post", "/make-reservation", "POST", []postData{
		{key: "first_name", value: "Abdullah"},
		{key: "last_name", value: "Nazeer"},
		{key: "email", value: "shareef@gmail.com"},
		{key: "phone", value: "0555513515"},
	}, http.StatusOK},
}

func TestHandlers(t *testing.T) {

	routes := getRoutes()

	ts := httptest.NewTLSServer(routes)
	defer ts.Close()

	for _, e := range theTest {

		if e.method == "GET" {

			resp, err := ts.Client().Get(ts.URL + e.url)

			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}
		} else { //else means method = POST
			// url.Values{} is part of standard library, holds info as a post request for a variable

			values := url.Values{}

			for _, z := range e.params {

				values.Add(z.key, z.value)
			}
			resp, err := ts.Client().PostForm(ts.URL+e.url, values)
			if err != nil {
				t.Log(err)
				t.Fatal(err)
			}

			if resp.StatusCode != e.expectedStatusCode {
				t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
			}

		}

	} //for close

}

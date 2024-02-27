package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"

	"testing"

	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
)

// performing table test

// type postData struct {
// 	key   string
// 	value string
// }

var theTest = []struct {
	name   string
	url    string
	method string
	// params             []postData
	expectedStatusCode int
}{

	{"home", "/", "GET", http.StatusOK},
	{"about", "/about", "GET", http.StatusOK},
	{"gs", "/generals-quarters", "GET", http.StatusOK},
	{"ms", "/majors-suite", "GET", http.StatusOK},
	{"contact", "/contact", "GET", http.StatusOK},
	{"sa", "/search-availability", "GET", http.StatusOK},
	{"rs", "/reservation-summary", "GET", http.StatusOK},
	{"non-existent", "/kch/bhi/link", "GET", http.StatusNotFound},

	// new routes
	{"login", "/user/login", "GET", http.StatusOK},
	{"logout", "/user/logout", "GET", http.StatusOK},
	{"dashboard", "/admin/dashboard", "GET", http.StatusOK},
	{"new res", "/admin/reservations-new", "GET", http.StatusOK},
	{"all res", "/admin/reservations-all", "GET", http.StatusOK},
	{"show res", "/admin/reservations/new/1/show", "GET", http.StatusOK},

	// {"post search-availability", "/search-availability", "POST", []postData{
	// 	{key: "start", value: "2023-03-01"},
	// 	{key: "start", value: "2023-03-02"},
	// }, http.StatusOK},
	// {"post search-availability-json", "/search-availability-json", "POST", []postData{
	// 	{key: "start", value: "2023-03-01"},
	// 	{key: "start", value: "2023-03-02"},
	// }, http.StatusOK},
	// {"make reservation post", "/make-reservation", "POST", []postData{
	// 	{key: "first_name", value: "Abdullah"},
	// 	{key: "last_name", value: "Nazeer"},
	// 	{key: "email", value: "shareef@gmail.com"},
	// 	{key: "phone", value: "0555513515"},
	// }, http.StatusOK},
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
		}
		// else { //else means method = POST
		// 	// url.Values{} is part of standard library, holds info as a post request for a variable

		// 	values := url.Values{}

		// 	for _, z := range e.params {

		// 		values.Add(z.key, z.value)
		// 	}
		// 	resp, err := ts.Client().PostForm(ts.URL+e.url, values)
		// 	if err != nil {
		// 		t.Log(err)
		// 		t.Fatal(err)
		// 	}

		// 	if resp.StatusCode != e.expectedStatusCode {
		// 		t.Errorf("for %s, expected %d but got %d", e.name, e.expectedStatusCode, resp.StatusCode)
		// 	}

		// }

	} //for close

}

func TestRepository_Reservation(t *testing.T) {

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/make-reservation", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusOK {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusOK)
	}

	// test case where reservation is not in session (reset everything)
	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}

	req, _ = http.NewRequest("GET", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	reservation.RoomID = 100
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.Reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusTemporaryRedirect {
		t.Errorf("Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusTemporaryRedirect)
	}
}

func TestRepository_PostReservation(t *testing.T) {
	reqBody := "start_date=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=khanbaba")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=khanbaba")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abcdef@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789987")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	// Built request and sent form data as post body, NewReader is io reader
	req, _ := http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()

	// it says to web server that request is form post
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded") //simulate a form submission

	handler := http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for missing post body, first if fail

	req, _ = http.NewRequest("POST", "/make-reservation", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// Session.Get() pass case
	reqBody = "start_date=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=khanbaba")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=khanbaba")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abcdef@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789987")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")
	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	reservation := models.Reservation{}
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.PostReservation)

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid data
	reqBody = "start_date=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abcdef@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789987")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	reservation = models.Reservation{
		FirstName: "a",
		LastName:  "a",
	}
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code for invalid form data: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for InsertReservation() failure
	reqBody = "start_date=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=khanbaba")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=khanbaba")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abcdef@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789987")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=3")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	reservation = models.Reservation{
		RoomID: 3,
	}
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code for InsertReservation() failure: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for InsertRoomRestriction() failure
	reqBody = "start_date=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end_date=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "first_name=khanbaba")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "last_name=khanbaba")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "email=abcdef@gmail.com")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "phone=123456789987")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=200")

	req, _ = http.NewRequest("POST", "/make-reservation", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostReservation)
	reservation = models.Reservation{
		RoomID: 200,
	}
	session.Put(ctx, "reservation", reservation)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code for InsertRoomRestriction()  failure: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_AvailabilityJSON(t *testing.T) {

	// Case where form is parsed
	reqBody := "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=1")

	req, _ := http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	rr := httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler := http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	var j jsonResponse
	err := json.Unmarshal(rr.Body.Bytes(), &j)

	if err != nil {
		t.Error("Could not parse JSON")
	}

	// case where form cannot be parsed

	req, _ = http.NewRequest("POST", "/search-availability-json", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Errorf("Could not parse JSON: %v", err)
	}

	// test for invalid start date
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code for invalid start date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid end date
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=2")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code for invalid end date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// test for invalid room id
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=invalid")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("Post Reservation handler returned wrong response code for room id: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case for SearchAvailabilityByRoomID() failure
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	reqBody = fmt.Sprintf("%s&%s", reqBody, "room_id=300")

	req, _ = http.NewRequest("POST", "/search-availability-json", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	rr = httptest.NewRecorder()
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.AvailabilityJSON)
	handler.ServeHTTP(rr, req)

	err = json.Unmarshal([]byte(rr.Body.String()), &j)

	if err != nil {
		t.Errorf("Could not parse JSON: %v", err)
	}

}

func TestRepository_ChooseRoom(t *testing.T) {

	// case where id is correct

	reservation := models.Reservation{
		RoomID: 1,
		Room: models.Room{
			ID:       1,
			RoomName: "General's Quarters",
		},
	}

	req, _ := http.NewRequest("GET", "/choose-room/1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	req.RequestURI = "/choose-room/1"
	rr := httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler := http.HandlerFunc(Repo.ChooseRoom)
	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case where id is incorrect
	req, _ = http.NewRequest("GET", "/choose-room/invald", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/invald"
	rr = httptest.NewRecorder()
	session.Put(ctx, "reservation", reservation)
	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// session failure
	req, _ = http.NewRequest("GET", "/choose-room/1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.RequestURI = "/choose-room/1"
	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.ChooseRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("ChooseRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_BookRoom(t *testing.T) {

	// case where date is parsed
	req, _ := http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=1", nil)
	ctx := getCtx(req)
	req = req.WithContext(ctx)

	rr := httptest.NewRecorder()

	handler := http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case where date is not parsed
	req, _ = http.NewRequest("GET", "/book-room?s=2050-01-01&e=2050-01-02&id=invalid", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case for invalid start date
	req, _ = http.NewRequest("GET", "/book-room?s=invalid&e=2050-01-02&id=1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case for invalid end date
	req, _ = http.NewRequest("GET", "/book-room?s=2050-01-01&e=invalid&id=1", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	rr = httptest.NewRecorder()

	handler = http.HandlerFunc(Repo.BookRoom)

	handler.ServeHTTP(rr, req)

	if rr.Code != http.StatusSeeOther {
		t.Errorf("BookRoom handler returned wrong response code: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

}

func TestRepository_PostAvailability(t *testing.T) {

	// case where form is parsed

	reqBody := "start=2060-01-02"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	req, _ := http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx := getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(Repo.PostAvailability)
	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code form is not parsed: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	// case where form is not parsed

	req, _ = http.NewRequest("POST", "/search-availability", nil)
	ctx = getCtx(req)
	req = req.WithContext(ctx)

	handler = http.HandlerFunc(Repo.PostAvailability)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code form is not parsed: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

	// case where start date is invalid
	reqBody = "start=invalid"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostAvailability)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code due to invalid start date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	// case where end date is invalid
	reqBody = "start=2060-01-01"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=invalid")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostAvailability)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code due to invalid end date: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}
	// case where SearchAvailabilityForAllRooms() fails
	reqBody = "start=2060-01-20"
	reqBody = fmt.Sprintf("%s&%s", reqBody, "end=2060-01-02")
	req, _ = http.NewRequest("POST", "/search-availability", strings.NewReader(reqBody))
	ctx = getCtx(req)
	req = req.WithContext(ctx)
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	handler = http.HandlerFunc(Repo.PostAvailability)

	rr = httptest.NewRecorder()

	handler.ServeHTTP(rr, req)
	if rr.Code != http.StatusSeeOther {
		t.Errorf("PostAvailability handler returned wrong response code because SearchAvailabilityForAllRooms failed: got %d, wanted %d", rr.Code, http.StatusSeeOther)
	}

}

var loginTests = []struct {
	name               string
	email              string
	expectedStatusCode int
	expectedHTML       string
	expectedLocation   string
}{
	{
		"valid-credentials",
		"iam@me.com",
		http.StatusSeeOther,
		"",
		"/",
	},
	{
		"invalid-credentials",
		"jackma@se.com",
		http.StatusSeeOther,
		"",
		"/user/login",
	},
	{
		"invalid-data",
		"j",
		http.StatusOK,
		`action="/user/login"`,
		"",
	},
}

func TestLogin(t *testing.T) {

	// range through all test

	for _, e := range loginTests {
		postedData := url.Values{}
		postedData.Add("email", e.email)
		postedData.Add("password", "password")

		// create request

		req, _ := http.NewRequest("POST", "/user/login", strings.NewReader(postedData.Encode()))
		ctx := getCtx(req)
		req = req.WithContext(ctx)

		// set Header
		req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
		rr := httptest.NewRecorder()

		// call the handler
		handler := http.HandlerFunc(Repo.PostLogin)
		handler.ServeHTTP(rr, req)

		if rr.Code != e.expectedStatusCode {
			t.Errorf("failed %s: expected code %d, but got %d", e.name, e.expectedStatusCode, rr.Code)
		}

		if e.expectedLocation != "" {
			realLoc, _ := rr.Result().Location()
			if realLoc.String() != e.expectedLocation {
				t.Errorf("failed %s: expected location %s, but got location %s", e.name, e.expectedLocation, realLoc.String())
			}
		}

		// checking for expected values in HTML

		if e.expectedHTML != "" {
			// read response body into string
			html := rr.Body.String()
			if !strings.Contains(html, e.expectedHTML) {
				t.Errorf("failed %s: expected to find %s but did not", e.name, e.expectedHTML)
			}
		}
	}
}

func getCtx(req *http.Request) context.Context {
	ctx, err := session.Load(req.Context(), req.Header.Get("X-Session"))
	if err != nil {
		log.Println(err)
	}

	return ctx
}

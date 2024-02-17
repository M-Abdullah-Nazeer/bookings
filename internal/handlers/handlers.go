package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"net/http"

	"github.com/M-Abdullah-Nazeer/bookings/internal/config"
	"github.com/M-Abdullah-Nazeer/bookings/internal/driver"
	"github.com/M-Abdullah-Nazeer/bookings/internal/forms"
	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
	"github.com/M-Abdullah-Nazeer/bookings/internal/render"
	"github.com/M-Abdullah-Nazeer/bookings/internal/repository"
	"github.com/M-Abdullah-Nazeer/bookings/internal/repository/dbrepo"
	// "github.com/go-chi/chi"
)

var Repo *Repository

type Repository struct {
	App *config.AppConfig
	DB  repository.DatabaseRepo
}

func NewRepo(a *config.AppConfig, db *driver.DB) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewPostgresRepo(db.SQL, a),
	}
}

func NewTestRepo(a *config.AppConfig) *Repository {
	return &Repository{
		App: a,
		DB:  dbrepo.NewTestingRepo(a),
	}
}

func NewHandlers(r *Repository) {
	Repo = r
}

func (m *Repository) Home(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "home.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Contact(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "contact.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Generals(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "generals.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Majors(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "majors.page.tmpl", &models.TemplateData{})
}

func (m *Repository) Reservation(w http.ResponseWriter, r *http.Request) {

	// session contains now room id, start and end dates
	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {

		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}

	room, err := m.DB.GetRoomByID(res.RoomID)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get room details from db")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
		return
	}
	//saving room name
	res.Room.RoomName = room.RoomName

	// session contains now room id, start and end dates and room name
	m.App.Session.Put(r.Context(), "reservation", res)

	//parsing time to string and saving in map
	sd := res.StartDate.Format("2006-01-02")
	ed := res.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = res

	render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
		Form:      forms.New(nil),
		Data:      data,
		StringMap: stringMap,
	})
}

// PostReservation handles the posting of a reservation form
func (m *Repository) PostReservation(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()

	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't parse form")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)
	if !ok {

		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	reservation.FirstName = r.Form.Get("first_name")
	reservation.LastName = r.Form.Get("last_name")
	reservation.Email = r.Form.Get("email")
	reservation.Phone = r.Form.Get("phone")

	form := forms.New(r.PostForm)

	// form.Has("first_name", r)
	form.MinLength("first_name", 3)
	form.Required("first_name", "last_name", "email", "phone")
	form.IsEmail("email")

	if !form.Valid() {
		data := make(map[string]interface{})
		data["reservation"] = reservation

		http.Error(w, "Form is not valid", http.StatusSeeOther)
		//re render the page
		render.Template(w, r, "make-reservation.page.tmpl", &models.TemplateData{
			Form: form,
			Data: data,
		})
		return
	}

	reservationID, err := m.DB.InsertReservation(reservation)
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't insert reservation to db")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// session contains now all values of reservation model
	m.App.Session.Put(r.Context(), "reservation", reservation)

	restriction := models.RoomRestriction{
		RoomID:        reservation.RoomID,
		StartDate:     reservation.StartDate,
		EndDate:       reservation.EndDate,
		ReservationID: reservationID,
		RestrictionID: 1,
	}

	err = m.DB.InsertRoomRestriction(restriction)
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't insert room restrictions to db")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	// send email notification to guest
	htmlMsg := fmt.Sprintf(`<h1>Hello %s </h1>`, reservation.FirstName)
	msg := models.MailData{
		To:       reservation.Email,
		From:     "me@google.com",
		Subject:  "Reservation Confirmation",
		Content:  htmlMsg,
		Template: "reservation-email-to-guest.html",
	}
	m.App.MailChan <- msg

	// send email notification to property owner
	htmlMsg = fmt.Sprintf(`<h1>Your Room %s has been booked by %s for %s to %s </h1>`, reservation.Room.RoomName, reservation.FirstName, reservation.StartDate, reservation.EndDate)
	msg = models.MailData{
		To:      "me@google.com",
		From:    "me@google.com",
		Subject: "Reservation Confirmation",
		Content: htmlMsg,
	}
	m.App.MailChan <- msg

	// Storing form object values in session to send them to reservation-summary template
	m.App.Session.Put(r.Context(), "reservation", reservation)

	http.Redirect(w, r, "/reservation-summary", http.StatusSeeOther)

}

// handles reservation summary page
func (m *Repository) ReservationSummary(w http.ResponseWriter, r *http.Request) {

	// Getting values from Session
	reservation, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation) //session doesnot know type variable it is saving so typecasting

	if !ok {

		m.App.Session.Put(r.Context(), "error", "Can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusTemporaryRedirect)
	}

	//Remove values from session after successful retrieval
	m.App.Session.Remove(r.Context(), "reservation")

	// always convert dates to string before displaying on template

	sd := reservation.StartDate.Format("2006-01-02")
	ed := reservation.EndDate.Format("2006-01-02")

	stringMap := make(map[string]string)
	stringMap["start_date"] = sd
	stringMap["end_date"] = ed

	data := make(map[string]interface{})
	data["reservation"] = reservation
	render.Template(w, r, "reservation-summary.page.tmpl", &models.TemplateData{
		Data:      data,
		StringMap: stringMap,
	})
}

func (m *Repository) About(w http.ResponseWriter, r *http.Request) {

	stringMap := make(map[string]string)
	stringMap["test"] = "Hi MAN"

	remoteIP := m.App.Session.GetString(r.Context(), "remote_ip")
	stringMap["remote_ip"] = remoteIP

	render.Template(w, r, "about.page.tmpl", &models.TemplateData{
		StringMap: stringMap,
	})
}

func (m *Repository) Availability(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "search-availability.page.tmpl", &models.TemplateData{})
}

// finds room avilability from db and redirects to choose-room page if available
func (m *Repository) PostAvailability(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse form!")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02" // YYYY-M-D

	start_date, err := time.Parse(layout, start)
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	end_date, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	rooms, err := m.DB.SearchAvailabilityForAllRooms(start_date, end_date)

	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't get all rooms from db")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	if len(rooms) == 0 {
		m.App.Session.Put(r.Context(), "error", "No Rooms Available")
		http.Redirect(w, r, "/search-availability", http.StatusSeeOther)
		return
	}

	data := make(map[string]interface{})
	data["rooms"] = rooms

	res := models.Reservation{
		StartDate: start_date,
		EndDate:   end_date,
	}

	m.App.Session.Put(r.Context(), "reservation", res)

	render.Template(w, r, "choose-room.page.tmpl", &models.TemplateData{
		Data: data,
	})

}

type jsonResponse struct {
	OK        bool   `json:"ok"`
	Message   string `json:"message"`
	RoomID    string `json:"room_id"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

// finds room avilability from db by json request without refreshing page from room page
func (m *Repository) AvailabilityJSON(w http.ResponseWriter, r *http.Request) {

	err := r.ParseForm()
	if err != nil {

		// can't parse form, so return appropriate json
		resp := jsonResponse{
			OK:      false,
			Message: "Internal Server Error",
		}

		// write json back to browser
		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	start := r.Form.Get("start")
	end := r.Form.Get("end")

	layout := "2006-01-02" // YYYY-M-D

	start_date, err := time.Parse(layout, start)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	end_date, err := time.Parse(layout, end)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	roomID, err := strconv.Atoi(r.Form.Get("room_id"))
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse room id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	available, err := m.DB.SearchAvailabilityByRoomID(start_date, end_date, roomID)

	if err != nil {
		// m.App.Session.Put(r.Context(), "error", "can't get room availability from db")
		// http.Redirect(w, r, "/", http.StatusSeeOther)
		// return

		resp := jsonResponse{
			OK:      false,
			Message: "Internal Server Error",
		}

		// write json back to browser
		out, _ := json.MarshalIndent(resp, "", "     ")
		w.Header().Set("Content-Type", "application/json")
		w.Write(out)
		return
	}

	respo := jsonResponse{
		OK:        available,
		Message:   "",
		StartDate: start,
		EndDate:   end,
		RoomID:    strconv.Itoa(roomID),
	}
	out, _ := json.MarshalIndent(respo, "", "     ")

	// if err != nil {
	// 	m.App.Session.Put(r.Context(), "error", "can't parse json")
	// 	http.Redirect(w, r, "/", http.StatusSeeOther)
	// 	return
	// }
	w.Header().Set("Content-Type", "application/json")
	w.Write(out)

}

// displays list of available rooms
func (m *Repository) ChooseRoom(w http.ResponseWriter, r *http.Request) {
	// roomID, err := strconv.Atoi(chi.URLParam(r, "id"))

	exploded := strings.Split(r.RequestURI, "/")
	roomID, err := strconv.Atoi(exploded[2])
	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't get room id from url")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res, ok := m.App.Session.Get(r.Context(), "reservation").(models.Reservation)

	if !ok {

		m.App.Session.Put(r.Context(), "error", "can't get reservation from session")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	res.RoomID = roomID

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)

}

// takes url parameters and build sessional variable and takes user to reservation form
func (m *Repository) BookRoom(w http.ResponseWriter, r *http.Request) {

	roomID, err := strconv.Atoi(r.URL.Query().Get("id"))

	if err != nil {

		m.App.Session.Put(r.Context(), "error", "can't parse room id")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	sd := r.URL.Query().Get("s")
	ed := r.URL.Query().Get("e")

	layout := "2006-01-02" // YYYY-M-D

	start_date, err := time.Parse(layout, sd)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse start date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}
	end_date, err := time.Parse(layout, ed)
	if err != nil {
		m.App.Session.Put(r.Context(), "error", "can't parse end date")
		http.Redirect(w, r, "/", http.StatusSeeOther)
		return
	}

	var res models.Reservation

	res.RoomID = roomID
	res.StartDate = start_date
	res.EndDate = end_date

	m.App.Session.Put(r.Context(), "reservation", res)

	http.Redirect(w, r, "/make-reservation", http.StatusSeeOther)
}

func (m *Repository) Login(w http.ResponseWriter, r *http.Request) {

	render.Template(w, r, "login.page.tmpl", &models.TemplateData{
		Form: forms.New(nil),
	})
}

// logging the user in
func (m *Repository) PostLogin(w http.ResponseWriter, r *http.Request) {

	// prevents session fixation attack, every session has a certain token associated with it so it is good practice to renew token while login or logout

	_ = m.App.Session.RenewToken(r.Context())

	err := r.ParseForm()
	if err != nil {
		log.Println(err)
	}

	email := r.Form.Get("email")
	password := r.Form.Get("password")

	form := forms.New(r.PostForm)

	form.Required("email", "password")
	form.IsEmail("email")

	if !form.Valid() {
		render.Template(w, r, "login.page.tmpl", &models.TemplateData{
			Form: form,
		})
	}
	id, _, err := m.DB.Authenticate(email, password)
	if err != nil {
		log.Println(err)
		m.App.Session.Put(r.Context(), "error", "Invalid Credentials")
		http.Redirect(w, r, "/user/login", http.StatusSeeOther)
		return
	}
	log.Println("worked")
	m.App.Session.Put(r.Context(), "user_id", id)
	m.App.Session.Put(r.Context(), "flash", "Logged in successfully")
	http.Redirect(w, r, "/", http.StatusSeeOther)
}

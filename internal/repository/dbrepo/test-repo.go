package dbrepo

import (
	"errors"
	"time"

	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
)

func (m *testDBRepo) AllUsers() bool {

	return true
}

func (m *testDBRepo) InsertReservation(res models.Reservation) (int, error) {
	if res.RoomID == 3 {
		return 0, errors.New("insert res error")
	}
	return 1, nil

}

// insert room restriction
func (m *testDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	if r.RoomID == 200 {
		return errors.New("insert res error")
	}
	return nil
}

// gives bool if room is available based on given dates and room id
func (m *testDBRepo) SearchAvailabilityByRoomID(start, end time.Time, roomId int) (bool, error) {
	if roomId == 300 {
		return true, errors.New("SearchAvailabilityByRoomID error")
	}

	return false, nil
}

// gives slice of available rooms if any for given date range
func (m *testDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {

	var rooms []models.Room
	if end.Day() < start.Day() {
		return rooms, errors.New("some error")
	}

	return rooms, nil
}

// get room data by id
func (m *testDBRepo) GetRoomByID(id int) (models.Room, error) {

	var room models.Room

	if id == 100 {
		return room, errors.New("some test Error")
	}
	return room, nil
}

func (m *testDBRepo) GetUserByID(id int) (models.User, error) {
	var u models.User

	return u, nil
}
func (m *testDBRepo) UpdateUser(user models.User) error {

	return nil
}

func (m *testDBRepo) Authenticate(email, testPassword string) (int, string, error) {

	if email == "iam@me.com" {
		return 1, "", nil
	} else {
		return 0, "", errors.New("some error")
	}
}

func (m *testDBRepo) AdminAllReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

func (m *testDBRepo) AllNewReservations() ([]models.Reservation, error) {

	var reservations []models.Reservation

	return reservations, nil
}

func (m *testDBRepo) GetReservationByID(id int) (models.Reservation, error) {

	var i models.Reservation

	return i, nil
}

func (m *testDBRepo) UpdateReservationByID(user models.Reservation) error {
	return nil
}

func (m *testDBRepo) UpdateProcessedForReservation(id, processed int) error {

	return nil
}

func (m *testDBRepo) DeleteReservation(id int) error {

	return nil
}

func (m *testDBRepo) AllRooms() ([]models.Room, error) {
	var rooms []models.Room

	return rooms, nil
}

func (m *testDBRepo) GetRoomRestrictionByDate(id int, start, end time.Time) ([]models.RoomRestriction, error) {

	var roomRestriction []models.RoomRestriction

	return roomRestriction, nil

}

func (m *testDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {

	return nil
}

func (m *testDBRepo) DeleteBlockForRoom(id int) error {

	return nil
}

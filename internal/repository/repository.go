package repository

import (
	"time"

	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool
	InsertReservation(models.Reservation) (int, error)
	InsertRoomRestriction(models.RoomRestriction) error
	SearchAvailabilityByRoomID(start, end time.Time, roomId int) (bool, error)
	SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error)
	GetRoomByID(id int) (models.Room, error)
	GetUserByID(id int) (models.User, error)
	UpdateUser(user models.User) error
	Authenticate(email, testPassword string) (int, string, error)
	AdminAllReservations() ([]models.Reservation, error)
	AllNewReservations() ([]models.Reservation, error)
	GetReservationByID(id int) (models.Reservation, error)
	UpdateReservationByID(user models.Reservation) error
	UpdateProcessedForReservation(id, processed int) error
	DeleteReservation(id int) error
	AllRooms() ([]models.Room, error)
	GetRoomRestrictionByDate(id int, start, end time.Time) ([]models.RoomRestriction, error)
	InsertBlockForRoom(id int, startDate time.Time) error
	DeleteBlockForRoom(id int) error
}

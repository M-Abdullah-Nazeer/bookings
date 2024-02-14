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
}

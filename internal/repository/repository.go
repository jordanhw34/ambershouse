package repository

import (
	"time"

	"github.com/jordanhw34/ambershouse/internal/models"
)

type DatabaseRepo interface {
	AllUsers() bool

	InsertReservation(res models.Reservation) (int, error)
	InsertRoomRestriction(res models.RoomRestriction) error
	AvailableByRoomIDAndDates(startDate, endDate time.Time, roomID int) (bool, error)
	AvailableByDates(startDate, endDate time.Time) ([]models.Room, error)
}

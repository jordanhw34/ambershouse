package dbrepo

import (
	"errors"
	"time"

	"github.com/jordanhw34/ambershouse/internal/models"
)

// Dummy function to show how we can access in handlers
func (m *dbTestRepo) AllUsers() bool {
	return true
}

// InsertReservation => inserts a reservation into the database
func (m *dbTestRepo) InsertReservation(res models.Reservation) (int, error) {
	return 1, nil
}

// InsertRoomRestriction => inserts a room restriction into our database
func (m *dbTestRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	return nil
}

// AvailableByRoomIDAndDates returns bool and error given a RoomID and startDate and endDate
func (m *dbTestRepo) AvailableByRoomIDAndDates(startDate, endDate time.Time, roomID int) (bool, error) {
	return false, nil
}

// AvailableByDates returns a slice of available rooms for a given date range
func (m *dbTestRepo) AvailableByDates(startDate, endDate time.Time) ([]models.Room, error) {
	var rooms []models.Room
	return rooms, nil
}

// GetRoomByID returns a Room given an ID
func (m *dbTestRepo) GetRoomByID(id int) (models.Room, error) {
	var room models.Room
	if id > 2 {
		return room, errors.New("error getting room")
	}
	return room, nil
}

func (m *dbTestRepo) GetUserByID(id int) (models.User, error) {
	var user models.User
	return user, nil
}
func (m *dbTestRepo) UpdateUser(u models.User) error {
	return nil
}
func (m *dbTestRepo) Authenticate(email, testPassword string) (int, string, error) {
	return 1, "hashedPassword", nil
}

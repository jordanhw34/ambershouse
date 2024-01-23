package dbrepo

import (
	"context"
	"time"

	"github.com/jordanhw34/ambershouse/internal/models"
)

// Dummy function to show how we can access in handlers
func (m *dbPostresRepo) AllUsers() bool {
	return true
}

// InsertReservation => inserts a reservation into the database
func (m *dbPostresRepo) InsertReservation(res models.Reservation) (int, error) {

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Updating SQL with "returning id", will use QueryRowContext() instead of ExecContext()
	var newID int
	stmt := `insert into reservations (first_name, last_name, email, phone, start_date, end_date, room_id, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7, $8, $9) returning id`
	//_, err := m.DB.ExecContext(
	err := m.DB.QueryRowContext(
		ctx,
		stmt,
		res.FirstName,
		res.LastName,
		res.Email,
		res.Phone,
		res.StartDate,
		res.EndDate,
		res.RoomID,
		time.Now(),
		time.Now(),
	).Scan(&newID)

	if err != nil {
		return 0, err
	}
	return newID, nil
}

// InsertRoomRestriction => inserts a room restriction into our database
func (m *dbPostresRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, reservation_id, restriction_id, created_at, updated_at)
				values ($1, $2, $3, $4, $5, $6, $7)`

	_, err := m.DB.ExecContext(
		ctx,
		stmt,
		r.StartDate,
		r.EndDate,
		r.RoomID,
		r.ReservationID,
		r.RestrictionID,
		time.Now(),
		time.Now(),
	)
	if err != nil {
		return err
	}

	return nil
}

// AvailableByRoomIDAndDates returns bool and error given a RoomID and startDate and endDate
func (m *dbPostresRepo) AvailableByRoomIDAndDates(startDate, endDate time.Time, roomID int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var numRows int
	defer cancel()

	query := `
		select count(id) 
		from room_restrictions 
		where room_id = $1
			and $2 < end_date and $3 > start_date;
	`

	row := m.DB.QueryRowContext(ctx, query, roomID, startDate, endDate)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// AvailableByDates returns a slice of available rooms for a given date range
func (m *dbPostresRepo) AvailableByDates(startDate, endDate time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	var rooms []models.Room
	defer cancel()

	query := `
		select r.id, r.room_name
		from rooms r
		where r.id not in (select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)
	`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return nil, err
	}

	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return nil, err
		}
		rooms = append(rooms, room)
	}

	err = rows.Err()
	if err != nil {
		return nil, err
	}
	return rooms, nil
}

package dbrepo

import (
	"context"
	"time"

	"github.com/jordanhw34/ambershouse/internal/models"
	"golang.org/x/crypto/bcrypt"
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
	defer cancel()
	var numRows int

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
	defer cancel()
	var rooms []models.Room

	query := `
		select r.id, r.room_name
		from rooms r
		where r.id not in (select room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date)
	`

	rows, err := m.DB.QueryContext(ctx, query, startDate, endDate)
	if err != nil {
		return rooms, err
	}

	defer rows.Close()

	for rows.Next() {
		var room models.Room
		err := rows.Scan(&room.ID, &room.RoomName)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)
	}

	err = rows.Err()
	if err != nil {
		return rooms, err
	}
	return rooms, nil
}

// GetRoomByID returns a Room given an ID
func (m *dbPostresRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	var room models.Room

	query := `select id, room_name, created_at, updated_at from rooms where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)
	err := row.Scan(
		&room.ID,
		&room.RoomName,
		&room.CreatedAt,
		&room.UpdatedAt,
	)
	if err != nil {
		return room, err
	}

	return room, nil
}

// GetUserByID returns a User by ID
func (m *dbPostresRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `select id, first_name, last_name, email, password, access_level, created_at, updated_at from users where id = $1`

	row := m.DB.QueryRowContext(ctx, query, id)

	var u models.User
	err := row.Scan(
		&u.ID,
		&u.FirstName,
		&u.LastName,
		&u.Email,
		&u.Password,
		&u.AccessLevel,
		&u.CreatedAt,
		&u.UpdatedAt,
	)
	if err != nil {
		return u, err
	}
	return u, nil
}

// UpdateUser in the database given a User object
func (m *dbPostresRepo) UpdateUser(u models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	query := `update users set first_name = $1, last_name = $2, access_level = $3, updated_at = $4`

	_, err := m.DB.ExecContext(ctx, query, u.FirstName, u.LastName, u.AccessLevel, time.Now())
	if err != nil {
		return err
	}
	return nil
}

// Authenticate a user using bcrypt
func (m *dbPostresRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	row := m.DB.QueryRowContext(ctx, "select id, password from users where email = $1", email)
	err := row.Scan(&id, &hashedPassword)
	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return id, "", bcrypt.ErrMismatchedHashAndPassword
	} else if err != nil {
		return id, "", err
	}
	return id, hashedPassword, nil
}

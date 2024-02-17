package dbrepo

import (
	"context"
	"errors"
	"time"

	"github.com/M-Abdullah-Nazeer/bookings/internal/models"
	"golang.org/x/crypto/bcrypt"
)

func (m *postgreDBRepo) AllUsers() bool {

	return true
}

func (m *postgreDBRepo) InsertReservation(res models.Reservation) (int, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var newID int

	stmt := `insert into reservation (first_name,last_name,email,phone,start_date,end_date,room_id,created_at,updated_at) 
									values($1,$2,$3,$4,$5,$6,$7,$8,$9) returning id`
	err := m.DB.QueryRowContext(ctx, stmt,
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

// insert room restriction
func (m *postgreDBRepo) InsertRoomRestriction(r models.RoomRestriction) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date,end_date,room_id,reservation_id,restriction_id,created_at,updated_at)
	values
	($1,$2,$3,$4,$5,$6,$7)`

	_, err := m.DB.ExecContext(ctx, stmt,
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

// gives bool if room is available based on given dates and room id
func (m *postgreDBRepo) SearchAvailabilityByRoomID(start, end time.Time, roomId int) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var numRows int
	stmt := `select count(id) from room_restrictions 
	where room_id = $1 and
	 $2 < end_date and $3 > start_date;`

	row := m.DB.QueryRowContext(ctx, stmt, roomId, start, end)
	err := row.Scan(&numRows)
	if err != nil {
		return false, err
	}

	if numRows == 0 {
		return true, nil
	}

	return false, nil
}

// gives slice of available rooms if any for given date range
func (m *postgreDBRepo) SearchAvailabilityForAllRooms(start, end time.Time) ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select r.id, r.room_name 
			from room r 
			where r.id NOT IN (select rr.room_id from room_restrictions rr where $1 < rr.end_date and $2 > rr.start_date);`

	rows, err := m.DB.QueryContext(ctx, stmt, start, end)

	var rooms []models.Room
	if err != nil {
		return rooms, err
	}

	for rows.Next() {
		var room models.Room

		err := rows.Scan(
			&room.ID,
			&room.RoomName,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, room)

	}

	if err = rows.Err(); err != nil {
		return rooms, err
	}
	return rooms, nil
}

// get room data by id
func (m *postgreDBRepo) GetRoomByID(id int) (models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var room models.Room
	stmt := `select id,room_name,created_at,updated_at from room where id=$1`

	row := m.DB.QueryRowContext(ctx, stmt, id)
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

func (m *postgreDBRepo) GetUserByID(id int) (models.User, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select id,first_name,last_name,email,password,access_level,created_at,updated_at 
	from users where id=$1`

	row := m.DB.QueryRowContext(ctx, stmt, id)

	var user models.User

	err := row.Scan(
		&user.ID,
		&user.FirstName,
		&user.LastName,
		&user.Email,
		&user.Password,
		&user.AccessLevel,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if err != nil {
		return user, err
	}

	return user, nil
}

// updates user
func (m *postgreDBRepo) UpdateUser(user models.User) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := ` update users set first_name=$1, 
	last_name=$2,email=$3,access_level=$4,updated_at=$5
	`

	_, err := m.DB.ExecContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Email,
		user.AccessLevel,
		time.Now(),
	)

	if err != nil {
		return err
	}

	return nil
}

// matches given password hash with hash in db
func (m *postgreDBRepo) Authenticate(email, testPassword string) (int, string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var id int
	var hashedPassword string

	stmt := `select id, password from users where email=$1`
	row := m.DB.QueryRowContext(ctx, stmt, email)
	err := row.Scan(&id, &hashedPassword)

	if err != nil {
		return id, "", err
	}

	err = bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(testPassword))

	// built in error of bcrypt in standard library
	if err == bcrypt.ErrMismatchedHashAndPassword {
		return 0, "", errors.New("incorrect password")
	} else if err != nil {
		return 0, "", err
	}

	return id, hashedPassword, nil
}

package dbrepo

import (
	"context"
	"errors"
	"log"
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

// returns slice of all reservation to admin
func (m *postgreDBRepo) AdminAllReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	stmt := `select r.id, r.first_name, r.last_name, r.email, r.phone, 
	r.start_date, r.end_date, r.room_id, rm.room_name, r.processed from reservation r
	left join room rm on (r.room_id = rm.id) 
	order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, stmt)

	if err != nil {
		return reservations, err
	}

	for rows.Next() {

		var i models.Reservation

		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.Room.RoomName,
			&i.Processed,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	defer rows.Close()

	return reservations, nil
}

// returns slice of all new reservation to admin
func (m *postgreDBRepo) AllNewReservations() ([]models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	var reservations []models.Reservation

	stmt := `select r.id, r.first_name, r.last_name, r.email, r.phone, 
	r.start_date, r.end_date, r.room_id, rm.room_name, r.processed from reservation r
	left join room rm on (r.room_id = rm.id) 
	where r.processed = 0
	order by r.start_date asc
	`

	rows, err := m.DB.QueryContext(ctx, stmt)

	if err != nil {
		return reservations, err
	}

	for rows.Next() {

		var i models.Reservation

		err = rows.Scan(
			&i.ID,
			&i.FirstName,
			&i.LastName,
			&i.Email,
			&i.Phone,
			&i.StartDate,
			&i.EndDate,
			&i.RoomID,
			&i.Room.RoomName,
			&i.Processed,
		)
		if err != nil {
			return reservations, err
		}
		reservations = append(reservations, i)
	}

	if err = rows.Err(); err != nil {
		return reservations, err
	}

	defer rows.Close()

	return reservations, nil
}

func (m *postgreDBRepo) GetReservationByID(id int) (models.Reservation, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select r.id, r.first_name, r.last_name, r.email, r.phone, 
	r.start_date, r.end_date, r.room_id, rm.room_name, r.processed from reservation r
	left join room rm on (r.room_id = rm.id) 
	where r.id = $1
	`

	rows := m.DB.QueryRowContext(ctx, stmt, id)

	var i models.Reservation

	err := rows.Scan(
		&i.ID,
		&i.FirstName,
		&i.LastName,
		&i.Email,
		&i.Phone,
		&i.StartDate,
		&i.EndDate,
		&i.RoomID,
		&i.Room.RoomName,
		&i.Processed,
	)

	if err != nil {
		return i, err
	}

	return i, nil
}

func (m *postgreDBRepo) UpdateReservationByID(user models.Reservation) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := ` update reservation set first_name=$1, 
	last_name=$2, email=$3, phone=$4, updated_at=$5
	where id=$6`

	_, err := m.DB.ExecContext(ctx, stmt,
		user.FirstName,
		user.LastName,
		user.Email,
		user.Phone,
		time.Now(),
		user.ID,
	)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgreDBRepo) UpdateProcessedForReservation(id, processed int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := ` update reservation set processed=$1 where id =$2`

	_, err := m.DB.ExecContext(ctx, stmt, processed, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgreDBRepo) DeleteReservation(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from reservation where id=$1`

	_, err := m.DB.ExecContext(ctx, stmt, id)

	if err != nil {
		return err
	}

	return nil
}

func (m *postgreDBRepo) AllRooms() ([]models.Room, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `select id, room_name, created_at, updated_at from room`
	var rooms []models.Room

	rows, err := m.DB.QueryContext(ctx, stmt)
	if err != nil {
		return rooms, err
	}

	for rows.Next() {

		var r models.Room

		err := rows.Scan(
			&r.ID,
			&r.RoomName,
			&r.CreatedAt,
			&r.UpdatedAt,
		)
		if err != nil {
			return rooms, err
		}
		rooms = append(rooms, r)
	}
	if err := rows.Err(); err != nil {
		return rooms, err
	}
	defer rows.Close()

	return rooms, nil
}

func (m *postgreDBRepo) GetRoomRestrictionByDate(id int, start, end time.Time) ([]models.RoomRestriction, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	// coalesce will replace null with 0
	stmt := `select id, start_date, end_date, room_id, coalesce(reservation_id,0), restriction_id 
	from room_restrictions where $1 < end_date and $2 >= start_date and room_id = $3 `

	var roomRestriction []models.RoomRestriction

	rows, err := m.DB.QueryContext(ctx, stmt, start, end, id)
	if err != nil {
		return roomRestriction, err
	}

	for rows.Next() {
		var rr models.RoomRestriction

		err := rows.Scan(
			&rr.ID,
			&rr.StartDate,
			&rr.EndDate,
			&rr.RoomID,
			&rr.ReservationID,
			&rr.RestrictionID,
		)
		if err != nil {
			return roomRestriction, err
		}

		roomRestriction = append(roomRestriction, rr)

	}

	if err := rows.Err(); err != nil {
		return roomRestriction, err
	}
	defer rows.Close()

	return roomRestriction, nil

}

// inserts room restriction
func (m *postgreDBRepo) InsertBlockForRoom(id int, startDate time.Time) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `insert into room_restrictions (start_date, end_date, room_id, restriction_id, created_at, updated_at) 
			values ($1, $2, $3, $4, $5, $6)`

	endDate := startDate.AddDate(0, 0, 0)

	_, err := m.DB.ExecContext(ctx, stmt, startDate, endDate, id, 2, time.Now(), time.Now())
	if err != nil {
		log.Println("yha ka error h ---------------")
		return err
	}

	return nil
}

// deletes room restrictions
func (m *postgreDBRepo) DeleteBlockForRoom(id int) error {
	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	stmt := `delete from room_restrictions where id=$1`

	_, err := m.DB.ExecContext(ctx, stmt, id)
	if err != nil {

		return err
	}

	return nil
}

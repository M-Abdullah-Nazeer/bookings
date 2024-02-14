package driver

import (
	"database/sql"
	"time"

	_ "github.com/jackc/pgx/v5/stdlib"
	// "github.com/jackc/pgx/v5/pgconn"
	// "github.com/jackc/pgx/v5/stdlib"
)

type DB struct {
	SQL *sql.DB
}

var dbCon = &DB{}

const maxOpenDbCon = 10
const maxIdleConn = 5
const maxDbLifetime = 5 * time.Minute

// ConnectSQL creates database connection pool

func ConnectSQL(dsn string) (*DB, error) {

	d, err := NewDatabase(dsn)
	if err != nil {
		panic(err)
	}
	d.SetConnMaxIdleTime(maxIdleConn)
	d.SetConnMaxLifetime(maxDbLifetime)
	d.SetMaxOpenConns(maxOpenDbCon)

	dbCon.SQL = d

	err = testDB(d)
	if err != nil {
		return nil, err
	}

	return dbCon, nil
}

// testDB tries to ping db

func testDB(d *sql.DB) error {

	err := d.Ping()
	if err != nil {
		return err
	}
	return nil
}

// NewDatabase creates new db for application
func NewDatabase(dsn string) (*sql.DB, error) {

	db, err := sql.Open("pgx", dsn)
	// db,err :=pgx.Connect(context.Background(),os.Getenv(dsn))
	if err != nil {
		return nil, err
	}
	if err = db.Ping(); err != nil {
		return nil, err
	}
	return db, nil
}

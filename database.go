package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

func NewDatabase() (*Database, error) {
	// db, err := sql.Open("sqlite3", "./db.sqlite3?cache=shared&mode=rwc&_busy_timeout=50000000")
	db, err := sql.Open("sqlite3", "./db.sqlite3")
	database := Database{db: db}
	if nil != err {
		return &database, err
	}

	err = database.Exec(`
		CREATE TABLE IF NOT EXISTS history (
			trip_id             TEXT,
			event_timestamp     TIMESTAMP,
			longitude           DOUBLE PRECISION,
			latitude            DOUBLE PRECISION,
			speed				REAL,
			heading				REAL
		);
	`)
	if nil != err {
		return &database, err
	}

	// err = database.Exec("PRAGMA journal_mode=WAL;")
	// if nil != err {
	// 	return &database, err
	// }

	return &database, err
}

type Database struct {
	db *sql.DB
}

func (self *Database) Close() error {
	return self.db.Close()
}

func (self *Database) Exec(query string, args ...interface{}) error {
	statement, err := self.db.Prepare(query)
	if nil != err {
		return err
	}
	_, err = statement.Exec(args...)
	return err
}

func (self *Database) Query(query string, args ...interface{}) (*sql.Rows, error) {
	return self.db.Query(query, args...)
}

func (self *Database) InsertEvent(event *Event) error {
	return self.Exec(`
		INSERT INTO history (trip_id, event_timestamp, longitude, latitude, speed, heading)
			VALUES (?, ?, ?, ?, ?, ?)`,
		event.TripID,
		event.Timestamp,
		event.Longitude,
		event.Latitude,
		event.Speed,
		event.Heading,
	)
}

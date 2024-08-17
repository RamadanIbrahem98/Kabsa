package db

import (
	"database/sql"
	"log"
	"os"
	"time"

	_ "github.com/mattn/go-sqlite3"
)

type DB struct {
	db *sql.DB
}

func New() (*DB, error) {
	dbPath := os.Getenv("DB_PATH")
	if dbPath == "" {
		dbPath = "./data/kabsa.db"
	}
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, err
	}

	database := &DB{
		db: db,
	}

	database.prepare()

	return database, nil
}

func (d *DB) prepare() {
	sqlStmt := `
	create table if not exists kabsa (id integer not null primary key, presses integer, start_at integer, end_at integer, wpm integer, created_at datetime);
	`

	_, err := d.db.Exec(sqlStmt)
	if err != nil {
		log.Printf("%q: %s\n", err, sqlStmt)
		return
	}
}

func (d *DB) Insert(presses, startAt, endAt, wpm int64) {
	stmt, err := d.db.Prepare("INSERT INTO kabsa(presses, start_at, end_at, wpm, created_at) values(?, ?, ?, ?, ?)")
	if err != nil {
		log.Fatal(err)
	}
	defer stmt.Close()

	_, err = stmt.Exec(presses, startAt, endAt, wpm, time.Now())
	if err != nil {
		log.Fatal(err)
	}
}

func (d *DB) Close() {
	d.db.Close()
}

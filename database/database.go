package database

import (
	"crypto/sha256"
	"database/sql"
	"errors"
	"sync"

	_ "github.com/mattn/go-sqlite3"
)

// reference to database
type Database struct {
	mu sync.Mutex
	db *sql.DB
}

type UserInfoEntry struct {
	ID int
	Username string
	PasswordHash string
}

// add row to database
func (db *Database) Insert(entry UserInfoEntry) (error) {
	passwordHash := sha256.Sum256([]byte(entry.PasswordHash))
	_, err := db.db.Exec("INSERT INTO userinfo VALUES(NULL, ?, ?);", entry.Username, string(passwordHash[:]))
	if err != nil {
		return err
	}

	return nil
}

// get row from database
func (db *Database) GetEntryByUsername(username string) (UserInfoEntry, error) {
	row := db.db.QueryRow("SELECT * FROM userinfo WHERE username=?", username)

	entry := UserInfoEntry{}

	if err := row.Scan(&entry.ID, &entry.Username, &entry.PasswordHash); err != nil {
		return UserInfoEntry{}, errors.New("No Entry Found") // if no entry is found return an error
	}

	return entry, nil
}

// string for creating table
// in database if it does not exist
const createString = `
	CREATE TABLE IF NOT EXISTS userinfo (
		id INTEGER NOT NULL PRIMARY KEY,
		username varchar(64),
		passwordHash varchar(64)
	);`

// initialise database instance
func InitDB() (*Database, error){
	db, err := sql.Open("sqlite3", "database.db")
	if err != nil {
		return nil, err
	}

	if _, err := db.Exec(createString); err != nil {
		return nil, err
	}

	return &Database{
		db: db,
	}, nil
}
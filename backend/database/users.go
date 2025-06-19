package database

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

// ==== The function will create the users table in the database =====
func CreateUsersTable(db *sql.DB) error {
	query := `
  CREATE TABLE IF NOT EXISTS users (
    id INTEGER PRIMARY KEY AUTOINCREMENT,
    uuid TEXT UNIQUE NOT NULL,
    nickname TEXT UNIQUE NOT NULL,
    age INTEGER NOT NULL,
    gender TEXT NOT NULL,
    email TEXT UNIQUE NOT NULL,
    first_name TEXT UNIQUE NOT NULL,
    last_name TEXT NOT NULL,
    password TEXT NOT NULL,
    last_message_time DATETIME
    );`
	if _, err := db.Exec(query); err != nil {
		return err
	}
	return nil
}

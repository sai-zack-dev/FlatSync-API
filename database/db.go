package database

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func ConnectDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./flatsync.db")
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}

	createUserTable()
}

func createUserTable() {
	createTable := `
    CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        name TEXT,             -- Name is nullable
        dob TEXT,              -- Date of birth is nullable
        avatar TEXT,           -- Avatar URL is nullable
        email TEXT NOT NULL UNIQUE,  -- Email is required (NOT NULL)
        password TEXT NOT NULL  -- Password is required (NOT NULL)
    );`
	if _, err := DB.Exec(createTable); err != nil {
		log.Fatal("Failed to create users table:", err)
	}
}

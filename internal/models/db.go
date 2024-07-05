package models

import (
	"database/sql"
	"log"

	_ "github.com/mattn/go-sqlite3"
)

var DB *sql.DB

func InitDB() {
	var err error
	DB, err = sql.Open("sqlite3", "./blog.db")
	if err != nil {
		log.Fatal(err)
	}

	createTable := `CREATE TABLE IF NOT EXISTS posts (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		title TEXT,
		content TEXT,
		image TEXT,
		user_id INTEGER,
		created_at DATETIME DEFAULT CURRENT_TIMESTAMP,
		FOREIGN KEY(user_id) REFERENCES users(id)
	);`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}
	createTable = `CREATE TABLE IF NOT EXISTS users (
        id INTEGER PRIMARY KEY AUTOINCREMENT,
        Username TEXT,
        Email TEXT,
		Password TEXT,
		Icon TEXT,
        created_at DATETIME DEFAULT CURRENT_TIMESTAMP
    );`

	_, err = DB.Exec(createTable)
	if err != nil {
		log.Fatal(err)
	}

}

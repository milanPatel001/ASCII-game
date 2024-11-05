package handlers

import (
	"database/sql"
	"log"
)

var db *sql.DB

func GetDB() *sql.DB {
	if db == nil {
		var err error
		db, err = sql.Open("sqlite3", "./game.db")
		if err != nil {
			log.Fatalf("Failed to connect to database: %v", err)
		}

		if err := db.Ping(); err != nil {
			log.Fatalf("Failed to ping database: %v", err)
		}
		log.Println("Sqlite database connection established")
	}

	return db
}

func CloseDB() {
	if db != nil {
		db.Close()
	}
}

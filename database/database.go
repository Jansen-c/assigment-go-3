package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func InitDB(connectionString string) (*sql.DB, error) { //? second () is the expected return. * inside is the address on memory.
	// Open database
	// if db, err := sql.Open("postgres", connectionString); err != nil { //? bisa dibikin mirip ternary style gitu, tapi makin susah diliat jadinya. PLUS gabisa dipake dibawah.
	// 	return nil, err
	// } else {
	// 	return db, nil
	// }
	db, err := sql.Open("postgres", connectionString)
	if err != nil {
		return nil, err
	}

	// Test connection
	err = db.Ping()
	if err != nil {
		return nil, err
	}

	// Set connection pool settings (optional tapi recommended)
	db.SetMaxOpenConns(25)
	db.SetMaxIdleConns(5)

	log.Println("Database connected successfully")
	return db, nil
}

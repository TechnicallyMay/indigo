package db

import (
	"database/sql"
	"log"

	"github.com/TechnicallyMay/indigo/appsettings"
	_ "modernc.org/sqlite"
)

var pool *sql.DB

func OpenDb(opts appsettings.AppSettings) *sql.DB {
	log.Println("Opening a connection to the database.")
	if pool != nil {
		log.Println("Connection already open, using existing connection.")
		return pool
	}

	var err error
	pool, err = sql.Open("sqlite", opts.DbPath)

	if err != nil {
		log.Fatal("Error when connecting to database.", err)
	}

	log.Println("Successfully connected to the database.")

	pool.SetMaxOpenConns(1)
	pool.SetConnMaxLifetime(0)

	err = pool.Ping()
	if err != nil {
		log.Fatal("Error when pinging database.", err)
	}
	log.Println("Successfully pinged db!")

	return pool
}

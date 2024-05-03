package db

import (
	"database/sql"
	"log"

	_ "modernc.org/sqlite"
)

var pool *sql.DB

func OpenDb() *sql.DB {
    log.Println("Opening a connection to the database.")
    if pool != nil {
        log.Println("Connection already open, using existing connection.")
        return pool
    }

    var err error
    pool, err = sql.Open("sqlite", ":memory:")

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


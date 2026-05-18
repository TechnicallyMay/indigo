package db

import (
	"database/sql"
	"log/slog"

	"github.com/TechnicallyMay/indigo/appsettings"
	_ "modernc.org/sqlite"
)

var pool *sql.DB

func OpenDb(opts appsettings.AppSettings) *sql.DB {
	slog.Info("Opening a connection to the database", "path", opts.DbPath)
	if pool != nil {
		slog.Info("Connection already open, using existing connection.")
		return pool
	}

	var err error
	pool, err = sql.Open("sqlite", opts.DbPath)

	if err != nil {
		slog.Error("Error when connecting to database.", "error", err)
		panic(err)
	}

	slog.Info("Successfully connected to the database.")

	pool.SetMaxOpenConns(1)
	pool.SetConnMaxLifetime(0)

	err = pool.Ping()
	if err != nil {
		slog.Error("Error when pinging database.", "error", err)
		panic(err)
	}
	slog.Info("Successfully pinged db!")

	return pool
}

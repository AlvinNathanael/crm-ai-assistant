package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func New(dsn string) (*sql.DB, error) {
	log.Println("Connecting to Postgres with DSN:", dsn)
	return sql.Open("postgres", dsn)
}

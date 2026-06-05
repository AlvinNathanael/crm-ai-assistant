package database

import "database/sql"

func New(dsn string) (*sql.DB, error) {
	return sql.Open("postgres", dsn)
}

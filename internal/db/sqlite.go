package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"
)

// NewSQLite
func NewSQLite() (*sql.DB, error) {

	dsn := os.Getenv("SQLITE_DSN")
	if dsn == "" {
		return nil, errors.New("$SQLITE_DSN is unset")
	}

	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("error opening sqlite database %v", err)
	}

	return db, nil
}

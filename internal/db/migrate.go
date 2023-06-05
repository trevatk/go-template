package db

import (
	"database/sql"
	"errors"
	"fmt"
	"os"

	"github.com/golang-migrate/migrate/v4"
	"github.com/golang-migrate/migrate/v4/database/sqlite"

	// file driver
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

// Migrate go-migrate database migration from file location
func Migrate(db *sql.DB) error {

	migrationDir := os.Getenv("SQLITE_MIGRATIONS_DIR")
	if migrationDir == "" {
		return errors.New("$SQLITE_MIGRATIONS_DIR is unset")
	}

	driver, err := sqlite.WithInstance(db, &sqlite.Config{})
	if err != nil {
		return fmt.Errorf("unable to create migration driver %v", err)
	}

	m, err := migrate.NewWithDatabaseInstance(
		fmt.Sprintf("file://%s", migrationDir),
		"sqlite", driver)
	if err != nil {
		return fmt.Errorf("failed to build database migration %v", err)
	}

	err = m.Up()
	if err != nil && !errors.Is(err, migrate.ErrNoChange) {
		return fmt.Errorf("failed to migrate database %v", err)
	}

	return nil
}

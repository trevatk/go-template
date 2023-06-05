package db_test

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/trevatk/go-template/internal/db"
)

func init() {
	os.Setenv("SQLITE_DSN", "./testfiles/sqlite/persons.db")
}

func TestNewSQLite(t *testing.T) {

	assert := assert.New(t)

	db, err := db.NewSQLite()
	assert.NoError(err)

	db.Close()
}

// Code generated by sqlc. DO NOT EDIT.
// versions:
//   sqlc v1.18.0

package persons

import (
	"database/sql"
	"time"
)

type Person struct {
	ID        int64
	Fname     string
	Lname     string
	Email     string
	CreatedAt time.Time
	UpdatedAt sql.NullTime
}
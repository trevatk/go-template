package domain

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/trevatk/go-template/internal/repository/persons"
)

// NewPerson
type NewPerson struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// Person
type Person struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdatePerson
type UpdatePerson struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// PersonService
type PersonService struct {
	queries *persons.Queries
}

// NewPersonService
func NewPersonService(ctx context.Context, db *sql.DB) (*PersonService, error) {

	conn, err := db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("unable to get database connection from pool %v", err)
	}

	return &PersonService{
		queries: persons.New(conn),
	}, nil
}

// CreatePerson
func (ps *PersonService) CreatePerson(ctx context.Context, newPerson *NewPerson) (*Person, error) {

	sqlPerson, err := ps.queries.InsertPerson(ctx, &persons.InsertPersonParams{
		Fname: newPerson.FirstName,
		Lname: newPerson.LastName,
		Email: newPerson.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert new person %v", err)
	}

	return &Person{
		ID:        sqlPerson.ID,
		FirstName: sqlPerson.Fname,
		LastName:  sqlPerson.Lname,
		Email:     sqlPerson.Email,
		CreatedAt: sqlPerson.CreatedAt,
		UpdatedAt: time.Time{},
	}, nil
}

// ReadPerson
func (ps *PersonService) ReadPerson(ctx context.Context, id int64) (*Person, error) {

	sqlPerson, err := ps.queries.ReadPerson(ctx, id)
	if err != nil {
		return nil, fmt.Errorf("error executing read person query %v", err)
	}

	return transformSqlPerson(sqlPerson), nil
}

// UpdatePerson
func (ps *PersonService) UpdatePerson(ctx context.Context, updatePerson *UpdatePerson) (*Person, error) {

	sqlPerson, err := ps.queries.UpdatePerson(ctx, &persons.UpdatePersonParams{
		Fname: updatePerson.FirstName,
		Lname: updatePerson.LastName,
		Email: updatePerson.Email,
		ID:    updatePerson.ID,
	})
	if err != nil {
		return nil, fmt.Errorf("error executing update person query %v", err)
	}

	return transformSqlPerson(sqlPerson), nil
}

func transformSqlPerson(sqlPerson *persons.Person) *Person {

	var person Person

	person.FirstName = sqlPerson.Fname
	person.LastName = sqlPerson.Lname
	person.Email = sqlPerson.Email
	person.CreatedAt = sqlPerson.CreatedAt

	if sqlPerson.UpdatedAt.Valid {
		person.UpdatedAt = sqlPerson.UpdatedAt.Time
	} else {
		person.UpdatedAt = time.Time{}
	}

	return &person
}

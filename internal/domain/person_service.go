package domain

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/trevatk/go-template/internal/repository/persons"
)

var (
	// ErrNotFound service level error message when resource is not found
	ErrNotFound = errors.New("resource id not found")
)

// NewPerson application layer model
type NewPerson struct {
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// NewPersonRequest request layer model used for validation of requests
// using chi render bind
type NewPersonRequest struct {
	*NewPerson
}

// Bind callback used to validate new person request model
func (npr *NewPersonRequest) Bind(r *http.Request) error {

	if npr.NewPerson == nil {
		return errors.New("no person details provided")
	}

	if npr.NewPerson.FirstName == "" {
		return errors.New("no first name provided")
	} else if npr.NewPerson.LastName == "" {
		return errors.New("no last name provided")
	} else if npr.NewPerson.Email == "" {
		return errors.New("no email provided")
	}

	return nil
}

// Person application layer model
type Person struct {
	ID        int64     `json:"id"`
	FirstName string    `json:"first_name"`
	LastName  string    `json:"last_name"`
	Email     string    `json:"email"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// UpdatePerson application layer model
type UpdatePerson struct {
	ID        int64  `json:"id"`
	FirstName string `json:"first_name"`
	LastName  string `json:"last_name"`
	Email     string `json:"email"`
}

// UpdatePersonRequest request layer model used for validation of requests
// using chi render bind
type UpdatePersonRequest struct {
	*UpdatePerson
}

// Bind callback used to validate update user request model
func (upr *UpdatePersonRequest) Bind(r *http.Request) error {

	if upr.UpdatePerson == nil {
		return errors.New("invalid request object")
	}

	return nil
}

// PersonService application layer to facilitate calls to business layer for all person related models
type PersonService struct {
	db *sql.DB
}

// NewPersonService create new person service instance
func NewPersonService(db *sql.DB) *PersonService {
	return &PersonService{
		db: db,
	}
}

// CreatePerson insert new person into database
func (ps *PersonService) Create(ctx context.Context, newPerson *NewPerson) (*Person, error) {

	conn, err := ps.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection %v", err)
	}
	defer conn.Close()

	sqlPerson, err := persons.New(conn).InsertPerson(ctx, &persons.InsertPersonParams{
		Fname: newPerson.FirstName,
		Lname: newPerson.LastName,
		Email: newPerson.Email,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to insert new person %v", err)
	}

	return transformSqlPerson(sqlPerson), nil
}

// ReadPerson retrieve person by id
func (ps *PersonService) Read(ctx context.Context, id int64) (*Person, error) {

	conn, err := ps.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection %v", err)
	}
	defer conn.Close()

	sqlPerson, err := persons.New(conn).ReadPerson(ctx, id)
	if err != nil {

		if err == sql.ErrNoRows {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("error executing read person query %v", err)
	}

	return transformSqlPerson(sqlPerson), nil
}

// UpdatePerson update existing person record
func (ps *PersonService) Update(ctx context.Context, updatePerson *UpdatePerson) (*Person, error) {

	conn, err := ps.db.Conn(ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to get database connection %v", err)
	}
	defer conn.Close()

	sqlPerson, err := persons.New(conn).UpdatePerson(ctx, &persons.UpdatePersonParams{
		Fname: updatePerson.FirstName,
		Lname: updatePerson.LastName,
		Email: updatePerson.Email,
		ID:    updatePerson.ID,
	})
	if err != nil {

		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrNotFound
		}

		return nil, fmt.Errorf("error executing update person query %v", err)
	}

	return transformSqlPerson(sqlPerson), nil
}

// Delete hard delete person record
func (ps *PersonService) Delete(ctx context.Context, id int64) error {

	conn, err := ps.db.Conn(ctx)
	if err != nil {
		return fmt.Errorf("failed to get database connection %v", err)
	}
	defer conn.Close()

	result, err := persons.New(conn).DeletePerson(ctx, id)
	if err != nil {
		return fmt.Errorf("error excuting delete person query %v", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("failed to get rows affected %v", err)
	}

	if affected == 0 {
		return ErrNotFound
	}

	return nil
}

// transform business model into application model
func transformSqlPerson(sqlPerson *persons.Person) *Person {

	var person Person

	person.ID = sqlPerson.ID
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

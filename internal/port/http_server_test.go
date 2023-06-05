package port

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/go-chi/chi/v5"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"github.com/trevatk/go-template/internal/db"
	"github.com/trevatk/go-template/internal/domain"
	"github.com/trevatk/go-template/internal/logging"
)

var (
	userId       int64 = 0
	readUserId   int64 = 0
	deleteUserId int64 = 0
)

func init() {
	os.Setenv("SQLITE_DSN", "./testfiles/sqlite/person.db")
	os.Setenv("SQLITE_MIGRATIONS_DIR", "./../../migrations")
}

type HttpServerSuite struct {
	suite.Suite
	mux *chi.Mux
}

func (suite *HttpServerSuite) SetupTest() {

	ctx := context.TODO()

	assert := assert.New(suite.T())

	logger, err := logging.New()
	assert.NoError(err)

	sqlite, err := db.NewSQLite()
	assert.NoError(err)

	err = db.Migrate(sqlite)
	assert.NoError(err)

	personService := domain.NewPersonService(sqlite)

	// preload database with two users
	readUser, err := personService.Create(ctx, &domain.NewPerson{
		FirstName: "read",
		LastName:  "person",
		Email:     "read.person@mailbox.com",
	})
	assert.NoError(err)
	readUserId = readUser.ID

	deletePerson, err := personService.Create(ctx, &domain.NewPerson{
		FirstName: "delete",
		LastName:  "person",
		Email:     "delete.person@mailbox.com",
	})
	assert.NoError(err)
	deleteUserId = deletePerson.ID

	bundle := domain.NewBundle(personService)

	server := NewHttpServer(logger, bundle)

	suite.mux = NewRouter(server)
}

func (suite *HttpServerSuite) TestCreatePerson() {

	assert := assert.New(suite.T())

	cases := []struct {
		newPerson *domain.NewPerson
		expected  int
	}{
		{
			// success
			newPerson: &domain.NewPerson{
				FirstName: "unit",
				LastName:  "test",
				Email:     "testing@mailbox.com",
			},
			expected: http.StatusCreated,
		},
		{
			// invalid body
			newPerson: &domain.NewPerson{
				FirstName: "", // intentionally leave empty
				LastName:  "test",
				Email:     "testing@mailbox.com",
			},
			expected: http.StatusBadRequest,
		},
	}

	for _, c := range cases {

		newPersonBytes, err := json.Marshal(c.newPerson)
		assert.NoError(err)

		req, err := http.NewRequest(http.MethodPost, "/api/v1/person/", bytes.NewBuffer(newPersonBytes))
		assert.NoError(err)

		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()
		suite.mux.ServeHTTP(rr, req)

		assert.Equal(c.expected, rr.Code)
	}
}

func (suite *HttpServerSuite) TestFetchPerson() {

	assert := assert.New(suite.T())

	cases := []struct {
		expected int
		endpoint string
	}{
		{
			// success
			expected: http.StatusAccepted,
			endpoint: fmt.Sprintf("/api/v1/person/%d", readUserId),
		},
		{
			// invalid URL parameter
			expected: http.StatusBadRequest,
			endpoint: "/api/v1/person/xyx",
		},
		{
			// not found
			expected: http.StatusNotFound,
			endpoint: fmt.Sprintf("/api/v1/person/%d", 99),
		},
	}

	for _, c := range cases {

		req, err := http.NewRequest(http.MethodGet, c.endpoint, nil)
		assert.NoError(err)

		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.mux.ServeHTTP(rr, req)

		assert.Equal(c.expected, rr.Code)
	}
}

func (suite *HttpServerSuite) TestUpdatePerson() {

	assert := assert.New(suite.T())

	cases := []struct {
		person   *domain.UpdatePersonRequest
		expected int
	}{
		{
			// success
			person: &domain.UpdatePersonRequest{
				UpdatePerson: &domain.UpdatePerson{
					ID:        readUserId,
					FirstName: "john",
					LastName:  "doe",
					Email:     "john.doe@mailbox.com",
				},
			},
			expected: http.StatusAccepted,
		},
		{
			// invalid request body
			person: &domain.UpdatePersonRequest{
				// intentioally leave nil
			},
			expected: http.StatusBadRequest,
		},
		{
			// not found
			person: &domain.UpdatePersonRequest{
				UpdatePerson: &domain.UpdatePerson{
					ID:        readUserId + 100,
					FirstName: "not",
					LastName:  "found",
					Email:     "not.found@mailbox.com",
				},
			},
			expected: http.StatusNotFound,
		},
	}

	for _, c := range cases {

		updatePersonByes, err := json.Marshal(c.person)
		assert.NoError(err)

		req, err := http.NewRequest(http.MethodPut, "/api/v1/person/", bytes.NewReader(updatePersonByes))
		assert.NoError(err)

		rr := httptest.NewRecorder()

		suite.mux.ServeHTTP(rr, req)

		assert.Equal(c.expected, rr.Code)
	}
}

func (suite *HttpServerSuite) TestDeletePerson() {

	assert := assert.New(suite.T())

	cases := []struct {
		expected int
		endpoint string
	}{
		{
			// success
			expected: http.StatusAccepted,
			endpoint: fmt.Sprintf("/api/v1/person/%d", deleteUserId),
		},
		{
			// invalid URL parameter
			expected: http.StatusBadRequest,
			endpoint: "/api/v1/person/xxx",
		},
		{
			// not found
			expected: http.StatusNotFound,
			endpoint: fmt.Sprintf("/api/v1/person/%d", deleteUserId+999),
		},
	}

	for _, c := range cases {

		req, err := http.NewRequest(http.MethodDelete, c.endpoint, nil)
		assert.NoError(err)

		req.Header.Add("Content-Type", "application/json")

		rr := httptest.NewRecorder()

		suite.mux.ServeHTTP(rr, req)

		assert.Equal(c.expected, rr.Code)
	}
}

func (suite *HttpServerSuite) TestHealth() {

	assert := assert.New(suite.T())

	cases := []struct {
		code int
		body string
	}{
		{
			// success case
			code: http.StatusOK,
			body: "OK",
		},
	}

	for _, c := range cases {

		req, err := http.NewRequest(http.MethodGet, "/health", nil)
		assert.NoError(err)

		rr := httptest.NewRecorder()

		suite.mux.ServeHTTP(rr, req)

		assert.Equal(c.code, rr.Code)

		body, err := io.ReadAll(rr.Body)
		assert.NoError(err)

		assert.Equal(c.body, string(body))
	}
}

func TestHttpServerSuite(t *testing.T) {
	suite.Run(t, new(HttpServerSuite))
}

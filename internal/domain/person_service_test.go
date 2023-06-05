package domain_test

import (
	"testing"

	"github.com/stretchr/testify/suite"
	"github.com/trevatk/go-template/internal/domain"
)

type PersonServiceSuite struct {
	suite.Suite
	service *domain.PersonService
}

func (suite *PersonServiceSuite) SetupTest() {}

func (suite *PersonServiceSuite) TestCreate() {}

func (suite *PersonServiceSuite) TestRead() {}

func (suite *PersonServiceSuite) TestUpdate() {}

func (suite *PersonServiceSuite) TestDelete() {}

func TestPersonServiceSuite(t *testing.T) {
	suite.Run(t, new(PersonServiceSuite))
}

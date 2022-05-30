package neo4jstorage_test

import (
	"goodisgood/internal/app/model"
	neo4jStorage "goodisgood/internal/app/storage/neo4j"
	"testing"

	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

func TestAccountRepository_Create(t *testing.T) {
	d := neo4jStorage.TestDB(t, DBuri, DBUsername, DBPassword)

	s := neo4jStorage.NewStorage(d)
	a := model.TestAccount()
	err := s.Account().Create(a)
	assert.NoError(t, err)
}

func TestAccountRepository_FindByEmail(t *testing.T) {
	d := neo4jStorage.TestDB(t, DBuri, DBUsername, DBPassword)

	s := neo4jStorage.NewStorage(d)

	e2 := "gopher@gopher.go"
	a, err := s.Account().FindByEmail(e2)
	assert.NoError(t, err)
	assert.NotNil(t, a)
}

func TestAccountRepository_Find(t *testing.T) {
	d := neo4jStorage.TestDB(t, DBuri, DBUsername, DBPassword)

	s := neo4jStorage.NewStorage(d)
	u1, err := s.Account().FindByEmail("gopher@gopher.go")
	if err != nil {
		t.Fatal(err)
	}

	e2 := u1.UUID
	a, err := s.Account().Find(e2)
	logrus.New().Info(err)

	assert.NoError(t, err)
	assert.NotNil(t, a)
}

package neo4jstorage_test

import (
	"goodisgood/internal/app/model"
	neo4jStorage "goodisgood/internal/app/storage/neo4j"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	d := neo4jStorage.TestDB(t, DBuri, DBUsername, DBPassword)

	s := neo4jStorage.NewStorage(d)
	a := model.TestAccount()
	u := model.TestUser()

	err := s.Account().Create(a)
	assert.NoError(t, err)

	err = s.User().Create(a.UUID, u)
	assert.NoError(t, err)
}

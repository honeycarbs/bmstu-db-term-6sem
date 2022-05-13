package storage_test

import (
	"goodisgood/internal/app/model"
	"goodisgood/internal/storage"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRepository_Create(t *testing.T) {
	s := storage.TestStore(t, DBuri, DBUsername, DBPassword)

	_, err := s.Account().Create(model.TestAccount())
	assert.NoError(t, err)
	// assert
}

func TestRepository_FindByEmail(t *testing.T) {
	s := storage.TestStore(t, DBuri, DBUsername, DBPassword)
	e1 := "gopher1@gopher.go"

	_, err := s.Account().FindByEmail(e1)
	assert.Error(t, err)

	e2 := "gopher@gopher.go"
	a, err := s.Account().FindByEmail(e2)
	assert.NoError(t, err)
	assert.NotNil(t, a)
}

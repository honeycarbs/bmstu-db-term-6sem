package teststorage_test

import (
	"errors"
	"testing"

	"goodisgood/internal/app/model"
	"goodisgood/internal/app/storage/teststorage"

	"github.com/stretchr/testify/assert"
)

func TestUserRepository_Create(t *testing.T) {
	s := teststorage.NewStorage()
	u := model.TestAccount()
	assert.NoError(t, s.Account().Create(u))
	assert.NotNil(t, u.UUID)
}

func TestUserRepository_FindByEmail(t *testing.T) {
	s := teststorage.NewStorage()
	u1 := model.TestAccount()
	_, err := s.Account().FindByEmail(u1.Email)
	assert.EqualError(t, err, errors.New("not found").Error())

	s.Account().Create(u1)
	u2, err := s.Account().FindByEmail(u1.Email)
	assert.NoError(t, err)
	assert.NotNil(t, u2)
}

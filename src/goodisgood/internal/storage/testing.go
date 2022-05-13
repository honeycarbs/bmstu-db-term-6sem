package storage

import (
	"testing"
)

func TestStore(t *testing.T, DBuri, DBUsername, DBPassword string) *Storage {
	t.Helper()
	config := NewConfig()

	config.DBuri = DBuri
	config.DBUsername = DBUsername
	config.DBPassword = DBPassword

	s := New(config)
	if err := s.Open(); err != nil {
		t.Fatal(err)
	}

	return s
}

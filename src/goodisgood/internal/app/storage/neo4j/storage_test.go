package neo4jstorage_test

import (
	"os"
	"testing"
)

var (
	DBuri,
	DBUsername,
	DBPassword string
)

func TestMain(m *testing.M) {
	DBuri = os.Getenv("DBURI")
	DBUsername = os.Getenv("DBUSERNAME")
	DBPassword = os.Getenv("DBPASSWORD")
	if DBuri == "" {
		DBuri = "bolt://localhost"
	}
	if DBUsername == "" {
		DBUsername = "neo4j"
	}
	if DBPassword == "" {
		DBPassword = "test"
	}

	os.Exit(m.Run())
}

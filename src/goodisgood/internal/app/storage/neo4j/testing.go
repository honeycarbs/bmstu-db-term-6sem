package neo4jstorage

import (
	"testing"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func TestDB(t *testing.T, DBuri, DBUsername, DBPassword string) neo4j.Driver {
	t.Helper()

	db, err := neo4j.NewDriver(DBuri, neo4j.BasicAuth(DBUsername, DBPassword, ""))
	if err != nil {
		t.Fatal(err)
	}
	return db
}

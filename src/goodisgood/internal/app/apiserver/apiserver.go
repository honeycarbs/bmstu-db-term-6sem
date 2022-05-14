package apiserver

import (
	neo4jstorage "goodisgood/internal/app/storage/neo4j"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func Start(config *Config) error {
	db, err := newDB(config.DBURI, config.DBUsername, config.DBPassword)
	if err != nil {
		return err
	}
	defer db.Close()

	storage := neo4jstorage.NewStorage(db)
	sessionStorage := sessions.NewCookieStore([]byte(config.SessionKey))
	s := newServer(storage, sessionStorage)
	return http.ListenAndServe(config.BindAddr, s)
}

func newDB(DBUri, DBUsername, DBPassword string) (neo4j.Driver, error) {
	db, err := neo4j.NewDriver(DBUri, neo4j.BasicAuth(DBUsername, DBPassword, ""))
	if err != nil {
		return nil, err
	}

	return db, nil
}

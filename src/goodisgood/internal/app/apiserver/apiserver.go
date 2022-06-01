package apiserver

import (
	neo4jstorage "goodisgood/internal/app/storage/neo4j"
	"net/http"

	"github.com/gorilla/sessions"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

func Start(conf *Config) error {
	uconf, err := NewDBConfig("config/dbuser.toml")
	if err != nil {
		return err
	}
	aconf, err := NewDBConfig("config/dbadmin.toml")
	if err != nil {
		return err
	}
	mconf, err := NewDBConfig("config/dbmoderator.toml")
	if err != nil {
		return err
	}
	dbu, err := newDB(uconf.DBURI, uconf.DBUsername, uconf.DBPassword)
	if err != nil {
		return err
	}
	defer dbu.Close()

	dba, err := newDB(aconf.DBURI, aconf.DBUsername, aconf.DBPassword)
	if err != nil {
		return err
	}
	defer dba.Close()

	dbm, err := newDB(mconf.DBURI, mconf.DBUsername, mconf.DBPassword)
	if err != nil {
		return err
	}
	defer dba.Close()

	ustorage := neo4jstorage.NewStorage(dbu)
	astorage := neo4jstorage.NewStorage(dba)
	mstorage := neo4jstorage.NewStorage(dbm)

	sessionStorage := sessions.NewCookieStore([]byte(conf.SessionKey))
	s := newServer(ustorage, mstorage, astorage, sessionStorage)

	return http.ListenAndServe(conf.BindAddr, s)
}

func newDB(DBUri, DBUsername, DBPassword string) (neo4j.Driver, error) {
	db, err := neo4j.NewDriver(DBUri, neo4j.BasicAuth(DBUsername, DBPassword, ""))
	if err != nil {
		return nil, err
	}

	return db, nil
}

package neo4jstorage

import (
	"goodisgood/internal/app/model"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type UserRepository struct {
	storage *Storage
}

func (r *UserRepository) Create(auuid string, u *model.User) (err error) {
	if err := u.Validate(); err != nil {
		return err
	}

	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		err = session.Close()
	}()
	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.createQuery(tx, auuid, u)
	}); err != nil {
		return err
	}
	return nil
}

func (r *UserRepository) createQuery(tx neo4j.Transaction, auuid string, u *model.User) (*model.User, error) {
	result, err := tx.Run(
		`match (a:account)
		where a.id=$auuid
		create (u:user{id: apoc.create.uuid(), race: $race, age: $age, gender: $gender})
		create (u)-[o:OWNS]->(a)
		return a.id as id`,
		map[string]interface{}{
			"auuid":  auuid,
			"race":   u.Race,
			"age":    u.Age,
			"gender": u.Gender,
		},
	)
	if err != nil {
		return nil, err
	}

	record, err := result.Single()
	if err != nil {
		return nil, err
	}

	id, _ := record.Get("id")
	u.UUID = id.(string)

	return u, nil
}

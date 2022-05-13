package storage

import (
	"encoding/json"
	"goodisgood/internal/app/model"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j/db"
)

type AccountRepository struct {
	storage *Storage
}

func (r *AccountRepository) Create(a *model.Account) (*model.Account, error) {
	if err := a.Validate(); err != nil {
		return nil, err
	}

	if err := a.EncryptPassword(); err != nil {
		return nil, err
	}
	sess := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer sess.Close()

	result, err := sess.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			"create (a: account{id: apoc.create.uuid(),username: $username, email: $email, password:$password}) return a.id",
			map[string]interface{}{"username": a.Username, "email": a.Email, "password": a.EncryptedPassword})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})
	if err != nil {
		return nil, err
	}

	a.UUID = result.(string)

	return a, nil
}

func (r *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	a := &model.Account{}
	sess := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer sess.Close()

	result, err := sess.WriteTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		result, err := transaction.Run(
			`match (a:account)
				where a.email=$email
				return apoc.convert.toJson({
				    id: a.id,
				    username: a.username,
				    password: a.password
				})`,
			map[string]interface{}{"email": email})
		if err != nil {
			return nil, err
		}

		if result.Next() {
			return result.Record().Values[0], nil
		}

		return nil, result.Err()
	})

	if err != nil {
		return nil, err
	}

	if result == nil {
		return nil, &db.Neo4jError{}
	}

	err = json.Unmarshal([]byte(result.(string)), &a)
	if err != nil {
		return nil, err
	}

	return a, nil
}

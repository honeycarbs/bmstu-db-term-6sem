package neo4jstorage

import (
	"goodisgood/internal/app/model"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
)

type AccountRepository struct {
	storage *Storage
}

func (r *AccountRepository) Create(a *model.Account) error {
	if err := a.Validate(); err != nil {
		return err
	}

	if err := a.EncryptPassword(); err != nil {
		return err
	}
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite,
		DatabaseName: "goodisgood"})
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.createQuery(tx, a)
	}); err != nil {
		logrus.New().Infof("Error is: %v", err != nil)
		return err
	}

	return nil
}

func (r *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite,
		DatabaseName: "goodisgood"})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.findByEmailQuery(tx, email)
	})

	if result == nil {
		return nil, err
	}
	a := &model.Account{}
	a = result.(*model.Account)

	return a, nil
}

func (r *AccountRepository) Find(uuid string) (*model.Account, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite,
		DatabaseName: "goodisgood"})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.findByUUIDQuery(tx, uuid)
	})

	if result == nil {
		return nil, err
	}

	a := &model.Account{}
	a = result.(*model.Account)

	return a, nil
}

func (r *AccountRepository) GetAll() ([]model.Account, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.getAllQuery(tx)
	})

	if err != nil {
		return nil, err
	}

	return result.([]model.Account), nil
}

func (r *AccountRepository) getAllQuery(tx neo4j.Transaction) ([]model.Account, error) {
	aarr := make([]model.Account, 0)
	result, err := tx.Run(
		`match (a:account) return a.id as id, a.email as email, a.password as password, a.username as username, a.role as role`,
		map[string]interface{}{},
	)

	for result.Next() {
		record := result.Record()
		if err != nil {
			return nil, err
		}

		id, _ := record.Get("id")
		email, _ := record.Get("email")
		password, _ := record.Get("password")
		username, _ := record.Get("username")

		l := model.Account{
			UUID:              id.(string),
			Username:          username.(string),
			Email:             email.(string),
			Password:          "",
			EncryptedPassword: password.(string),
		}
		aarr = append(aarr, l)
	}

	return aarr, nil
}

func (r *AccountRepository) createQuery(tx neo4j.Transaction, a *model.Account) (*model.Account, error) {
	result, err := tx.Run(
		`create (a: account{id: apoc.create.uuid(),username: $username, email: $email, password:$password, role:"user"}) 
		return a.id as id`,
		map[string]interface{}{
			"username": a.Username,
			"email":    a.Email,
			"password": a.EncryptedPassword,
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
	a.UUID = id.(string)

	return a, nil
}

func (r *AccountRepository) findByEmailQuery(tx neo4j.Transaction, email string) (*model.Account, error) {
	result, err := tx.Run(
		`match (a:account)
		where a.email=$email
		return a.id as id, a.username as username, a.password as password`,
		map[string]interface{}{
			"email": email,
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
	username, _ := record.Get("username")
	password, _ := record.Get("password")

	a := &model.Account{
		UUID:              id.(string),
		Username:          username.(string),
		Email:             email,
		Password:          "",
		EncryptedPassword: password.(string),
	}

	return a, nil
}

func (r *AccountRepository) findByUUIDQuery(tx neo4j.Transaction, uuid string) (*model.Account, error) {
	result, err := tx.Run(
		`match (a:account)
		where a.id=$uuid
		return a.email as email, a.username as username, a.password as password, a.role as role`,
		map[string]interface{}{
			"uuid": uuid,
		},
	)
	if err != nil {
		return nil, err
	}

	record, err := result.Single()
	if err != nil {
		return nil, err
	}

	email, _ := record.Get("email")
	username, _ := record.Get("username")
	password, _ := record.Get("password")
	role, _ := record.Get("role")

	a := &model.Account{
		UUID:              uuid,
		Username:          username.(string),
		Email:             email.(string),
		Password:          "",
		Role:              role.(string),
		EncryptedPassword: password.(string),
	}

	return a, nil
}

package neo4jstorage

import (
	"goodisgood/internal/app/model"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"github.com/sirupsen/logrus"
)

type AccountRepository struct {
	storage *Storage
}

func (r *AccountRepository) Create(a *model.Account) (err error) {
	if err := a.Validate(); err != nil {
		return err
	}

	if err := a.EncryptPassword(); err != nil {
		return err
	}
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer func() {
		err = session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.createQuery(tx, a)
	}); err != nil {
		return err
	}
	return nil
}

func (r *AccountRepository) FindByEmail(email string) (*model.Account, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
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
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite})
	defer session.Close()

	result, err := session.ReadTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.findByUUIDQuery(tx, uuid)
	})

	if result == nil {
		logrus.New().Infof("Error inside from find: %v", err)
		return nil, err
	}

	a := &model.Account{}
	a = result.(*model.Account)

	return a, nil
}

func (r *AccountRepository) createQuery(tx neo4j.Transaction, a *model.Account) (*model.Account, error) {
	result, err := tx.Run(
		`create (a: account{id: apoc.create.uuid(),username: $username, email: $email, password:$password}) 
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
		return a.email as email, a.username as username, a.password as password`,
		map[string]interface{}{
			"uuid": uuid,
		},
	)
	if err != nil {
		return nil, err
	}

	record, err := result.Single()
	// logrus.New().Info(err)
	if err != nil {
		return nil, err
	}

	email, _ := record.Get("email")
	username, _ := record.Get("username")
	password, _ := record.Get("password")

	a := &model.Account{
		UUID:              uuid,
		Username:          username.(string),
		Email:             email.(string),
		Password:          "",
		EncryptedPassword: password.(string),
	}

	return a, nil
}

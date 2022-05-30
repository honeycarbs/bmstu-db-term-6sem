package neo4jstorage

import (
	"goodisgood/internal/app/model"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type LocationRepository struct {
	storage *Storage
}

func (r *LocationRepository) Create(l *model.Location) (err error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer func() {
		err = session.Close()
	}()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.getLocation(tx, l)
	}); err != nil {
		return err
	}
	return nil
}

func (r *LocationRepository) Assign(auuid string, l *model.Location) error {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.assignToUser(tx, auuid, l)
	}); err != nil {
		return err
	}
	return nil

}

func (r *LocationRepository) Get() ([]model.Location, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.getLocationsQuery(tx)
	})

	if err != nil {
		return nil, err
	}

	return result.([]model.Location), nil
}

func (r *LocationRepository) getLocationsQuery(tx neo4j.Transaction) ([]model.Location, error) {
	larr := make([]model.Location, 0)

	result, err := tx.Run(
		`match (l:location) return 
		l.id as id, l.district as district, l.name as name, l.region as region`,
		map[string]interface{}{},
	)

	if err != nil {
		return nil, err
	}

	for result.Next() {
		record := result.Record()
		if err != nil {
			return nil, err
		}

		id, _ := record.Get("id")
		name, _ := record.Get("name")
		region, _ := record.Get("region")
		district, _ := record.Get("district")

		l := model.Location{
			UUID:     id.(string),
			Name:     name.(string),
			Region:   region.(string),
			District: district.(string),
		}
		larr = append(larr, l)
	}

	return larr, nil
}

func (r *LocationRepository) getLocation(tx neo4j.Transaction, l *model.Location) (*model.Location, error) {
	result, err := tx.Run(
		`match (l:location{name: $name})
		return l.id as id`,
		map[string]interface{}{
			"name":     l.Name,
			"region":   l.Region,
			"district": l.District,
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
	l.UUID = id.(string)

	return l, nil
}

func (r *LocationRepository) assignToUser(tx neo4j.Transaction, auuid string, l *model.Location) (*model.Location, error) {
	_, err := tx.Run(
		`match (u:user)-[o:OWNS]->(a:account{id:$auuid})
		match (l:location{id: $luuid})
		merge (u)-[lw:LIVES]->(l)`,
		map[string]interface{}{
			"auuid": auuid,
			"luuid": l.UUID,
		},
	)
	if err != nil {
		return nil, err
	}
	return l, nil
}

package neo4jstorage

import (
	"goodisgood/internal/app/model"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type EducationPlaceRepository struct {
	storage *Storage
}

func (r *EducationPlaceRepository) Assign(auuid string, e *model.EducationPlace, p *model.EducationProgram) error {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.assignToUser(tx, auuid, e, p)
	}); err != nil {
		return err
	}
	return nil

}

func (r *EducationPlaceRepository) Get() ([]model.EducationPlace, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.getEducationPlacesQuery(tx)
	})

	if err != nil {
		return nil, err
	}

	return result.([]model.EducationPlace), nil
}

func (r *EducationPlaceRepository) getEducationPlacesQuery(tx neo4j.Transaction) ([]model.EducationPlace, error) {
	larr := make([]model.EducationPlace, 0)

	result, err := tx.Run(
		`match (e:education_place)
		return e.id as id, e.name as name`,
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

		l := model.EducationPlace{
			UUID: id.(string),
			Name: name.(string),
		}
		larr = append(larr, l)
	}

	return larr, nil
}

func (r *EducationPlaceRepository) assignToUser(tx neo4j.Transaction, auuid string, er *model.EducationPlace, ep *model.EducationProgram) (*model.EducationPlace, error) {
	_, err := tx.Run(
		`match (u:user)-[o:OWNS]->(a:account{id:$auuid})
		match(e:education_place{name: $ename})
		merge((u)-[s:STUDIES{field:$field, level: $level}]->(e))`,
		map[string]interface{}{
			"auuid": auuid,
			"ename": er.Name,
			"field": ep.Field,
			"level": ep.Level,
		},
	)
	if err != nil {
		return nil, err
	}
	return er, nil
}

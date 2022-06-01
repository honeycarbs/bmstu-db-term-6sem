package neo4jstorage

import (
	"errors"
	"goodisgood/internal/app/model"

	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
)

type PollRepository struct {
	storage *Storage
}

func (r *PollRepository) Submit(auuid string, a *model.Answer) error {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	if _, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.submitPoll(tx, auuid, a)
	}); err != nil {
		return err
	}
	return nil
}

func (r *PollRepository) GetUserResult(auuid string) (*model.Poll, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.getPollResults(tx, auuid)
	})

	if err != nil {
		return nil, err
	}

	return result.(*model.Poll), nil
}

func (r *PollRepository) GetPollStats() ([]model.Stats, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.getPollStats(tx)
	})

	if err != nil {
		return nil, err
	}

	return result.([]model.Stats), nil
}

func (r *PollRepository) GetWordsList() ([]string, error) {
	session := r.storage.db.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: "goodisgood"})
	defer session.Close()

	result, err := session.WriteTransaction(func(tx neo4j.Transaction) (interface{}, error) {
		return r.getWordsList(tx)
	})

	if err != nil {
		return nil, err
	}

	return result.([]string), nil
}

func (r *PollRepository) getWordsList(tx neo4j.Transaction) ([]string, error) {
	result, err := tx.Run(
		`match (w:word)
		return w.name as name`,
		map[string]interface{}{},
	)

	if err != nil {
		return nil, err
	}

	warr := make([]string, 0)

	for result.Next() {
		record := result.Record()
		if err != nil {
			return nil, err
		}

		w, _ := record.Get("name")
		warr = append(warr, w.(string))

	}
	return warr, nil
}

func (r *PollRepository) getPollResults(tx neo4j.Transaction, auuid string) (*model.Poll, error) {
	result, err := tx.Run(
		`match (u:user)-[o:OWNS]->(a:account{id:$auuid})
		optional match (u)-[m:MARKED]->(w:word)
			call apoc.do.when(
				m is null,
				'return null',
				'return w.name as name, m.mark as mark',
				{w: w, m: m})
		yield value
		return value.name as name, value.mark as mark`,
		map[string]interface{}{
			"auuid": auuid,
		},
	)

	// result, err := tx.Run(
	// 	`match (u:user)-[o:OWNS]->(a:account{id:$auuid})
	// 	match (w:word{name: $wname})
	// 	merge (u)-[m:MARKED{mark: $mark}]->(w)
	// 	return m is not null as ok`,
	// 	map[string]interface{}{
	// 		"auuid": auuid,
	// 	},
	// )

	if err != nil {
		return nil, err
	}

	p := &model.Poll{
		Answer: make([]model.Answer, 0),
	}

	for result.Next() {
		record := result.Record()
		if err != nil {
			return nil, err
		}

		word, _ := record.Get("name")
		mark, _ := record.Get("mark")
		a := &model.Answer{
			Word: word.(string),
			Mark: int(mark.(int64)),
		}

		p.Answer = append(p.Answer, *a)
	}
	return p, nil
}

func (r *PollRepository) submitPoll(tx neo4j.Transaction, auuid string, a *model.Answer) (*model.Answer, error) {
	result, err := tx.Run(
		`match (u:user)-[o:OWNS]->(a:account{id:$auuid})
		match (w:word{name: $wname})
		optional match (u)-[m:MARKED]->(w)
		call apoc.do.when(
		    m is null,
		    'create (u)-[m:MARKED{mark: mrk}]->(w) return m is not null as ok',
		    'return true as ok',
		    {u:u, w:w, mrk: $mark}) 
		yield value
		return value.ok as ok`,
		map[string]interface{}{
			"auuid": auuid,
			"mark":  a.Mark,
			"wname": a.Word,
		},
	)
	if err != nil {
		return nil, err
	}

	record, err := result.Single()
	if err != nil {
		return nil, err
	}

	ok, _ := record.Get("ok")

	if !(ok.(bool)) {
		return nil, errors.New("cant't create relationship")
	}

	return a, nil
}

func (r *PollRepository) getPollStats(tx neo4j.Transaction) ([]model.Stats, error) {
	sarr := make([]model.Stats, 0)

	result, err := tx.Run(
		`MATCH ()-[r:MARKED]->(w:word)
		WITH w.name AS name, collect(r.mark) AS marklist, COUNT(r.mark) AS cnt
		unwind marklist as ml
		RETURN DISTINCT name, avg(ml) as avg`,
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

		stat, _ := record.Get("avg")
		name, _ := record.Get("name")

		l := model.Stats{
			Word:  name.(string),
			Stats: stat.(float64),
		}
		sarr = append(sarr, l)
	}

	return sarr, nil
}

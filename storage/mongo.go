package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
	"github.com/DanShu93/martialarts-tracker/query"
)

type MongoRepository struct {
	database *mgo.Database
}

func NewMongoRepository(url, db string) (MongoRepository, error) {
	session, err := mgo.Dial(url)

	if err != nil {
		return MongoRepository{}, err
	}

	database := session.DB(db)

	return MongoRepository{database: database}, nil
}

func (s MongoRepository) Create(collectionName string, data interface{}) error {
	data.(map[string]interface{})["_id"] = data.(map[string]interface{})["ID"]
	delete(data.(map[string]interface{}), "ID")
	err := s.database.C(collectionName).Insert(data)
	if err != nil {
		return err
	}

	return nil
}

func (s MongoRepository) Read(collectionName, id string, result *interface{}) error {
	q := s.database.C(collectionName).Find(bson.M{"_id": id})

	n, err := q.Count()
	if err != nil {
		return Invalid
	}

	if n == 0 {
		return NotFound
	}

	err = q.One(&result)
	if err != nil {
		return Invalid
	}

	(*result).(bson.M)["ID"] = (*result).(bson.M)["_id"]
	delete((*result).(bson.M), "_id")

	return nil
}

func (s MongoRepository) Update(collectionName, id string, data interface{}) error {
	data.(map[string]interface{})["_id"] = data.(map[string]interface{})["ID"]
	delete(data.(map[string]interface{}), "ID")

	err := s.database.C(collectionName).Update(bson.M{"_id": id}, data)
	if err != nil {
		return err
	}

	return nil
}

func (s MongoRepository) Delete(collectionName, id string) error {
	err := s.database.C(collectionName).Remove(bson.M{"_id": id})
	if err != nil {
		return err
	}

	return nil
}

func (s MongoRepository) ReadAll(collectionName string, query query.Query, result *[]interface{}) error {
	mq := createMongoQuery(query)
	q := s.database.C(collectionName).Find(mq)

	n, err := q.Count()
	if err != nil {
		return err
	}

	if n == 0 {
		return NotFound
	}

	err = q.All(result)
	if err != nil {
		return err
	}

	resultSlice := *result

	for i := range resultSlice {
		resultSlice[i].(bson.M)["ID"] = resultSlice[i].(bson.M)["_id"]
		delete(resultSlice[i].(bson.M), "_id")
	}

	return nil
}

func createMongoQuery(q query.Query) bson.M {
	mq := bson.M{}
	or := make([]bson.M, 0)
	and := make([]bson.M, 0)
	for k, v := range q.Q {
		if k == "ID" {
			k = "_id"
		}

		fieldQueries := make([]bson.M, len(v.Values))
		for i, currentValue := range v.Values {
			fieldQueries[i] = bson.M{k: currentValue}
		}

		switch v.Kind {
		case query.KindAnd:
			and = append(and, fieldQueries...)
		case query.KindOr:
			or = append(or, fieldQueries...)
		case query.KindContains:
			mq[k] = bson.M{"$in": v.Values}
		}
	}

	if len(or) != 0 {
		mq["$or"] = or
	}
	if len(and) != 0 {
		mq["$and"] = and
	}

	return mq
}

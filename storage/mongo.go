package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
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
	query := s.database.C(collectionName).Find(bson.M{"_id": id})

	n, err := query.Count()
	if err != nil {
		return Invalid
	}

	if n != 1 {
		return NotFound
	}

	err = query.One(&result)
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

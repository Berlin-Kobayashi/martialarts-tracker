package storage

import (
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoRepository struct {
	collection *mgo.Collection
}

func NewMongoRepository(url, db, collectionName string) (MongoRepository, error) {
	session, err := mgo.Dial(url)

	if err != nil {
		return MongoRepository{}, err
	}

	database := session.DB(db)

	collection := database.C(collectionName)

	return MongoRepository{collection}, nil
}

func (s MongoRepository) Save(data interface{}) error {
	data.(map[string]interface{})["_id"] = data.(map[string]interface{})["ID"]
	delete(data.(map[string]interface{}), "ID")
	err := s.collection.Insert(data)
	if err != nil {
		return err
	}

	return nil
}

func (s MongoRepository) Read(id string, result interface{}) error {
	query := s.collection.Find(bson.M{"_id": id})

	n, err := query.Count()
	if err != nil {
		return Invalid
	}

	if n != 1 {
		return NotFound
	}

	err = query.One(result)
	if err != nil {
		return Invalid
	}

	return nil
}

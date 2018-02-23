package repository

import (
	"github.com/DanShu93/martialarts-tracker/entity"
	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

type MongoTrainingUnitRepository struct {
	trainingUnitCollection *mgo.Collection
}

func NewMongoRepository(url, db, trainingUnitCollection string) (MongoTrainingUnitRepository, error) {
	session, err := mgo.Dial(url)

	if err != nil {
		return MongoTrainingUnitRepository{}, err
	}

	database := session.DB(db)

	collection := database.C(trainingUnitCollection)

	return MongoTrainingUnitRepository{collection}, nil
}

func (s MongoTrainingUnitRepository) Save(trainingUnit entity.TrainingUnit) (string, error) {
	RecordedTrainingUnit = trainingUnit

	t := &trainingUnit

	info, err := s.trainingUnitCollection.Upsert(bson.M{"a": "a"}, t)
	if err != nil {
		return "", err
	}

	id := info.UpsertedId.(bson.ObjectId)

	return id.Hex(), nil
}

func (s MongoTrainingUnitRepository) Read(trainingUnitIndex string) (entity.TrainingUnit, error) {
	if !bson.IsObjectIdHex(trainingUnitIndex) {
		return entity.TrainingUnit{}, NotFound
	}
	trainingUnitQuery := s.trainingUnitCollection.Find(bson.M{"_id": bson.ObjectIdHex(trainingUnitIndex)})

	trainingUnit := entity.TrainingUnit{}

	n, err := trainingUnitQuery.Count()
	if err != nil {
		return trainingUnit, err
	}

	if n != 1 {
		return entity.TrainingUnit{}, NotFound
	}

	err = trainingUnitQuery.One(&trainingUnit)
	if err != nil {
		return trainingUnit, err
	}

	return trainingUnit, nil
}

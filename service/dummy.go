package service

import (
	"github.com/DanShu93/martialarts-tracker/entity"
)

var v4Fixture = "b5e57615-0f40-404e-bbe0-6ae81fe8080a"

type dummyUUIDGenerator struct {
}

func (g dummyUUIDGenerator) Generate() string {
	return v4Fixture
}

var recordedData interface{}

type dummyRepository struct {
}

func (s dummyRepository) Save(data interface{}) error {
	recordedData = data

	return nil
}

func (s dummyRepository) Read(id string, result interface{}) error {
	switch resultPtr := result.(type) {
	case *entity.TrainingUnit:
		*resultPtr = trainingUnitFixture
	}

	return nil
}

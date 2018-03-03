package service

import (
	"testing"
	"reflect"
	"github.com/DanShu93/martialarts-tracker/entity"
)

var indexingRepository = IndexingRepository{
	Generator:  dummyUUIDGenerator{},
	Repository: dummyRepository{},
}

func TestIndexingRepository_Save(t *testing.T) {
	actual := trainingUnitFixture
	err := indexingRepository.Save(&actual)
	if err != nil {
		t.Fatal(err)
	}

	expected := actual
	expected.ID = uuidV4Fixture

	if !reflect.DeepEqual(recordedData, &expected) {
		t.Errorf("Does not save expected data. Actual %v Expected %v", recordedData, &expected)
	}
}

func TestIndexingRepository_Read(t *testing.T) {
	actual := entity.TrainingUnit{}
	err := indexingRepository.Read(trainingUnitFixture.ID, &actual)
	if err != nil {
		t.Fatal(err)
	}

	if !reflect.DeepEqual(actual, trainingUnitFixture) {
		t.Errorf("Does not read expected data. Actual %v Expected %v", actual, trainingUnitFixture)
	}
}

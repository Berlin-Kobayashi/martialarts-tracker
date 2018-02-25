package service

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"github.com/DanShu93/martialarts-tracker/entity"
	"fmt"
	"encoding/json"
	"bytes"
	"reflect"
)

var entityName = "training-unit"

var storageService = StorageService{
	UUIDGenerator: dummyUUIDGenerator{},
	EntityDefinitions: EntityDefinitions{
		entityName: {
			Entity:     &entity.TrainingUnit{},
			Repository: dummyRepository{},
		},
	},
}

func TestStorageService_ServeHTTPGET(t *testing.T) {
	req := httptest.NewRequest("GET", fmt.Sprintf("/%s/%s", entityName, trainingUnitFixture.ID), nil)
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	expectedBody := getTrainingUnitFixtureJSON(t)

	content, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != expectedBody {
		t.Errorf("Does not produce expected response Actual %q Expected %q", content, expectedBody)
	}
}

func TestStorageService_ServeHTTPPOST(t *testing.T) {
	req := httptest.NewRequest("POST", fmt.Sprintf("/%s", entityName), bytes.NewReader([]byte(getTrainingUnitFixtureJSON(t))))
	w := httptest.NewRecorder()
	storageService.ServeHTTP(w, req)

	if !reflect.DeepEqual(recordedData, &trainingUnitFixture) {
		t.Errorf("Does not save expected data Actual %v Expected %v", recordedData, trainingUnitFixture)
	}
}

func getTrainingUnitFixtureJSON(t *testing.T) string {
	fixtureJSON, err := json.Marshal(trainingUnitFixture)

	if err != nil {
		t.Errorf("Could not create training unit fixture")
	}

	return string(fixtureJSON)
}

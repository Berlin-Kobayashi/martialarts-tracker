package service

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"bytes"
	"reflect"
	"github.com/DanShu93/martialarts-tracker/repository"
	"encoding/json"
)

func TestTrainingUnitService_ServeHTTPGET(t *testing.T) {
	trainingUnitService := TrainingUnitService{repository.DummyTrainingUnitRepository{}}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	trainingUnitService.ServeHTTP(w, req)

	expectedBody := getTrainingUnitFixtureJSON(t)

	content, err := ioutil.ReadAll(w.Result().Body)
	if err != nil {
		t.Fatal(err)
	}

	if string(content) != expectedBody {
		t.Errorf("TrainingUnitService GET does not produce expected response Actual %q Expected %q", content, expectedBody)
	}
}

func TestTrainingUnitService_ServeHTTPPOST(t *testing.T) {
	repo := repository.DummyTrainingUnitRepository{}
	trainingUnitService := TrainingUnitService{repo}

	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(getTrainingUnitFixtureJSON(t))))
	w := httptest.NewRecorder()
	trainingUnitService.ServeHTTP(w, req)

	if !reflect.DeepEqual(repository.RecordedTrainingUnit, repository.TrainingUnitFixture) {
		t.Errorf("TrainingUnitService POST does not save expected data Actual %v Expected %v", repository.RecordedTrainingUnit, repository.TrainingUnitFixture)
	}
}

func getTrainingUnitFixtureJSON(t *testing.T) string {
	fixtureJSON, err := json.Marshal(repository.TrainingUnitFixture)

	if err != nil {
		t.Errorf("Could not create training unit fixture")
	}

	return string(fixtureJSON)
}

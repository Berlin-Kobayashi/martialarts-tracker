package service

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"bytes"
	"reflect"
	"encoding/json"
)

func TestTrainingUnitService_ServeHTTPGET(t *testing.T) {
	repo := dummyRepository{}
	uuidGenerator := dummyUUIDGenerator{}
	trainingUnitService := TrainingUnitService{repo, uuidGenerator}

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
	repo := dummyRepository{}
	uuidGenerator := dummyUUIDGenerator{}
	trainingUnitService := TrainingUnitService{repo, uuidGenerator}

	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(getTrainingUnitFixtureJSON(t))))
	w := httptest.NewRecorder()
	trainingUnitService.ServeHTTP(w, req)

	if !reflect.DeepEqual(recordedData, trainingUnitFixture) {
		t.Errorf("TrainingUnitService POST does not save expected data Actual %v Expected %v", recordedData, trainingUnitFixture)
	}
}

func getTrainingUnitFixtureJSON(t *testing.T) string {
	fixtureJSON, err := json.Marshal(trainingUnitFixture)

	if err != nil {
		t.Errorf("Could not create training unit fixture")
	}

	return string(fixtureJSON)
}

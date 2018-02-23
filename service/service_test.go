package service

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"bytes"
	"reflect"
	"github.com/DanShu93/martialarts-tracker/repository"
)

var trainingUnitJSONFixture = "{\"Series\":\"JKD I\",\"Techniques\":[{\"Kind\":\"counter\",\"Name\":\"pak sao\",\"Description\":\"Means slapping hand\\nCounter a jab by slapping the elbow of the opponent into his body, destroying his structure\\nAt the same time perform a jab\"}],\"Methods\":[{\"Kind\":\"counter\",\"Name\":\"Pak Sao drill\",\"Description\":\"\",\"Covers\":[{\"Kind\":\"counter\",\"Name\":\"pak sao\",\"Description\":\"Means slapping hand\\nCounter a jab by slapping the elbow of the opponent into his body, destroying his structure\\nAt the same time perform a jab\"}]}],\"Exercises\":[{\"Kind\":\"Sparring\",\"Name\":\"Lead hand sparring\",\"Description\":\"Sparing with lead hand punches only\"}]}"

func TestTrainingUnitService_ServeHTTPGET(t *testing.T) {
	trainingUnitService := TrainingUnitService{repository.DummyTrainingUnitRepository{}}

	req := httptest.NewRequest("GET", "/", nil)
	w := httptest.NewRecorder()
	trainingUnitService.ServeHTTP(w, req)

	expectedBody := trainingUnitJSONFixture

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

	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(trainingUnitJSONFixture)))
	w := httptest.NewRecorder()
	trainingUnitService.ServeHTTP(w, req)

	if !reflect.DeepEqual(repository.RecordedTrainingUnit, repository.TrainingUnitFixture) {
		t.Errorf("TrainingUnitService POST does not save expected data Actual %v Expected %v", repository.RecordedTrainingUnit, repository.TrainingUnitFixture)
	}
}

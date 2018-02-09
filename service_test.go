package martialarts

import (
	"testing"
	"net/http/httptest"
	"io/ioutil"
	"bytes"
	"reflect"
)

var pakSaoFixture = Technique{
	Kind:        "counter",
	Name:        "pak sao",
	Description: "Means slapping hand\nCounter a jab by slapping the elbow of the opponent into his body, destroying his structure\nAt the same time perform a jab",
}

var trainingUnitFixture = TrainingUnit{
	Series: "JKD I",
	Techniques: []Technique{
		pakSaoFixture,
	},
	Methods: []Method{
		{
			Kind:        "counter",
			Name:        "Pak Sao drill",
			Description: "",
			Covers:      []Technique{pakSaoFixture},
		},
	},
	Exercises: []Exercise{
		{
			Kind:        "Sparring",
			Name:        "Lead hand sparring",
			Description: "Sparing with lead hand punches only",
		},
	},
}

var trainingUnitJSONFixture = "{\"Series\":\"JKD I\",\"Techniques\":[{\"Kind\":\"counter\",\"Name\":\"pak sao\",\"Description\":\"Means slapping hand\\nCounter a jab by slapping the elbow of the opponent into his body, destroying his structure\\nAt the same time perform a jab\"}],\"Methods\":[{\"Kind\":\"counter\",\"Name\":\"Pak Sao drill\",\"Description\":\"\",\"Covers\":[{\"Kind\":\"counter\",\"Name\":\"pak sao\",\"Description\":\"Means slapping hand\\nCounter a jab by slapping the elbow of the opponent into his body, destroying his structure\\nAt the same time perform a jab\"}]}],\"Exercises\":[{\"Kind\":\"Sparring\",\"Name\":\"Lead hand sparring\",\"Description\":\"Sparing with lead hand punches only\"}]}"

var recordedTrainingUnit TrainingUnit

type DummyTrainingUnitRepository struct {
}

func (s DummyTrainingUnitRepository) Save(trainingUnit TrainingUnit) error {
	recordedTrainingUnit = trainingUnit

	return nil
}

func TestTrainingUnitService_ServeHTTPGET(t *testing.T) {
	trainingUnitService := TrainingUnitService{DummyTrainingUnitRepository{}}

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
	repository := DummyTrainingUnitRepository{}
	trainingUnitService := TrainingUnitService{repository}

	req := httptest.NewRequest("POST", "/", bytes.NewReader([]byte(trainingUnitJSONFixture)))
	w := httptest.NewRecorder()
	trainingUnitService.ServeHTTP(w, req)

	if !reflect.DeepEqual(recordedTrainingUnit, trainingUnitFixture) {
		t.Errorf("TrainingUnitService POST does not save expected data Actual %v Expected %v", recordedTrainingUnit, trainingUnitFixture)
	}
}

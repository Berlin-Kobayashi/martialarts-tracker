package martialarts

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type TrainingUnitService struct {
	Repository TrainingUnitRepository
}

func (s TrainingUnitService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		s.get(rw, r)
	case http.MethodPost:
		s.post(rw, r)
	}
}

func (s TrainingUnitService) get(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "application/json")

	pakSao := Technique{
		Kind:        "counter",
		Name:        "pak sao",
		Description: "Means slapping hand\nCounter a jab by slapping the elbow of the opponent into his body, destroying his structure\nAt the same time perform a jab",
	}

	trainingUnit := TrainingUnit{
		Series: "JKD I",
		Techniques: []Technique{
			pakSao,
		},
		Methods: []Method{
			{
				Kind:        "counter",
				Name:        "Pak Sao drill",
				Description: "",
				Covers:      []Technique{pakSao},
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

	jsonString, err := json.Marshal(trainingUnit)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	rw.Write(jsonString)
}

func (s TrainingUnitService) post(rw http.ResponseWriter, r *http.Request) {
	content, err := ioutil.ReadAll(r.Body)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}

	trainingUnit := TrainingUnit{}
	err = json.Unmarshal(content, &trainingUnit)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
	}

	err = s.Repository.Save(trainingUnit)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
	}
}

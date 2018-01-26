package main

import (
	"net/http"
	"encoding/json"
	"io/ioutil"
)

type TrainingUnit struct {
	Series     string
	Techniques []Technique
	Methods    []Method
	Exercises  []Exercise
}

type Technique struct {
	Kind        string
	Name        string
	Description string
}

type Method struct {
	Kind        string
	Name        string
	Description string
	Covers      []Technique
}

type Exercise struct {
	Kind        string
	Name        string
	Description string
}

type FavIconService struct {
}

func (s FavIconService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "image/x-icon")
	http.ServeFile(rw, r, "/go/src/github.com/DanShu93/martialarts-tracker/favicon.ico")
}

type MainService struct {
}

func (s MainService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	rw.Header().Add("Content-Type", "text/html")
	rw.Write([]byte("<link rel=\"shortcut icon\" href=\"http://localhost:8888/favicon.ico\" type=\"image/x-icon\">"))
}

type TrainingUnitService struct {
}

func (s TrainingUnitService) ServeHTTP(rw http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
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
	case http.MethodPost:
		content, err := ioutil.ReadAll(r.Body)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		trainingUnit := TrainingUnit{}
		err = json.Unmarshal(content, &trainingUnit)
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
		}

		jsonString, err := json.Marshal(trainingUnit)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}

		err = ioutil.WriteFile("/go/src/github.com/DanShu93/martialarts-tracker/data/training1.json", jsonString, 0644)
		if err != nil {
			rw.WriteHeader(http.StatusInternalServerError)
		}
	}
}

func main() {
	http.Handle("/favicon.ico", FavIconService{})
	http.Handle("/index.html", MainService{})
	http.Handle("/training-unit", TrainingUnitService{})
	http.ListenAndServe(":80", nil)
}

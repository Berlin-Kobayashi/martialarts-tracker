package main

import (
	"net/http"
	"github.com/DanShu93/jsonmancer/mongo"
	"github.com/DanShu93/jsonmancer/storage"
	"gopkg.in/mgo.v2"
	"github.com/DanShu93/jsonmancer/uuid"
	"reflect"
)

type TrainingUnit struct {
	ID     string `json:"id"`
	Start  string `json:"start"`
	End    string `json:"end"`
	Series string `json:"series"`
}

type Technique struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Method struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

type Exercise struct {
	ID          string `json:"id"`
	Kind        string `json:"kind"`
	Name        string `json:"name"`
	Description string `json:"description"`
}

func main() {
	mongoURL := "db:27017"
	mongoDB := "martialarts"

	storageService, err := build(mongoURL, mongoDB)
	if err != nil {
		panic(err)
	}

	http.Handle("/", storageService)

	http.ListenAndServe(":80", nil)
}

func build(mongoURL, mongoDB string) (http.Handler, error) {
	session, err := mgo.Dial(mongoURL)
	if err != nil {
		panic(err)
	}

	err = session.DB(mongoDB).DropDatabase()
	if err != nil {
		panic(err)
	}

	repository, err := mongo.New(
		mongoURL,
		mongoDB,
	)
	if err != nil {
		panic(err)
	}

	technique := storage.Entity{
		Name:       "technique",
		Data:       reflect.TypeOf(Technique{}),
		References: map[string]storage.Entity{},
	}

	exercise := storage.Entity{
		Name:       "exercise",
		Data:       reflect.TypeOf(Exercise{}),
		References: map[string]storage.Entity{},
	}

	method := storage.Entity{
		Name: "method",
		Data: reflect.TypeOf(Method{}),
		References: map[string]storage.Entity{
			"covers": technique,
		},
	}
	trainingUnit := storage.Entity{
		Name: "trainingunit",
		Data: reflect.TypeOf(TrainingUnit{}),
		References: map[string]storage.Entity{
			"techniques": technique,
			"exercises":  exercise,
			"methods":    method,
		},
	}

	entities := []storage.Entity{
		technique,
		exercise,
		method,
		trainingUnit,
	}

	store, err := storage.New(entities, repository, uuid.V4{})
	if err != nil {
		panic(err)
	}

	return storage.Service{Storage: store, Info: storage.Info{Title: "martialarts-tracker", Version: "0.0.1"}}, nil
}

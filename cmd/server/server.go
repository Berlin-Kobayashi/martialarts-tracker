package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker/storage"
	"github.com/DanShu93/martialarts-tracker/service"
	"github.com/DanShu93/martialarts-tracker/entity"
	"reflect"
	"github.com/DanShu93/martialarts-tracker/uuid"
)

func main() {
	mongoURL := "martialarts-tracker-db:27017"
	mongoDB := "martialarts"

	storageService, err := build(mongoURL, mongoDB)
	if err != nil {
		panic(err)
	}

	http.Handle("/", storageService)

	http.ListenAndServe(":80", nil)
}

func build(mongoURL, mongoDB string) (service.StorageService, error) {
	repository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
	)
	if err != nil {
		return service.StorageService{}, err
	}

	entityDefinitions := service.EntityDefinitions{
		"training-unit": reflect.TypeOf(entity.TrainingUnit{}),
		"technique":     reflect.TypeOf(entity.Technique{}),
		"method":        reflect.TypeOf(entity.Method{}),
		"exercise":      reflect.TypeOf(entity.Exercise{}),
	}

	return service.NewStorageService(entityDefinitions, uuid.V4{}, repository)
}

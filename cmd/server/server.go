package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker/storage"
	"github.com/DanShu93/martialarts-tracker/service"
	"github.com/DanShu93/martialarts-tracker/entity"
	"reflect"
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
	trainingUnitRepository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
		"training_units",
	)
	if err != nil {
		return service.StorageService{}, err
	}

	techniqueRepository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
		"techniques",
	)
	if err != nil {
		return service.StorageService{}, err
	}

	methodRepository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
		"methods",
	)
	if err != nil {
		return service.StorageService{}, err
	}

	exerciseRepository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
		"exercises",
	)
	if err != nil {
		return service.StorageService{}, err
	}

	entityDefinitions := service.EntityDefinitions{
		"training-unit": {
			T: reflect.TypeOf(entity.TrainingUnit{}),
			R: trainingUnitRepository,
		},
		"technique": {
			T: reflect.TypeOf(entity.Technique{}),
			R: techniqueRepository,
		},
		"method": {
			T: reflect.TypeOf(entity.Method{}),
			R: methodRepository,
		},
		"exercise": {
			T: reflect.TypeOf(entity.Exercise{}),
			R: exerciseRepository,
		},
	}

	return service.StorageService{
		EntityDefinitions: entityDefinitions,
	}, nil
}

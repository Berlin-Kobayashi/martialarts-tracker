package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker/storage"
	"github.com/DanShu93/martialarts-tracker/service"
	"github.com/DanShu93/martialarts-tracker/uuid"
	"github.com/DanShu93/martialarts-tracker/entity"
)

func main() {
	mongoURL := "martialarts-tracker-db:27017"
	mongoDB := "martialarts"

	trainingUnitRepository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
		"training_units",
	)
	if err != nil {
		panic(err)
	}

	techniqueRepository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
		"techniques",
	)
	if err != nil {
		panic(err)
	}

	methodRepository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
		"methods",
	)
	if err != nil {
		panic(err)
	}

	exerciseRepository, err := storage.NewMongoRepository(
		mongoURL,
		mongoDB,
		"exercises",
	)
	if err != nil {
		panic(err)
	}

	entityDefinitions := service.EntityDefinitions{
		"training-unit": {
			Entity:     &entity.TrainingUnit{},
			Repository: trainingUnitRepository,
		},
		"technique": {
			Entity:     &entity.Technique{},
			Repository: techniqueRepository,
		},
		"method": {
			Entity:     &entity.Method{},
			Repository: methodRepository,
		},
		"exercise": {
			Entity:     &entity.Exercise{},
			Repository: exerciseRepository,
		},
	}

	uuidGenerator := uuid.V4{}

	trackingService := service.StorageService{
		EntityDefinitions: entityDefinitions,
		UUIDGenerator:     uuidGenerator,
	}

	http.Handle("/training-unit/log", service.LogService{})
	http.Handle("/", trackingService)

	http.ListenAndServe(":80", nil)
}

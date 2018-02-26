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

	uuidGenerator := uuid.V4{}

	entityDefinitions := service.EntityDefinitions{
		"training-unit": {
			Entity: &entity.TrainingUnit{},
			Repository: service.IndexingRepository{
				Repository: trainingUnitRepository,
				Generator:  uuidGenerator,
			},
		},
		"technique": {
			Entity: &entity.Technique{},
			Repository: service.IndexingRepository{
				Repository: techniqueRepository,
				Generator:  uuidGenerator,
			},
		},
		"method": {
			Entity: &entity.Method{},
			Repository: service.IndexingRepository{
				Repository: methodRepository,
				Generator:  uuidGenerator,
			},
		},
		"exercise": {
			Entity: &entity.Exercise{},
			Repository: service.IndexingRepository{
				Repository: exerciseRepository,
				Generator:  uuidGenerator,
			},
		},
		"training-unit/log": {
			Entity: &service.Log{},
			Repository: service.LogRepository{
				TrainingUnitRepository: trainingUnitRepository,
				TechniqueRepository:    techniqueRepository,
				ExerciseRepository:     exerciseRepository,
				MethodRepository:       methodRepository,
			},
		},
	}

	return service.StorageService{
		EntityDefinitions: entityDefinitions,
	}, nil
}

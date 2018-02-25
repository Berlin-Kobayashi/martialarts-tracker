package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker/storage"
	"github.com/DanShu93/martialarts-tracker/service"
	"github.com/DanShu93/martialarts-tracker/uuid"
	"github.com/DanShu93/martialarts-tracker/entity"
)

func main() {
	trainingUnitRepository, err := storage.NewMongoRepository(
		"martialarts-tracker-db:27017",
		"martialarts",
		"training_units",
	)

	if err != nil {
		panic(err)
	}

	entityDefinitions := service.EntityDefinitions{
		"training-unit": {
			Entity:     &entity.TrainingUnit{},
			Repository: trainingUnitRepository,
		},
	}

	uuidGenerator := uuid.V4{}

	trackingService := service.StorageService{
		EntityDefinitions: entityDefinitions,
		UUIDGenerator:     uuidGenerator,
	}

	http.Handle("/", trackingService)

	http.ListenAndServe(":80", nil)
}

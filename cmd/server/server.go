package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker/storage"
	"github.com/DanShu93/martialarts-tracker/service"
	"github.com/DanShu93/martialarts-tracker/uuid"
)

func main() {
	repo, err := storage.NewMongoRepository(
		"martialarts-tracker-db:27017",
		"martialarts",
		"training_units",
	)

	if err != nil {
		panic(err)
	}

	uuidGenerator := uuid.V4{}

	trainingUnitService := service.TrainingUnitService{Repository: repo, UUIDGenerator: uuidGenerator}

	http.Handle("/training-unit/", trainingUnitService)
	http.Handle("/training-unit", trainingUnitService)

	http.ListenAndServe(":80", nil)
}

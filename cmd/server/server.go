package main

import (
	"net/http"
	"github.com/DanShu93/martialarts-tracker/repository"
	"github.com/DanShu93/martialarts-tracker/service"
)

func main() {
	repo, err := repository.NewMongoRepository(
		"martialarts-tracker-db:27017",
		"martialarts",
		"training_units",
	)

	if err != nil {
		panic(err)
	}

	http.Handle("/training-unit/", service.TrainingUnitService{Repository: repo})
	http.Handle("/training-unit", service.TrainingUnitService{Repository: repo})

	http.ListenAndServe(":80", nil)
}
